package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

func NewCoprRepo(repoPath string) (CoprRepo, error) {
	repo := CoprRepo{
		User:    strings.Split(repoPath, "/")[0],
		Project: strings.Split(repoPath, "/")[1],
	}
	return repo, nil
}

func FedoraReleaseVersion() string {
	osRelease, err := ini.Load("/etc/os-release")
	if err != nil {
		log.Fatal("Fail to read file: ", err)
	}

	return osRelease.Section("").Key("VERSION_ID").String()
}

func RepoFileUrl(r CoprRepo) *url.URL {
	fedoraRelease := "fedora-" + FedoraReleaseVersion()
	repoName := r.User + "-" + r.Project + "-" + fedoraRelease + ".repo"
	base, err := url.Parse(CoprUrl)
	if err != nil {
		log.Fatal(err)
	}
	repoUrl := base.JoinPath(r.User, r.Project, "repo", fedoraRelease, repoName)
	return repoUrl
}

func RepoFileName(r CoprRepo) string {
	fileName := strings.Join([]string{"_copr", CoprHost, r.User, r.Project + ".repo"}, ":")
	return fileName
}

func RepoFilePath(r CoprRepo) string {
	return ReposDir + RepoFileName(r)
}

func RepoExists(r CoprRepo) bool {
	_, err := os.Stat(RepoFilePath(r))
	return !os.IsNotExist(err)
}

func GetLocalRepoFileLines(r CoprRepo) ([]string, error) {
	repoFile := RepoFilePath(r)
	contents, err := os.ReadFile(repoFile)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(contents), "\n"), nil
}

func WriteRepoToFile(r CoprRepo, content []byte) error {
	err := os.WriteFile(RepoFilePath(r), content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ToggleRepo(r CoprRepo, desiredState RepoState) error {
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

func AddRepo(r CoprRepo) error {
	resp, err := http.Get(RepoFileUrl(r).String())
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

func DeleteRepo(r CoprRepo) error {
	if RepoExists(r) {
		err := os.Remove(RepoFilePath(r))
		if err != nil {
			return err
		}
		fmt.Printf("Repository %s/%s deleted.\n", r.User, r.Project)
	} else {
		fmt.Printf("Repository %s/%s does not exist locally. Nothing to delete.\n", r.User, r.Project)
	}
	return nil
}
