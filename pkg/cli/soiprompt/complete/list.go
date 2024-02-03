package complete

import (
	"golang.org/x/net/context"
	"log"
	"soi-go/pkg/cli/soiprompt/utils"
	"strings"

	"soi-go/pkg/cli/soiprompt/view"

	"soi-go/pkg/cli/soiprompt/complete/soisort"

	"github.com/c-bata/go-prompt"
)

// listCmd はlistコマンド系のSuggestを提示します
func (c *Completer) listCmd(d prompt.Document) []prompt.Suggest {
	// option探索
	if utils.IsOptionWord(d) {
		return browserOptSuggests
	}
	soisDir, err := c.Bucket.Path()
	if err != nil {
		log.Fatal(err)
	}
	input := removeOption(removeCmd(d.TextBeforeCursor()))
	if len(c.cache.ListSoiCache) == 0 {
		sois, err := c.service.LoadAll(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		c.cache.ListSoiCache = sois
	}

	soisort.Exec(c.cache.ListSoiCache, d)

	var sgs []prompt.Suggest
	for _, s := range c.cache.ListSoiCache {
		sgs = append(sgs, prompt.Suggest{
			Text:        view.ToLine(s, soisDir),
			Description: "",
		})
	}

	words := strings.Split(input, " ")
	return filterByMultiWords(words, sgs)
}

var browserOptSuggests = []prompt.Suggest{
	{Text: "-p", Description: "open in private mode"},
	{Text: "-c", Description: "open w/ chrome"},
	{Text: "-f", Description: "open w/ firefox"},
	{Text: "-s", Description: "open w/ safari"},
	{Text: "-n", Description: "sort by num-views"},
	{Text: "-a", Description: "sort by add-day"},
	{Text: "-v", Description: "sort by view-day"},
}

func removeCmd(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimLeft(text, "list ")
	text = strings.TrimLeft(text, "l ")
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
