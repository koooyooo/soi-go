package execute

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/koooyooo/soi-go/pkg/cli/constant"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/common"
	"github.com/koooyooo/soi-go/pkg/soi"
)

// open は指定されたSoiを元にブラウザを開きます
func open(in string) error {
	flags := flag.NewFlagSet("open", flag.PanicOnError)
	firefox := flags.Bool("f", false, "use firefox")
	safari := flags.Bool("s", false, "use safari")
	if err := flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}

	// Soiファイルを特定
	soisDir, err := constant.LocalBucket.Path()
	if err != nil {
		return err
	}

	relPath := addJSON(common.RemoveHeader(filepath.Join(flags.Args()...)))

	fullPath := filepath.Join(soisDir, relPath)

	// Soiファイルを読み込み
	b, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	var s soi.SoiData
	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	// 閲覧履歴を追記
	s.NumViews++

	// 利用ログを記載
	s.UsageLogs = append(s.UsageLogs, soi.UsageLog{
		Type:   soi.UsageTypeOpen,
		UsedAt: time.Now(),
	})

	// Soiファイルを再登録
	ub, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(fullPath, ub, 0600); err != nil {
		return err
	}

	// 環境設定に応じてブラウザオープン
	if *firefox {
		return exec.Command("open", "-a", "Firefox", s.URI).Start()
	}
	if *safari {
		return exec.Command("open", "-a", "Safari", s.URI).Start()
	}
	return exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", s.URI).Start()
}

// 末尾に ".json" を追加します
func addJSON(path string) string {
	if !strings.HasSuffix(path, ".json") {
		path = path + ".json"
	}
	return path
}
