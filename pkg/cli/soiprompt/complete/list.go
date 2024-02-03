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

var listOptSuggests = []prompt.Suggest{
	{Text: "-p", Description: "open in private mode"},
	{Text: "-c", Description: "open w/ chrome"},
	{Text: "-f", Description: "open w/ firefox"},
	{Text: "-s", Description: "open w/ safari"},
	{Text: "-n", Description: "sort by num-views"},
	{Text: "-a", Description: "sort by add-day"},
	{Text: "-v", Description: "sort by view-day"},
}
