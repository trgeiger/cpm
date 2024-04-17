package app

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-ini/ini"
)

const (
	CoprUrl  string    = "https://copr.fedorainfracloud.org/coprs/"
	CoprHost string    = "copr.fedorainfracloud.org"
	ReposDir string    = "/etc/yum.repos.d/"
	Enabled  RepoState = "enabled=1"
	Disabled RepoState = "enabled=0"
)

func FedoraReleaseVersion() string {
	osRelease, err := ini.Load("/etc/os-release")
	if err != nil {
		log.Fatal("Fail to read file: ", err)
	}

	return osRelease.Section("").Key("VERSION_ID").String()
}

func HandleError(err error) {
	if errors.Is(err, fs.ErrPermission) {
		fmt.Printf("This command must be run with superuser privileges.\nError: %s\n", err)
	} else {
		fmt.Println(err)
	}
}

func GetLocalRepoFileLines(r *CoprRepo) ([]string, error) {
	repoFile := r.LocalFilePath()
	contents, err := os.ReadFile(repoFile)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(contents), "\n"), nil
}

func WriteRepoToFile(r *CoprRepo, content []byte) error {
	err := os.WriteFile(r.LocalFilePath(), content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ToggleRepo(r *CoprRepo, desiredState RepoState) error {
	fileLines, err := GetLocalRepoFileLines(r)
	if err != nil {
		return err
	}
	var statusMessage string
	if desiredState == Enabled {
		statusMessage = "enabled"
	} else {
		statusMessage = "disabled"
	}

	for i, line := range fileLines {
		if strings.Contains(line, "enabled=") {
			if line == string(desiredState) {
				fmt.Printf("Repository is already %s.\n", statusMessage)
				return nil
			} else {
				fileLines[i] = string(desiredState)
			}
		}
	}
	output := strings.Join(fileLines, "\n")
	err = WriteRepoToFile(r, []byte(output))
	if err != nil {
		return err
	}
	fmt.Printf("Repository %s/%s %s.\n", r.User, r.Project, statusMessage)
	return nil
}

func AddRepo(r *CoprRepo) error {
	resp, err := http.Get(r.RepoConfigUrl())
	if err != nil {
		return err
	}
	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = WriteRepoToFile(r, []byte(output))
	if err != nil {
		return err
	}
	fmt.Printf("Repository %s/%s added.\n", r.User, r.Project)
	return nil
}

func DeleteRepo(r *CoprRepo) error {
	if r.LocalFileExists() {
		err := os.Remove(r.LocalFilePath())
		if err != nil {
			return err
		}
		fmt.Printf("Repository %s/%s deleted.\n", r.User, r.Project)
	} else {
		fmt.Printf("Repository %s/%s does not exist locally. Nothing to delete.\n", r.User, r.Project)
	}
	return nil
}
