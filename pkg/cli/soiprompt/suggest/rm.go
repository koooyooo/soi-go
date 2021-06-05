package suggest

import (
	"log"
	"sort"
	"strings"

	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/utils"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/constant"
)

// RmCmd はrmコマンド系のSuggestを提示します
func RmCmd(d prompt.Document) []prompt.Suggest {
	word := strings.TrimPrefix(d.GetWordBeforeCursor(), "rm ")

	dir, err := constant.LocalBucket.Path()
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
	dirs, err := utils.ListDirsRecursively(dir, false)
	if err != nil {
		log.Fatal(err)
	}
	fileDirs = append(fileDirs, dirs...)
	sort.Strings(fileDirs)

	return utils.FilePathsToSuggests(dir, fileDirs, word)
}
