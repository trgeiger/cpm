package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
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
		if errors.Is(err, afero.ErrFileNotFound) {
			return fmt.Errorf("repository %s/%s is not installed", r.User, r.Project)
		}
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

func AddRepo(r *CoprRepo, fs afero.Fs, out io.Writer, multi bool) error {
	resp, err := http.Get(r.RepoConfigUrl(fs, multi))
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
		err := fs.Remove(r.LocalFilePath())
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "Repository %s/%s deleted.\n", r.User, r.Project)
	} else {
		fmt.Fprintf(out, "Repository %s/%s does not exist locally. Nothing to delete.\n", r.User, r.Project)
	}
	return nil
}

func GetReposList(fs afero.Fs, out io.Writer, state RepoState) ([]*CoprRepo, error) {
	files, err := afero.ReadDir(fs, ReposDir)
	if err != nil {
		return nil, err
	}
	var reposStrings []string
	var repos []*CoprRepo
	for _, file := range files {
		if !file.IsDir() {
			var repoName string
			addToResult := false
			isCoprRepo := false
			ioFile, err := fs.Open(ReposDir + file.Name())

			if err != nil {
				return nil, err
			}

			scanner := bufio.NewScanner(ioFile)

			for scanner.Scan() {
				// If we see Copr repo, store name in repoName
				if strings.Contains(scanner.Text(), "[copr:copr") {
					t := strings.Split(strings.Trim(scanner.Text(), "[]"), ":")
					repoName = t[len(t)-2] + "/" + t[len(t)-1]
					isCoprRepo = true
				}
				// If we see our desired state, flip our flag
				if strings.Contains(scanner.Text(), string(state)) && isCoprRepo {
					addToResult = true
				}
			}
			if addToResult && !slices.Contains(reposStrings, repoName) {
				r, err := NewCoprRepo(repoName)
				if err != nil {
					return nil, err
				}
				repos = append(repos, r)
				reposStrings = append(reposStrings, repoName)
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(out, "Issue reading repo files: ", err)
			}
		}
	}
	return repos, nil
}
