package soi

import (
	"time"

	"github.com/koooyooo/soi-go/pkg/common/hash"
)

type SoiData struct {
	Name          string     `json:"name"`
	Title         string     `json:"title"`
	Hash          string     `json:"hash"`
	URI           string     `json:"uri"`
	Tags          []string   `json:"tags"`
	Rate          float32    `json:"rate"`
	NumViews      int32      `json:"num_views"`     // ページを開いた回数
	NumReads      float32    `json:"num_reads"`     // ページを実際に読んだ回数(2回半なら0.5)
	Comprehension int32      `json:"comprehension"` // 理解度 1 - 100
	CreatedAt     time.Time  `json:"created_at"`
	UsageLogs     []UsageLog `json:"usage_log"`
}

func (s SoiData) GetHash() string {
	if s.Hash == "" {
		return hash.Sha1(s.URI)
	}
	return s.Hash
}

type UsageLog struct {
	Type   UsageType `json:"type"`
	UsedAt time.Time `json:"used_at"`
}

type UsageType string

var (
	UsageTypeUpdate = UsageType("update")
	UsageTypeOpen   = UsageType("open")
)
