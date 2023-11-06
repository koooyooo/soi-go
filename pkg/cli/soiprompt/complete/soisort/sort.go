package soisort

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"

	"soi-go/pkg/model"
)

var yearsAgo100 = time.Now().Add(24 * 365 * -100 * time.Hour)

func Exec(sois []*model.SoiData, d prompt.Document) {
	sorters := []soiSorter{NumViewSorter, AddDaySorter, ViewDaySorter}
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

var NumViewSorter = soiSorter{
	opt: "-n",
	sortFunc: func(sois []*model.SoiData) func(i, j int) bool {
		return func(i, j int) bool {
			return sois[i].NumViews > sois[j].NumViews
		}
	},
}

var AddDaySorter = soiSorter{
	opt: "-a",
	sortFunc: func(sois []*model.SoiData) func(i, j int) bool {
		return func(i, j int) bool {
			return sois[i].CreatedAt.After(sois[j].CreatedAt)
		}
	},
}

var ViewDaySorter = soiSorter{
	opt: "-v",
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
