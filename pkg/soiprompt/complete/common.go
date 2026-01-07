package complete

import (
	"context"
	"log"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/soiprompt/complete/soisort"
	"github.com/koooyooo/soi-go/pkg/soiprompt/view"
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
		text = strings.TrimPrefix(text, cmd)
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
	// スペースで分割して単語ごとに処理
	words := strings.Fields(text)
	var filteredWords []string

	// オプション一覧を作成
	optionSet := make(map[string]bool)
	for _, s := range listOptSuggests {
		optionSet[s.Text] = true
	}

	// オプションでない単語のみを残す
	for _, word := range words {
		if !optionSet[word] {
			filteredWords = append(filteredWords, word)
		}
	}

	return strings.Join(filteredWords, " ")
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
