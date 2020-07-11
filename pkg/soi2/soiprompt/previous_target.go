package soiprompt

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

// PreviousTarget は事前の実行結果を記録します
type PreviousTarget struct {
	Path     string
	Suggests []prompt.Suggest
}

// filter は入力を元にSuggestをフィルタリングします
func (p PreviousTarget) filter(lastInput string) []prompt.Suggest {
	var filtered []prompt.Suggest
	for _, s := range p.Suggests {
		if strings.HasPrefix(s.Text, lastInput) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

var previousTarget PreviousTarget
