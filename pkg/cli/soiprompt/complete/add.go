package complete

import (
	"golang.org/x/net/context"
	"log"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/utils"
)

// addCmd はaddコマンド系のSuggestを提示します
func (c *Completer) addCmd(d prompt.Document) []prompt.Suggest {
	// remove cache
	c.cache.Clear()

	// option探索
	if utils.IsOptionWord(d) {
		return []prompt.Suggest{
			{Text: "-n", Description: "name of the url"},
			{Text: "-d", Description: "dir to which model store"},
			{Text: "-t", Description: "tags to the url (allows multiple options)"},
		}
	}
	if strings.HasSuffix(d.Text, "-n ") {
		return EmptySuggests
	}
	// dir探索
	if strings.HasSuffix(d.Text, "-d ") {
		var suggests []prompt.Suggest
		soiRoot, err := c.Bucket.Path()
		if err != nil {
			log.Fatal(err)
		}
		dirs, err := c.service.ListPath(context.Background(), "", false) // TODO fix context flow
		if err != nil {
			log.Fatal(err)
		}
		for _, d := range dirs {
			suggests = append(suggests, prompt.Suggest{
				Text:        strings.TrimPrefix(d, soiRoot+"/"),
				Description: "",
			})
		}
		return suggests
	}
	if strings.HasSuffix(d.Text, " ") {
		return []prompt.Suggest{
			{Text: "https://", Description: "target url"},
			{Text: "-n", Description: "name of the url"},
			{Text: "-d", Description: "dir to which model store"},
			{Text: "-t", Description: "tags to the url (allows multiple options)"},
		}
	}
	return EmptySuggests
}
