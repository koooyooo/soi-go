package complete

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/utils"
	"github.com/koooyooo/soi-go/pkg/common/file"
)

// digCmd はppコマンド系のSuggestを提示します
func (c *Completer) digCmd(d prompt.Document) []prompt.Suggest {
	// option探索
	if utils.IsOptionWord(d) {
		return browserOptSuggests
	}

	soisDir, err := c.Bucket.Path()
	if err != nil {
		log.Fatal(err)
	}

	digPath := removeOption(strings.Split(d.TextBeforeCursor(), " ")[1])
	return suggestByPath(soisDir, filepath.Join(soisDir, digPath), d.GetWordBeforeCursor(), true)
}

func suggestByPath(soisDir, path, input string, showDir bool) []prompt.Suggest {
	var pathErr bool
	isDir, err := file.IsDir(path)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			pathErr = true
		default:
		}
	}
	if !pathErr && !isDir {
		return EmptySuggests
	}
	if pathErr {
		path = toLeafDirPath(path)
	}
	dirs, err := utils.ListFileDirs(path, showDir, true)
	if err != nil {
		return EmptySuggests
	}
	return utils.FilePathsToSuggestsNoEx(soisDir, dirs, input)
}

// toLeafDirPath はPathを末端ディレクトリのPathに変換します
func toLeafDirPath(path string) string {
	lastSlashIdx := strings.LastIndex(path, "/")
	if lastSlashIdx == -1 {
		return ""
	}
	return path[0:lastSlashIdx]
}
