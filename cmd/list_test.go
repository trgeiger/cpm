package cmd

import (
	"bytes"
	"testing"

	"github.com/trgeiger/copr-tool/internal/testutil"
)

func TestListCmd(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		repoFiles  [][]string // format: file/reponame, test directory folder
		otherFiles [][]string // format: filename, path, test directory folder
		expected   string
	}{
		{
			name: "List existing repos",
			repoFiles: [][]string{
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite.repo", "enabled"},
				{"_copr:copr.fedorainfracloud.org:bieszczaders:kernel-cachyos.repo", "enabled"},
			},
			expected: "- List of enabled repositories:\nbieszczaders/kernel-cachyos\nkylegospo/bazzite\n",
		},
		{
			name:     "No repos to list",
			expected: "- No enabled repositories\n",
		},
		{
			name: "List disabled and enabled repos",
			args: []string{"--all"},
			repoFiles: [][]string{
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite.repo", "disabled"},
				{"_copr:copr.fedorainfracloud.org:bieszczaders:kernel-cachyos.repo", "enabled"},
			},
			expected: "- List of enabled repositories:\n" +
				"bieszczaders/kernel-cachyos\n" +
				"\n" +
				"- List of disabled repositories:\n" +
				"kylegospo/bazzite\n",
		},
	}

	for _, test := range tests {

		b := new(bytes.Buffer)
		fs := testutil.AssembleTestFs(test.repoFiles, test.otherFiles)
		cmd := NewListCmd(fs, b)
		cmd.SetOut(b)
		cmd.SetArgs(test.args)

		cmd.Execute()
		if b.String() != test.expected {
			t.Fatalf("Test \"%s\" failed", test.name)
		}
	}

}
