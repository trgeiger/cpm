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
	"github.com/spf13/afero"
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

func HandleError(err error, out io.Writer) {
	if errors.Is(err, fs.ErrPermission) {
		fmt.Fprintf(out, "This command must be run with superuser privileges.\nError: %s\n", err)
	} else {
		fmt.Fprintln(out, err)
	}
}

func GetLocalRepoFileLines(r *CoprRepo, fs afero.Fs) ([]string, error) {
	repoFile := r.LocalFilePath()
	contents, err := afero.ReadFile(fs, repoFile)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(contents), "\n"), nil
}

func WriteRepoToFile(r *CoprRepo, fs afero.Fs, content []byte) error {
	err := afero.WriteFile(fs, r.LocalFilePath(), content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ToggleRepo(r *CoprRepo, fs afero.Fs, out io.Writer, desiredState RepoState) error {
	fileLines, err := GetLocalRepoFileLines(r, fs)
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
				fmt.Fprintf(out, "Repository %s is already %s.\n", r.Name(), statusMessage)
				return nil
			} else {
				fileLines[i] = string(desiredState)
			}
		}
	}
	output := strings.Join(fileLines, "\n")
	err = WriteRepoToFile(r, fs, []byte(output))
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "Repository %s/%s %s.\n", r.User, r.Project, statusMessage)
	return nil
}

func AddRepo(r *CoprRepo, fs afero.Fs, out io.Writer) error {
	resp, err := http.Get(r.RepoConfigUrl())
	if err != nil {
		return err
	}
	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = WriteRepoToFile(r, fs, []byte(output))
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "Repository %s/%s added.\n", r.User, r.Project)
	return nil
}

func DeleteRepo(r *CoprRepo, fs afero.Fs, out io.Writer) error {
	if r.LocalFileExists(fs) {
		err := os.Remove(r.LocalFilePath())
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "Repository %s/%s deleted.\n", r.User, r.Project)
	} else {
		fmt.Fprintf(out, "Repository %s/%s does not exist locally. Nothing to delete.\n", r.User, r.Project)
	}
	return nil
}
