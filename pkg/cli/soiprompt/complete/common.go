package complete

import (
	"context"
	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/complete/soisort"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/view"
	"log"
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
