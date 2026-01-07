package execute

import (
	"errors"
	"flag"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/koooyooo/soi-go/pkg/opener"

	"golang.org/x/net/context"

	"github.com/koooyooo/soi-go/pkg/model"
	"github.com/koooyooo/soi-go/pkg/soiprompt/view"
)

// open は指定されたSoiを元にブラウザを開きます
func (e *Executor) open(in string) error {
	flags := flag.NewFlagSet("open", flag.PanicOnError)
	chrome := flags.Bool("c", false, "use chrome")
	firefox := flags.Bool("f", false, "use firefox")
	safari := flags.Bool("s", false, "use safari")
	edge := flags.Bool("e", false, "use edge")
	private := flags.Bool("p", false, "private mode")

	sortByNumViews := flags.Bool("n", false, "sort by num-views")
	sortByAddDay := flags.Bool("a", false, "sort by add-day")
	sortByViewDay := flags.Bool("v", false, "sort by view-day")

	if err := flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}

	ctx := context.Background()
	
	// ソート処理を適用
	sortedSois := e.Cache.ListSoiCache
	if *sortByNumViews || *sortByAddDay || *sortByViewDay {
		sortedSois = applySorting(e.Cache.ListSoiCache, *sortByNumViews, *sortByAddDay, *sortByViewDay)
	}
	
	s, err := findSoi(sortedSois, flags.Args())
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
	switch os {
	case "darwin":
		return opener.NewMacOpener(), true
	case "windows":
		return opener.NewWindowsOpener(), true
	case "linux":
		if isWSL() {
			return opener.NewWSLOpener(), true
		}
		return opener.NewLinuxOpener(), true
	}
	return nil, false
}

// isWSL checks if the current Linux environment is running under WSL
func isWSL() bool {
	content, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(content)), "microsoft")
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

// applySorting はソートオプションに基づいてSoiDataをソートします
func applySorting(sois []*model.SoiData, sortByNumViews, sortByAddDay, sortByViewDay bool) []*model.SoiData {
	// コピーを作成してソート
	sorted := make([]*model.SoiData, len(sois))
	copy(sorted, sois)
	
	if sortByNumViews {
		// 閲覧回数でソート（降順）
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i].NumViews < sorted[j].NumViews {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
	} else if sortByAddDay {
		// 追加日でソート（新しい順）
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i].CreatedAt.Before(sorted[j].CreatedAt) {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
	} else if sortByViewDay {
		// 最終閲覧日でソート（新しい順）
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				iLastView := getLastViewTime(sorted[i])
				jLastView := getLastViewTime(sorted[j])
				if iLastView.Before(jLastView) {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
	}
	
	return sorted
}

// getLastViewTime は最後の閲覧時間を取得します
func getLastViewTime(soi *model.SoiData) time.Time {
	if len(soi.UsageLogs) == 0 {
		return soi.CreatedAt
	}
	
	var lastView time.Time
	for _, log := range soi.UsageLogs {
		if log.Type == model.UsageTypeOpen && log.UsedAt.After(lastView) {
			lastView = log.UsedAt
		}
	}
	
	if lastView.IsZero() {
		return soi.CreatedAt
	}
	
	return lastView
}
