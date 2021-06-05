package meta

import (
	"fmt"
	"strings"

	"github.com/koooyooo/soi-go/pkg/soi"
)

// Create はメタヘッダを作成します
func Create(sd *soi.SoiData) string {
	return fmt.Sprintf("[%3d %04.1f]", sd.NumViews, sd.NumReads)
}

// Remove は文字列からメタヘッダを除去します
func Remove(s string) string {
	if strings.HasPrefix(s, "[") && strings.Contains(s, "]") {
		idx := strings.Index(s, "]")
		s = s[idx+len("]"):]
	}
	return s
}
