package complete

import (
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/soiprompt/utils"
)

// listCmd はlistコマンド系のSuggestを提示します
func (c *Completer) listCmd(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()

	// オプション入力中の場合はオプション候補のみ表示
	if utils.IsOptionWord(d) {
		return listOptSuggests
	}

	// オプションが含まれているがオプション入力中でない場合は、
	// 元の候補とオプション候補の両方を表示
	if containsOptions(text) && strings.HasSuffix(text, " ") {
		baseSuggests := c.baseList(d, "list", "ls", "l")
		// 元の候補とオプション候補をマージ
		allSuggests := make([]prompt.Suggest, 0, len(baseSuggests)+len(listOptSuggests))
		allSuggests = append(allSuggests, baseSuggests...)
		allSuggests = append(allSuggests, listOptSuggests...)
		return allSuggests
	}

	// 通常の場合は元の候補のみ表示
	return c.baseList(d, "list", "ls", "l")
}

// containsOptions はテキストにオプションが含まれているかチェック
func containsOptions(text string) bool {
	words := strings.Fields(text)
	for i := 1; i < len(words); i++ {
		if strings.HasPrefix(words[i], "-") {
			return true
		}
	}
	return false
}
