package soisort

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"

	"github.com/koooyooo/soi-go/pkg/model"
)

var yearsAgo100 = time.Now().Add(24 * 365 * -100 * time.Hour)

func Exec(sois []*model.SoiData, d prompt.Document) {
	sorters := []soiSorter{ViewSorter, MostRecentlyUsedSorter, MostRecentlyCreatedSorter}
	for _, s := range sorters {
		if strings.Contains(d.Text, fmt.Sprintf(" %s ", s.opt)) {
			sort.SliceStable(sois, s.sortFunc(sois))
		}
	}
}

type soiSorter struct {
	opt      string
	sortFunc func([]*model.SoiData) func(i, j int) bool
}

var ViewSorter = soiSorter{
	opt: "-v",
	sortFunc: func(sois []*model.SoiData) func(i, j int) bool {
		return func(i, j int) bool {
			return sois[i].NumViews > sois[j].NumViews
		}
	},
}

var MostRecentlyUsedSorter = soiSorter{
	opt: "-u",
	sortFunc: func(sois []*model.SoiData) func(i, j int) bool {
		return func(i, j int) bool {
			lastUsageTime := func(sois []*model.SoiData, idx int) time.Time {
				soi := sois[idx]
				usageLogs := soi.UsageLogs
				if len(usageLogs) == 0 {
					return yearsAgo100
				}
				return soi.UsageLogs[len(soi.UsageLogs)-1].UsedAt
			}
			return lastUsageTime(sois, i).After(lastUsageTime(sois, j))
		}
	},
}

var MostRecentlyCreatedSorter = soiSorter{
	opt: "-r",
	sortFunc: func(sois []*model.SoiData) func(i, j int) bool {
		return func(i, j int) bool {
			return sois[i].CreatedAt.After(sois[j].CreatedAt)
		}
	},
}
