package view

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/koooyooo/soi-go/pkg/model"
)

type SoiLine struct {
	Hash     string
	NumViews int
	Path     string
	Tags     []string
}

func ParseLine4Hash(elms []string) (string, bool) {
	for _, e := range elms {
		if isHash(e) {
			return e[1 : len(e)-1], true
		}
	}
	return "", false
}

func ParseLine4Path(elms []string) (string, bool) {
	// dig であることを期待しオプションを除いた次の要素を取得
	for _, elm := range elms {
		if isOption(elm) {
			continue
		}
		return elm, true
	}
	return "", false
}

// コンソール表示をパースします。
func ParseLine(s string) (*SoiLine, error) {
	s = strings.Replace(s, " -v", "", -1)
	s = strings.Replace(s, " -u", "", -1)
	s = strings.Replace(s, " -r", "", -1)

	const spaceWidth = 1
	const bracketWidth = 1
	l := new(SoiLine)

	// indexes
	idxStNumView := strings.Index(s, "[")
	idxEdNumView := strings.Index(s, "]")

	idxStTags := strings.LastIndex(s, "[")
	idxEdTags := strings.LastIndex(s, "]")

	// Hash
	l.Hash = strings.TrimSpace(s[0 : idxStNumView-spaceWidth])

	// NumView
	numViewsStr := strings.TrimSpace(s[idxStNumView+bracketWidth : idxEdNumView])
	numViews, err := strconv.Atoi(numViewsStr)
	if err != nil {
		return nil, err
	}
	l.NumViews = numViews

	// Path
	l.Path = strings.TrimSpace(s[idxEdNumView+bracketWidth+spaceWidth : idxStTags-spaceWidth])

	// Tags
	tagStr := s[idxStTags+bracketWidth : idxEdTags]
	tags := strings.Split(tagStr, " ")
	if len(tags) == 1 && tags[0] == "" {
		tags = nil
	}
	l.Tags = tags

	return l, nil
}

func ToLine(s *model.SoiData, soisDir string) string {
	relPath := strings.TrimPrefix(s.Path, soisDir+"/")
	baseName := strings.TrimSuffix(relPath, ".json")
	return fmt.Sprintf("|%s| %s %s %s", s.Hash[0:8], metaString(s), baseName, tagsString(s))
}

// metaString はメタヘッダを作成します
func metaString(sd *model.SoiData) string {
	return fmt.Sprintf("[%3d]", sd.NumViews)
}

func tagsString(s *model.SoiData) string {
	var tags []string
	for _, t := range s.Tags {
		tags = append(tags, fmt.Sprintf("#%s", t))
	}
	for _, kvt := range s.KVTags {
		tags = append(tags, fmt.Sprintf("#%s=%s", kvt.Key, kvt.Value))
	}

	var buff strings.Builder
	buff.WriteString("[")
	buff.WriteString(strings.Join(tags, " "))
	buff.WriteString("]")

	return buff.String()
}
