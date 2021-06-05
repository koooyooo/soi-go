package suggest

import (
	"log"
	"strings"

	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/utils"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/constant"
)

// mvCmd はmvコマンド系のSuggestを提示します
func mvCmd(d prompt.Document) []prompt.Suggest {
	text := d.Text
	is2ndArg := 2 < len(strings.Split(text, " "))

	word := strings.TrimPrefix(d.GetWordBeforeCursor(), "mv ")

	dir, err := constant.LocalBucket.Path()
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	if is2ndArg {
		files, err = utils.ListDirsRecursively(dir, true)
	} else {
		files, err = utils.ListFilesRecursively(dir)
	}
	if err != nil {
		log.Fatal(err)
	}
	return utils.FilePathsToSuggests(dir, files, word)
}
