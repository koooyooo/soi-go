package complete

import (
	"strings"

	"golang.org/x/net/context"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/soiprompt/utils"
)

// addCmd はaddコマンド系のSuggestを提示します
func (c *Completer) addCmd(d prompt.Document) []prompt.Suggest {
	c.cache.Clear()

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
	if strings.HasSuffix(d.Text, "-d ") {
		var suggests []prompt.Suggest
		soiRoot, err := c.Bucket.Path()
		if err != nil {
			return EmptySuggests
		}
		dirs, err := c.service.ListPath(context.Background(), "", false)
		if err != nil {
			return EmptySuggests
		}
		for _, dir := range dirs {
			suggests = append(suggests, prompt.Suggest{
				Text:        strings.TrimPrefix(dir, soiRoot+"/"),
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
