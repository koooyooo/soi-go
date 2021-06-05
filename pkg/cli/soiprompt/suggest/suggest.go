package suggest

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

var EmptySuggests []prompt.Suggest

// Completer はSuggest候補を提示することで補完を実施します
func Completer(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()
	switch {
	case hasPrefixes(text, "add ", "a "):
		return addCmd(d)
	case hasPrefixes(text, "mv "):
		return mvCmd(d)
	case hasPrefixes(text, "rm "):
		return rmCmd(d)
	case hasPrefixes(text, "dig ", "d "):
		return digCmd(d)
	case hasPrefixes(text, "list ", "ls", "l "):
		return listCmd(d)
	case hasPrefixes(text, "cb ", "c "):
		return cbCmd(d)
	case hasPrefixes(text, "help ", "h "):
		return helpCmd(d)
	default:
		s := []prompt.Suggest{
			{Text: "add", Description: "(a)dd url"},
			{Text: "dig", Description: "(d)ig urls"},
			{Text: "list", Description: "(l)ist and filter urls"},
			{Text: "cb", Description: "change bucket"},
			{Text: "tag", Description: "add tags (TODO)"},
			{Text: "mv", Description: "move file to dir"},
			{Text: "rm", Description: "remove file or dir"},
			//{Text: "tags", Description: "lists up all tags"}, TODO implements as "list -t"
			{Text: "pull", Description: "pull urls from repos"},
			{Text: "push", Description: "push urls to repos"},
			{Text: "help", Description: "(h)elp document"},
			{Text: "quit", Description: "(q)uit soi"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}
}

// hasPrefixes は引数の文字列に接頭語が含まれているものがあるかを調査します
func hasPrefixes(in string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(in, p) {
			return true
		}
	}
	return false
}
