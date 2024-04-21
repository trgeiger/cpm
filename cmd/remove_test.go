package cmd

import (
	"bytes"
	"testing"

	"github.com/trgeiger/cpm/internal/testutil"
)

func TestRemoveCmd(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		repoFiles  [][]string // format: file/reponame, test directory folder
		otherFiles [][]string // format: filename, path, test directory folder
		expected   string
	}{
		{
			name: "Remove invalid repo name",
			args: []string{
				"cpm",
			},
			expected: "invalid repository name: cpm\n",
		},
		{
			name: "Remove uninstalled repo",
			args: []string{
				"example/example",
			},
			expected: "Repository example/example does not exist locally. Nothing to delete.\n",
		},
		{
			name: "Remove installed repo",
			args: []string{
				"kylegospo/bazzite",
			},
			repoFiles: [][]string{
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite.repo", "enabled"},
			},
			expected: "Repository kylegospo/bazzite deleted.\n",
		},
	}

	for _, test := range tests {

		b := new(bytes.Buffer)
		fs := testutil.AssembleTestFs(test.repoFiles, test.otherFiles)
		cmd := NewRemoveCmd(fs, b)
		cmd.SetOut(b)
		cmd.SetArgs(test.args)

		cmd.Execute()

		if b.String() != test.expected {
			t.Fatalf("Test \"%s\" failed", test.name)
		}
	}
}
