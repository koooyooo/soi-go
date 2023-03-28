package complete

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/c-bata/go-prompt"
	"github.com/stretchr/testify/assert"
)

func TestSuggestByPath(t *testing.T) {
	wd, _ := os.Getwd()
	prjDir := strings.TrimSuffix(wd, "/pkg/cli/soiprompt/complete")
	soisDir4Test := prjDir + "/testfiles/bucket1"

	tests := []struct {
		name     string
		soisDir  string
		path     string
		input    string
		showDir  bool
		expected []prompt.Suggest
	}{
		{
			name:    "known-dir",
			soisDir: soisDir4Test,
			path:    soisDir4Test + "/portal",
			input:   "google",
			showDir: true,
			expected: []prompt.Suggest{
				{
					Text:        "portal/google",
					Description: "",
				},
			},
		},
		{
			name:     "unknown-dir",
			soisDir:  "/unknown/bucket",
			path:     "/unknown/bucket/dig-target",
			input:    "target",
			showDir:  true,
			expected: []prompt.Suggest(nil),
		},
	}

	for _, test := range tests {
		fmt.Println(os.Getwd())
		t.Run(test.name, func(t *testing.T) {
			suggests := suggestByPath(test.soisDir, test.path, test.input, test.showDir)
			assert.Equal(t, test.expected, suggests)
		})
	}
}
