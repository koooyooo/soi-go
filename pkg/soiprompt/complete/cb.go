package complete

import (
	"golang.org/x/net/context"

	"github.com/c-bata/go-prompt"
)

// cbCmd はバケット変更時のSuggestを提示します
func (c *Completer) cbCmd(d prompt.Document) []prompt.Suggest {
	buckets, err := c.service.ListBucket(context.Background())
	if err != nil {
		return EmptySuggests
	}
	var s []prompt.Suggest
	for _, b := range buckets {
		s = append(s, prompt.Suggest{
			Text:        b,
			Description: "existing bucket",
		})
	}
	s = append(s, prompt.Suggest{
		Text:        "<<new bucket name>>",
		Description: `input new bucket name (local private one should start with "_")`,
	})
	return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
}
