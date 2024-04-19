package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/trgeiger/copr-tool/internal/testutil"
)

func TestEnableCmd(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		repoFiles  [][]string // format: file/reponame, test directory folder
		otherFiles [][]string // format: filename, path, test directory folder
		expected   string
	}{
		{
			name: "Add valid repo",
			args: []string{
				"kylegospo/bazzite",
			},
			expected: "Repository kylegospo/bazzite added.\n",
			otherFiles: [][]string{
				{"os-release", "/etc/", "f40"},
			},
		},
		{
			name: "Add invalid repo name",
			args: []string{
				"copr-tool",
			},
			expected: "invalid repository name: copr-tool\n",
		},
		{
			name: "Repo does not exist",
			args: []string{
				"example/example",
			},
			expected: "repository does not exist, https://copr.fedorainfracloud.org/coprs/example/example returned 404\n",
		},
		{
			name: "Repo does not support Fedora version",
			args: []string{
				"kylegospo/bazzite",
			},
			otherFiles: [][]string{
				{"os-release", "/etc/", "f30"},
			},
			expected: "repository kylegospo/bazzite does not support Fedora release 30\n",
		},
		{
			name: "Repo already exists and already enabled",
			args: []string{
				"kylegospo/bazzite",
			},
			repoFiles: [][]string{
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite.repo", "enabled"},
			},
			expected: "Repository kylegospo/bazzite is already enabled.\n",
		},
		{
			name: "Repo already exists but not enabled",
			args: []string{
				"kylegospo/bazzite",
			},
			repoFiles: [][]string{
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite.repo", "disabled"},
			},
			expected: "Repository kylegospo/bazzite enabled.\n",
		},
	}

	for _, test := range tests {

		b := new(bytes.Buffer)
		fs := testutil.AssembleTestFs(test.repoFiles, test.otherFiles)
		cmd := NewEnableCmd(fs, b)
		cmd.SetOut(b)
		cmd.SetArgs(test.args)

		cmd.Execute()

		outB := b.String()
		fmt.Print(outB)
		if b.String() != test.expected {
			t.Fatalf("Test: \"%s\" failed", test.name)
		}
	}

}
