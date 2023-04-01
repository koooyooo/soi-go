package complete

import (
	"golang.org/x/net/context"
	"log"
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/utils"
)

// rmCmd はrmコマンド系のSuggestを提示します
func (c *Completer) rmCmd(d prompt.Document) []prompt.Suggest {
	word := strings.TrimPrefix(d.GetWordBeforeCursor(), "rm ")

	dir, err := c.Bucket.Path()
	if err != nil {
		log.Fatal(err)
	}
	var fileDirs []string
	ctx := context.Background() // TODO fix context flow
	paths, err := c.service.ListPath(ctx, "", true)
	if err != nil {
		log.Fatal(err)
	}
	fileDirs = append(fileDirs, paths...)
	sort.Strings(fileDirs)

	return utils.FilePathsToSuggests(dir, fileDirs, word)
}
