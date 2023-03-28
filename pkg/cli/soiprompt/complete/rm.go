package complete

import (
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
	// ファイル系を追加
	files, err := utils.ListFilesRecursively(dir)
	if err != nil {
		log.Fatal(err)
	}
	fileDirs = append(fileDirs, files...)
	// ディレクトリ系を追加
	dirs, err := utils.ListDirsRecursively(dir, c.Bucket, false)
	if err != nil {
		log.Fatal(err)
	}
	fileDirs = append(fileDirs, dirs...)
	sort.Strings(fileDirs)

	return utils.FilePathsToSuggests(dir, fileDirs, word)
}
