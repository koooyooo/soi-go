package constant

import (
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// LocalBucket はローカルバケットを表現するグローバル変数
var LocalBucket = localBucket{
	name: "default",
}

// localBucket はローカルバケットを管理する構造体
type localBucket struct {
	name string
}

// IsLocalOnly はローカル限定のバケット
func (l localBucket) IsLocalOnly() bool {
	return strings.HasPrefix(l.name, "_")
}

// SetName はローカルバケットの名前を設定する
func (l *localBucket) SetName(n string) {
	l.name = n
}

// GetName はローカルバケットの名前を取得する
func (l localBucket) GetName() string {
	return l.name
}

// Path はローカルバケット毎のルートパスを取得する
func (l localBucket) Path() (string, error) {
	soisDir, err := soisDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(soisDir, l.name), nil
}

// soisDir はSoiのルートディレクトリを取得します
func soisDir() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".soi"), nil
}
