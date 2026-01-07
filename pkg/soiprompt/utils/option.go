package utils

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

func IsOptionWord(d prompt.Document) bool {
	text := d.TextBeforeCursor()

	// ハイフンで終わる場合（例: "list -" または "list -n -"）
	// これは明らかにオプション入力中
	if strings.HasSuffix(text, "-") {
		return true
	}

	// スペース + ハイフンで終わる場合（例: "list -n -"）
	// これも明らかにオプション入力中
	if strings.HasSuffix(text, " -") {
		return true
	}

	// それ以外の場合（例: "list -n " や "list -n google"）は
	// オプション入力中ではなく、通常の候補を表示すべき
	return false
}
