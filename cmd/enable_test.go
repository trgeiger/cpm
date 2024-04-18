package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/trgeiger/copr-tool/app"
)

func TestEnableCmd(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		localFiles [][]string
		expected   string
	}{
		{
			name: "Add valid repo",
			args: []string{
				"kylegospo/bazzite",
			},
			expected: "Repository kylegospo/bazzite added.\n",
		},
		{
			name: "Repo already exists and already enabled",
			args: []string{
				"kylegospo/bazzite",
			},
			localFiles: [][]string{
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite.repo", "enabled"},
			},
			expected: "Repository kylegospo/bazzite is already enabled.\n",
		},
		{
			name: "Repo already exists but not enabled",
			args: []string{
				"kylegospo/bazzite",
			},
			localFiles: [][]string{
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite.repo", "disabled"},
			},
			expected: "Repository kylegospo/bazzite enabled.\n",
		},
	}

	for _, test := range tests {
		fs := afero.NewMemMapFs()
		fs.Mkdir("/etc/yum.repos.d/", 0755)
		b := new(bytes.Buffer)
		cmd := NewEnableCmd(fs, b)
		cmd.SetOut(b)
		cmd.SetArgs(test.args)
		if test.localFiles != nil {
			for _, file := range test.localFiles {
				localFs := afero.NewOsFs()
				testFile, _ := afero.ReadFile(localFs, "./test/"+file[0]+"-"+file[1])
				_ = afero.WriteFile(fs, app.ReposDir+file[0], testFile, 0755)
			}
		}
		cmd.Execute()
		poop := b.String()
		fmt.Print(poop)
		if b.String() != test.expected {
			t.Fatalf("Test: \"%s\" failed", test.name)
		}
	}

}
