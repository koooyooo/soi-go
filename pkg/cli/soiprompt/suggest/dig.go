package suggest

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/suggest/common"

	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/utils"
	"github.com/koooyooo/soi-go/pkg/fileio"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/constant"
)

// DigCmd はppコマンド系のSuggestを提示します
func DigCmd(d prompt.Document) []prompt.Suggest {
	// option探索
	if utils.IsOptionWord(d) {
		return browserOptSuggests
	}
	input := d.TextBeforeCursor()
	inputs := strings.Split(input, " ")

	flags := flag.NewFlagSet("dig", flag.PanicOnError)
	flags.Bool("f", false, "open w/ firefox")
	flags.Bool("s", false, "open w/ safari")
	if err := flags.Parse(inputs[1:]); err != nil {
		log.Fatal(err)
	}

	soisDir, err := constant.LocalBucket.Path()
	if err != nil {
		log.Fatal(err)
	}
	return suggestByPath(soisDir, filepath.Join(soisDir, flags.Arg(0)), d.GetWordBeforeCursor(), true)
}

func suggestByPath(soisDir, path, input string, showDir bool) []prompt.Suggest {
	var found bool
	isDir, err := fileio.IsDir(strings.TrimSuffix(path, "/"))
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			found = false
			// 対象ファイルが見つからないだけの場合はスルー
		default:
			log.Fatal(err)
		}
	} else {
		found = true
	}
	switch {
	case !found || isDir:
		if !found {
			path = toLeafDirPath(path)
		}
		dirs, err := utils.ListFileDirs(path, showDir, true)
		if err != nil {
			log.Fatal(err)
		}
		return utils.FilePathsToSuggests(soisDir, dirs, input)
	}
	return common.EmptySuggests
}

// toLeafDirPath はPathを末端ディレクトリのPathに変換します
func toLeafDirPath(path string) string {
	lastSlashIdx := strings.LastIndex(path, "/")
	if lastSlashIdx == -1 {
		return ""
	}
	return path[0:lastSlashIdx]
}
