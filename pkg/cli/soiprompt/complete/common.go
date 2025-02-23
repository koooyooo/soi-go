package complete

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/complete/soisort"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/view"
)

func (c *Completer) baseList(d prompt.Document, commands ...string) []prompt.Suggest {
	soisDir, err := c.Bucket.Path()
	if err != nil {
		log.Fatal(err)
	}
	input := removeOption(removeCmd(d.TextBeforeCursor(), commands...))
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

func removeCmd(text string, commands ...string) string {
	text = strings.TrimSpace(text)
	for _, cmd := range commands {
		strings.TrimPrefix(text, cmd+" ")
	}
	return text
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

func removeOption(text string) string {
	// TODO 通常の順に並べ文字数順にソートするロジックに変更する
	// TODO 両端にスペースを付けて引っ掛けて、半角スペースと置換する
	var options []string
	for _, s := range listOptSuggests {
		options = append(options, s.Text)
	}
	for _, opt := range options {
		text = strings.ReplaceAll(text, fmt.Sprintf(" %s ", opt), " ")
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
