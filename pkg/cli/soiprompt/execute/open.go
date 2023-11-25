package execute

import (
	"errors"
	"flag"
	"runtime"
	"soi-go/pkg/cli/opener"
	"soi-go/pkg/cli/opener/linux"
	"soi-go/pkg/cli/opener/macos"
	"soi-go/pkg/cli/opener/windows"
	"strings"
	"time"

	"golang.org/x/net/context"

	"soi-go/pkg/cli/soiprompt/view"
	"soi-go/pkg/model"
)

// open は指定されたSoiを元にブラウザを開きます
func (e *Executor) open(in string) error {
	flags := flag.NewFlagSet("open", flag.PanicOnError)
	chrome := flags.Bool("c", false, "use chrome")
	firefox := flags.Bool("f", false, "use firefox")
	safari := flags.Bool("s", false, "use safari")
	edge := flags.Bool("e", false, "use edge")
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

	opn, ok := getOpener(runtime.GOOS)
	if !ok {
		return errors.New("unsupported os :" + runtime.GOOS)
	}

	if *chrome {
		return opn.OpenChrome(s, *private)
	}
	if *firefox {
		return opn.OpenFirefox(s, *private)
	}
	if *safari {
		return opn.OpenSafari(s, *private)
	}
	if *edge {
		return opn.OpenEdge(s, *private)
	}

	defB := strings.ToLower(e.Conf.DefaultBrowser)
	switch defB {
	case "chrome":
		return opn.OpenChrome(s, *private)
	case "firefox":
		return opn.OpenFirefox(s, *private)
	case "safari":
		return opn.OpenSafari(s, *private)
	case "edge":
		return opn.OpenEdge(s, *private)
	default:
		return opn.OpenChrome(s, *private)
	}
}

func getOpener(os string) (opener.Opener, bool) {
	switch runtime.GOOS {
	case "darwin":
		return macos.NewOpener(), true
	case "windows":
		return windows.NewOpener(), true
	case "linux":
		return linux.NewLinuxOpener(), true
	}
	return nil, false
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
