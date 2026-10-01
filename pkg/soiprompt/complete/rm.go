package complete

import (
	"sort"
	"strings"

	"golang.org/x/net/context"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/soiprompt/utils"
)

// rmCmd はrmコマンド系のSuggestを提示します
func (c *Completer) rmCmd(d prompt.Document) []prompt.Suggest {
	word := strings.TrimPrefix(d.GetWordBeforeCursor(), "rm ")

	dir, err := c.Bucket.Path()
	if err != nil {
		return EmptySuggests
	}
	var fileDirs []string
	ctx := context.Background()
	paths, err := c.service.ListPath(ctx, "", true)
	if err != nil {
		return EmptySuggests
	}
	fileDirs = append(fileDirs, paths...)
	sort.Strings(fileDirs)

	return utils.FilePathsToSuggests(dir, fileDirs, word)
}
