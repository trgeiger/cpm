package cmd

import (
	"bytes"
	"testing"

	"github.com/trgeiger/copr-tool/internal/testutil"
)

func TestDisableCmd(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		repoFiles  [][]string // format: file/reponame, test directory folder
		otherFiles [][]string // format: filename, path, test directory folder
		expected   string
	}{
		{
			name: "Disable invalid repo name",
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
			expected: "repository example/example is not installed\n",
		},
		{
			name: "Repo already exists and already disabled",
			args: []string{
				"kylegospo/bazzite",
			},
			repoFiles: [][]string{
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite.repo", "disabled"},
			},
			expected: "Repository kylegospo/bazzite is already disabled.\n",
		},
		{
			name: "Repo already exists but not disabled",
			args: []string{
				"kylegospo/bazzite",
			},
			repoFiles: [][]string{
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite.repo", "enabled"},
			},
			expected: "Repository kylegospo/bazzite disabled.\n",
		},
	}

	for _, test := range tests {

		b := new(bytes.Buffer)
		fs := testutil.AssembleTestFs(test.repoFiles, test.otherFiles)
		cmd := NewDisableCmd(fs, b)
		cmd.SetOut(b)
		cmd.SetArgs(test.args)

		cmd.Execute()

		if b.String() != test.expected {
			t.Fatalf("Test \"%s\" failed", test.name)
		}
	}
}
