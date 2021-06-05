package soiprompt

import "strings"

// toStorableName はファイル名を保存可能な形式に変換します
func toStorableName(n string) string {
	// pathの予約語系を変換します
	n = strings.ReplaceAll(n, " ", "_")
	n = strings.ReplaceAll(n, "/", "／")
	// 拡張子がなければ追加します
	if !strings.HasSuffix(n, ".json") {
		n = n + ".json"
	}
	return n
}
