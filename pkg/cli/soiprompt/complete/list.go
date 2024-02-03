package complete

import (
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/utils"
	"strings"

	"github.com/c-bata/go-prompt"
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

func removeCmd(text string, commands ...string) string {
	text = strings.TrimSpace(text)
	for _, cmd := range commands {
		text = strings.TrimLeft(text, cmd+" ")
	}
	return text
}

func removeOption(text string) string {
	// TODO 通常の順に並べ文字数順にソートするロジックに変更する
	options := []string{"-p", "-c", "-f", "-s", "-n", "-a", "-v"}
	for _, opt := range options {
		text = strings.ReplaceAll(text, opt, "")
	}
	return strings.TrimSpace(text)
}

func filterByMultiWords(words []string, filtered []prompt.Suggest) []prompt.Suggest {
	for _, word := range words {
		if word == "" {
			continue
		}
		filtered = prompt.FilterContains(filtered, word, true)
	}
	return filtered
}
