package complete

import (
	"context"
	"github.com/c-bata/go-prompt"
	"log"
	"soi-go/pkg/cli/soiprompt/complete/soisort"
	"soi-go/pkg/cli/soiprompt/view"
	"strings"
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
