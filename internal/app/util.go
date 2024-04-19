package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

const (
	CoprUrl  string    = "https://copr.fedorainfracloud.org/coprs/"
	CoprHost string    = "copr.fedorainfracloud.org"
	ReposDir string    = "/etc/yum.repos.d/"
	Enabled  RepoState = "enabled=1"
	Disabled RepoState = "enabled=0"
)

func FedoraReleaseVersion(fs afero.Fs) string {
	// osRelease, err := ini.Load("/etc/os-release")
	reader := viper.New()
	reader.SetFs(fs)
	reader.SetConfigName("os-release")
	reader.SetConfigType("ini")
	reader.AddConfigPath("/etc/")
	reader.ReadInConfig()
	osRelease := reader.GetString("default.version_id")

	return osRelease
}

func SudoMessage(err error, out io.Writer) {
	if errors.Is(err, fs.ErrPermission) {
		fmt.Fprintf(out, "This command must be run with superuser privileges.\nError: %s\n", err)
	} else {
		fmt.Fprintln(out, err)
	}
}

func WriteRepoToFile(r *CoprRepo, fs afero.Fs, content []byte) error {
	err := afero.WriteFile(fs, r.LocalFilePath(), content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ToggleRepo(r *CoprRepo, fs afero.Fs, out io.Writer, desiredState RepoState) error {
	repoFile := r.LocalFilePath()
	contents, err := afero.ReadFile(fs, repoFile)
	if err != nil {
		return err
	}
	fileLines := strings.Split(string(contents), "\n")

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
	resp, err := http.Get(r.RepoConfigUrl(fs))
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

func GetAllRepos(fs afero.Fs) ([]*CoprRepo, error) {
	files, err := os.ReadDir(ReposDir)
	if err != nil {
		return nil, err
	}
	var reposStrings []string
	var repos []*CoprRepo
	for _, file := range files {
		if !file.IsDir() {
			ioFile, err := os.Open(ReposDir + file.Name())

			if err != nil {
				return nil, err
			}

			scanner := bufio.NewScanner(ioFile)
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), "[copr:copr") {
					t := strings.Split(strings.Trim(scanner.Text(), "[]"), ":")
					// r, _ := app.NewCoprRepo(t[len(t)-2] + "/" + t[len(t)-1])
					repoName := t[len(t)-2] + "/" + t[len(t)-1]
					if !slices.Contains(reposStrings, repoName) {
						r, err := NewCoprRepo(repoName)
						if err != nil {
							return nil, err
						}
						repos = append(repos, r)
						reposStrings = append(reposStrings, repoName)
					}
					break
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Issue reading repo files: ", err)
			}
		}
	}
	return repos, nil
}
