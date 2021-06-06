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
	"github.com/koooyooo/soi-go/pkg/soi"
)

// open は指定されたSoiを元にブラウザを開きます
func open(in string) error {
	flags := flag.NewFlagSet("open", flag.PanicOnError)
	chrome := flags.Bool("c", false, "use chrome")
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

	relPath := addJSONSuffix(flags.Arg(flags.NArg() - 1))
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
	if *chrome {
		return openChrome(s)
	}
	if *firefox {
		return openFirefox(s)
	}
	if *safari {
		return openSafari(s)
	}
	defB := strings.ToLower(constant.EnvKeyDefaultBrowser.Get())
	if defB == "chrome" {
		return openChrome(s)
	}
	if defB == "firefox" {
		return openFirefox(s)
	}
	if defB == "safari" {
		return openSafari(s)
	}
	return exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", s.URI).Start()
}

func openChrome(s soi.SoiData) error {
	return exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", s.URI).Start()
}

func openFirefox(s soi.SoiData) error {
	return exec.Command("open", "-a", "Firefox", s.URI).Start()
}

func openSafari(s soi.SoiData) error {
	return exec.Command("open", "-a", "Safari", s.URI).Start()
}

// 末尾に ".json" を追加します
func addJSONSuffix(path string) string {
	if !strings.HasSuffix(path, ".json") {
		path = path + ".json"
	}
	return path
}
