package soiprompt

import (
	"fmt"
	"strings"

	"github.com/koooyooo/soi-go/pkg/soi"
)

func createHeader(sd *soi.SoiData) string {
	return fmt.Sprintf("[%3d %04.1f]", sd.NumViews, sd.NumReads)
}

// 相対パス内のヘッダ部分を除去
func removeHeader(s string) string {
	if strings.HasPrefix(s, "[") && strings.Contains(s, "]") {
		idx := strings.Index(s, "]")
		s = s[idx+len("]"):]
	}
	return s
}
