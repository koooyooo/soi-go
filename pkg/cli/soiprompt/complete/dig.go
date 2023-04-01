package complete

import (
	"golang.org/x/net/context"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/utils"
	"github.com/koooyooo/soi-go/pkg/common/file"
)

// digCmd はppコマンド系のSuggestを提示します
func (c *Completer) digCmd(d prompt.Document) []prompt.Suggest {
	if utils.IsOptionWord(d) {
		return browserOptSuggests
	}
	digPath := removeOption(strings.Split(d.TextBeforeCursor(), " ")[1])

	if len(c.digPathCache) == 0 {
		paths, err := c.service.ListPath(context.Background(), digPath, true)
		if err != nil {
			log.Fatal(err)
		}
		c.digPathCache = paths
	}

	var suggests []prompt.Suggest
	for _, nextPath := range nextElmPath(c.digPathCache, digPath) {
		suggests = append(suggests, prompt.Suggest{Text: nextPath})
	}
	return suggests
}

func nextElmPath(paths []string, part string) []string {
	numElms := strings.Count(part, "/") + 1
	var nextPathMap = make(map[string]struct{})
	for _, path := range paths {
		if !strings.HasPrefix(path, part) {
			continue
		}
		elms := strings.Split(path, "/")
		nextPath := strings.Join(elms[:numElms], "/")
		if len(nextPath) < len(path) {
			nextPath = nextPath + "/"
		}
		nextPathMap[nextPath] = struct{}{}
	}
	var nextPaths []string
	for k, _ := range nextPathMap {
		nextPaths = append(nextPaths, k)
	}
	sort.Strings(nextPaths)
	return nextPaths
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
