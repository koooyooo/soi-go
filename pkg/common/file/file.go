/*
Package file offers file related functions
*/
package file

import (
	"io"
	"os"
	"strings"
)

// Exists returns true if the file or dir exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func IsDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}

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
