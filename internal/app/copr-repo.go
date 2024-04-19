package app

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

type (
	RepoState string
)

type CoprRepo struct {
	User       string
	Project    string
	LocalFiles []string
}

func NewCoprRepo(repoName string) (*CoprRepo, error) {
	if matched, _ := regexp.MatchString(`\w*\/\w*`, repoName); !matched {
		return nil, fmt.Errorf("invalid repository name: %s", repoName)
	}
	repo := &CoprRepo{
		User:    strings.Split(repoName, "/")[0],
		Project: strings.Split(repoName, "/")[1],
	}
	return repo, nil
}

func (c *CoprRepo) Name() string {
	return strings.Join([]string{c.User, c.Project}, "/")
}

func (c *CoprRepo) RepoUrl() string {
	base, err := url.Parse(CoprUrl)
	if err != nil {
		log.Fatal(err)
	}
	return base.JoinPath(c.Name()).String()
}

func (c *CoprRepo) RemoteFileName(fs afero.Fs) string {
	return strings.Join([]string{c.User, c.Project, FedoraReleaseVersion(fs)}, "-") + ".repo"
}

func (c *CoprRepo) RepoConfigUrl(fs afero.Fs) string {
	fedoraRelease := "fedora-" + FedoraReleaseVersion(fs)
	base, err := url.Parse(c.RepoUrl())
	if err != nil {
		log.Fatal(err)
	}
	repoUrl := base.JoinPath("repo", fedoraRelease, c.RemoteFileName(fs))
	return repoUrl.String()
}

func (c *CoprRepo) DefaultLocalFileName() string {
	fileName := strings.Join([]string{"_copr", CoprHost, c.User, c.Project + ".repo"}, ":")
	return fileName
}

func (c *CoprRepo) LocalFilePath() string {
	return ReposDir + c.DefaultLocalFileName()
}

func (c *CoprRepo) LocalFileExists(fs afero.Fs) bool {
	_, err := fs.Stat(c.LocalFilePath())
	return !os.IsNotExist(err)
}

func (c *CoprRepo) FindLocalFiles(fs afero.Fs) error {
	files, err := afero.ReadDir(fs, ReposDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		result, err := afero.FileContainsBytes(fs, ReposDir+file.Name(), []byte(c.Name()))
		if err != nil {
			return err
		}
		if result {
			c.LocalFiles = append(c.LocalFiles, file.Name())
		}
	}
	return nil
}

func (c *CoprRepo) PruneDuplicates(fs afero.Fs, out io.Writer) (bool, error) {
	if len(c.LocalFiles) == 0 {
		fmt.Fprintf(out, "Repository %s is not installed.", c.Name())
	} else if len(c.LocalFiles) > 1 {
		if _, err := fs.Open(ReposDir + c.DefaultLocalFileName()); err != nil {
			err := fs.Rename(ReposDir+c.LocalFiles[0], ReposDir+c.DefaultLocalFileName())
			if err != nil {
				return false, err
			}
			c.LocalFiles[0] = c.DefaultLocalFileName()
		}
		pruneCount := 0
		for _, fileName := range c.LocalFiles {
			if fileName != c.DefaultLocalFileName() {
				err := fs.Remove(ReposDir + fileName)
				if err != nil {
					return true, err
				}
				pruneCount++
				//TODO remove the element from LocalFiles
			}
		}
		fmt.Fprintf(out, "Pruned %d duplicate entries for %s\n", pruneCount, c.Name())
		return true, nil
	}
	return false, nil
}
