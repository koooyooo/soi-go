package complete

import (
	"strings"

	"golang.org/x/net/context"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/soiprompt/utils"
)

// mvCmd はmvコマンド系のSuggestを提示します
func (c *Completer) mvCmd(d prompt.Document) []prompt.Suggest {
	word := strings.TrimPrefix(d.GetWordBeforeCursor(), "mv ")
	is2ndArg := 2 < len(strings.Split(d.Text, " "))
	paths, err := c.service.ListPath(context.Background(), word, !is2ndArg)
	if err != nil {
		return EmptySuggests
	}
	dir, err := c.Bucket.Path()
	if err != nil {
		return EmptySuggests
	}
	return utils.FilePathsToSuggests(dir, paths, word)
}
