package complete

import (
	"log"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/utils"
)

// mvCmd はmvコマンド系のSuggestを提示します
func (c *Completer) mvCmd(d prompt.Document) []prompt.Suggest {
	text := d.Text
	is2ndArg := 2 < len(strings.Split(text, " "))

	word := strings.TrimPrefix(d.GetWordBeforeCursor(), "mv ")

	dir, err := c.Bucket.Path()
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	if is2ndArg {
		files, err = utils.ListDirsRecursively(dir, c.Bucket, true)
	} else {
		files, err = utils.ListFilesRecursively(dir)
	}
	if err != nil {
		log.Fatal(err)
	}
	return utils.FilePathsToSuggests(dir, files, word)
}
