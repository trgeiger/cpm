package testutil

import (
	"github.com/spf13/afero"
	"github.com/trgeiger/cpm/internal/app"
)

func AssembleTestFs(repoFiles [][]string, otherFiles [][]string) afero.Fs {
	fs := afero.NewMemMapFs()
	fs.Mkdir("/etc/yum.repos.d/", 0755)
	localFs := afero.NewOsFs()
	for _, file := range repoFiles {
		testFile, _ := afero.ReadFile(localFs, "./test/"+file[1]+"/"+file[0])
		_ = afero.WriteFile(fs, app.ReposDir+file[0], testFile, 0755)
	}
	for _, file := range otherFiles {
		testFile, _ := afero.ReadFile(localFs, "./test/"+file[2]+"/"+file[0])
		_ = afero.WriteFile(fs, file[1]+file[0], testFile, 0755)
	}
	return fs
}
