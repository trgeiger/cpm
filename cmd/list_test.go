package cmd

import (
	"bytes"
	"testing"

	"github.com/trgeiger/copr-tool/internal/testutil"
)

func TestListCmd(t *testing.T) {
	tests := []struct {
		name       string
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
			expected: "bieszczaders/kernel-cachyos\nkylegospo/bazzite\n",
		},
		{
			name:     "No repos to list",
			expected: "No installed Copr repositories.\n",
		},
	}

	for _, test := range tests {

		b := new(bytes.Buffer)
		fs := testutil.AssembleTestFs(test.repoFiles, test.otherFiles)
		cmd := NewListCmd(fs, b)
		cmd.SetOut(b)

		cmd.Execute()

		if b.String() != test.expected {
			t.Fatalf("Test \"%s\" failed", test.name)
		}
	}

}
