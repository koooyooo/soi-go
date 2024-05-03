package complete

import (
	"golang.org/x/net/context"
	"log"
	"strings"

	"github.com/c-bata/go-prompt"
	"soi-go/pkg/cli/soiprompt/utils"
)

// digCmd はppコマンド系のSuggestを提示します
func (c *Completer) digCmd(d prompt.Document) []prompt.Suggest {
	if utils.IsOptionWord(d) {
		return listOptSuggests
	}
	digPath := removeOption(strings.Split(d.TextBeforeCursor(), " ")[1])

	if len(c.cache.DigPathCache) == 0 {
		paths, err := c.service.ListPath(context.Background(), digPath, true)
		if err != nil {
			log.Fatal(err)
		}
		c.cache.DigPathCache = paths
	}

	var suggests []prompt.Suggest
	for _, nextPath := range nextElmPath(c.cache.DigPathCache, digPath) {
		suggests = append(suggests, prompt.Suggest{Text: nextPath})
	}
	return suggests
}

func nextElmPath(paths []string, part string) []string {
	numElmsOfPart := strings.Count(part, "/") + 1
	var result []string
	for _, path := range paths {
		if !strings.HasPrefix(path, part) {
			continue
		}
		numElmsOfPath := len(strings.Split(strings.TrimSuffix(path, "/"), "/"))
		if numElmsOfPath != numElmsOfPart {
			continue
		}
		result = append(result, path)
	}
	return result
}

// toLeafDirPath はPathを末端ディレクトリのPathに変換します
func toLeafDirPath(path string) string {
	lastSlashIdx := strings.LastIndex(path, "/")
	if lastSlashIdx == -1 {
		return ""
	}
	return path[0:lastSlashIdx]
}
