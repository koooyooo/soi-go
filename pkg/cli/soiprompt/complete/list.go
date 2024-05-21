package complete

import (
	"github.com/c-bata/go-prompt"
	"soi-go/pkg/cli/soiprompt/utils"
)

// listCmd はlistコマンド系のSuggestを提示します
func (c *Completer) listCmd(d prompt.Document) []prompt.Suggest {
	if utils.IsOptionWord(d) {
		return listOptSuggests
	}
	return c.baseList(d, "list", "ls", "l")
}
