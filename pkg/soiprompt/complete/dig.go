package complete

import (
	"strings"

	"golang.org/x/net/context"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/soiprompt/utils"
)

// digCmd はppコマンド系のSuggestを提示します
func (c *Completer) digCmd(d prompt.Document) []prompt.Suggest {
	if utils.IsOptionWord(d) {
		return listOptSuggests
	}
	parts := strings.Split(d.TextBeforeCursor(), " ")
	digPath := ""
	if len(parts) > 1 {
		digPath = removeOption(parts[1])
	}

	if len(c.cache.DigPathCache) == 0 {
		paths, err := c.service.ListPath(context.Background(), digPath, true)
		if err != nil {
			return EmptySuggests
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
