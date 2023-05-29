package execute

import (
	"errors"
	"flag"
	"github.com/koooyooo/soi-go/pkg/cli/opener"
	"golang.org/x/net/context"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/koooyooo/soi-go/pkg/cli/constant"
	"github.com/koooyooo/soi-go/pkg/cli/loader"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/view"
	"github.com/koooyooo/soi-go/pkg/model"
)

// open は指定されたSoiを元にブラウザを開きます
func (e *Executor) open(in string) error {
	flags := flag.NewFlagSet("open", flag.PanicOnError)
	chrome := flags.Bool("c", false, "use chrome")
	firefox := flags.Bool("f", false, "use firefox")
	safari := flags.Bool("s", false, "use safari")

	_ = flags.Bool("n", false, "sort by num-views")
	_ = flags.Bool("a", false, "sort by add-day")
	_ = flags.Bool("v", false, "sort by view-day")

	if err := flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}

	// Soiファイルを特定
	soisDir, err := e.Bucket.Path()
	if err != nil {
		return err
	}

	s, err := findSoi(soisDir, flags.Args())
	if err != nil {
		return err
	}

	// 閲覧履歴を追記
	s.NumViews++

	// 利用ログを記載
	s.UsageLogs = append(s.UsageLogs, model.UsageLog{
		Type:   model.UsageTypeOpen,
		UsedAt: time.Now(),
	})

	ctx := context.Background()
	if err := e.Service.Store(ctx, s); err != nil {
		return err
	}

	// 環境設定に応じてブラウザオープン
	if *chrome {
		return opener.OpenChrome(s)
	}
	if *firefox {
		return opener.OpenFirefox(s)
	}
	if *safari {
		return opener.OpenSafari(s)
	}
	defB := strings.ToLower(constant.EnvKeyDefaultBrowser.Get())
	if defB == "chrome" {
		return opener.OpenChrome(s)
	}
	if defB == "firefox" {
		return opener.OpenFirefox(s)
	}
	if defB == "safari" {
		return opener.OpenSafari(s)
	}
	return exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", s.URI).Start()
}

func findSoi(soisDir string, args []string) (*model.SoiData, error) {
	hash, findHash := view.ParseLine4Hash(args)
	if findHash {
		sois, err := loader.LoadSois(soisDir)
		if err != nil {
			return nil, err
		}
		for _, soi := range sois {
			if strings.HasPrefix(soi.Hash, hash) {
				return soi, nil
			}
		}
	}
	pathTail, findPath := view.ParseLine4Path(args)
	if findPath {
		p := path.Join(soisDir, pathTail+".json")
		soi, err := loader.LoadSoiData(p)
		if err != nil {
			return nil, err
		}
		return soi, nil
	}
	return nil, errors.New("no path found")
}
