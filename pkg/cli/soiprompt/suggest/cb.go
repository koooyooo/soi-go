package suggest

import (
	"log"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/constant"
)

// CBCmd はバケット変更時のSuggestを提示します
func CBCmd(d prompt.Document) []prompt.Suggest {
	buckets, err := constant.ListBuckets()
	if err != nil {
		log.Fatal(err)
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
