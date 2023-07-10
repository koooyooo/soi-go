package execute

import (
	"errors"
	"flag"
	"github.com/koooyooo/soi-go/pkg/cli/opener"
	"golang.org/x/net/context"
	"os/exec"
	"strings"
	"time"

	"github.com/koooyooo/soi-go/pkg/cli/constant"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/view"
	"github.com/koooyooo/soi-go/pkg/model"
)

// open は指定されたSoiを元にブラウザを開きます
func (e *Executor) open(in string) error {
	flags := flag.NewFlagSet("open", flag.PanicOnError)
	chrome := flags.Bool("c", false, "use chrome")
	firefox := flags.Bool("f", false, "use firefox")
	safari := flags.Bool("s", false, "use safari")
	private := flags.Bool("p", false, "private mode")

	_ = flags.Bool("n", false, "sort by num-views")
	_ = flags.Bool("a", false, "sort by add-day")
	_ = flags.Bool("v", false, "sort by view-day")

	if err := flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}

	ctx := context.Background()
	s, err := findSoi(e.Cache.ListSoiCache, flags.Args())
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

	if err := e.Service.Store(ctx, s); err != nil {
		return err
	}

	if *chrome {
		return opener.OpenChrome(s, *private)
	}
	if *firefox {
		return opener.OpenFirefox(s, *private)
	}
	if *safari {
		return opener.OpenSafari(s, *private)
	}
	defB := strings.ToLower(constant.EnvKeyDefaultBrowser.Get())
	if defB == "chrome" {
		return opener.OpenChrome(s, *private)
	}
	if defB == "firefox" {
		return opener.OpenFirefox(s, *private)
	}
	if defB == "safari" {
		return opener.OpenSafari(s, *private)
	}
	return exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", s.URI).Start()
}

func findSoi(sois []*model.SoiData, args []string) (*model.SoiData, error) {
	hash, findHash := view.ParseLine4Hash(args)
	if findHash {
		for _, soi := range sois {
			if strings.HasPrefix(soi.Hash, hash) {
				return soi, nil
			}
		}
	}
	pathTail, findPath := view.ParseLine4Path(args)
	if findPath {
		for _, soi := range sois {
			if strings.Contains(soi.Path+"/"+soi.Name, pathTail) {
				return soi, nil
			}
		}
	}
	return nil, errors.New("no path found")
}
