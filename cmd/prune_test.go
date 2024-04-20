package cmd

import (
	"bytes"
	"testing"

	"github.com/trgeiger/copr-tool/internal/testutil"
)

func TestPruneCmd(t *testing.T) {
	tests := []struct {
		name       string
		repoFiles  [][]string // format: file/reponame, test directory folder
		otherFiles [][]string // format: filename, path, test directory folder
		expected   string
	}{
		{
			name:     "No repositories installed",
			expected: "Nothing to prune.\n",
		},
		{
			name: "Remove 1 duplicate",
			repoFiles: [][]string{
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite.repo", "enabled"},
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite-copy.repo", "enabled"},
			},
			expected: "Removed 1 duplicate entry for kylegospo/bazzite.\n",
		},
		{
			name: "Remove multiple duplicates",
			repoFiles: [][]string{
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite.repo", "enabled"},
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite-copy.repo", "enabled"},
				{"_copr:copr.fedorainfracloud.org:kylegospo:bazzite-copy2.repo", "enabled"},
			},
			expected: "Removed 2 duplicate entries for kylegospo/bazzite.\n",
		},
	}

	for _, test := range tests {

		b := new(bytes.Buffer)
		fs := testutil.AssembleTestFs(test.repoFiles, test.otherFiles)
		cmd := NewPruneCmd(fs, b)
		cmd.SetOut(b)

		cmd.Execute()

		if b.String() != test.expected {
			t.Fatalf("Test \"%s\" failed", test.name)
		}
	}
}
