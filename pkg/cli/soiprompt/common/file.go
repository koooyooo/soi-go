package common

import "strings"

// ToStorableName はファイル名を保存可能な形式に変換します
func ToStorableName(n string) string {
	// pathの予約語系を変換します
	n = strings.ReplaceAll(n, " ", "_")
	n = strings.ReplaceAll(n, "/", "／")
	// 拡張子がなければ追加します
	if !strings.HasSuffix(n, ".json") {
		n = n + ".json"
	}
	return n
}
