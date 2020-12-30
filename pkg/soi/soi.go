package soi

import (
	"encoding/json"
	"time"
)

var (
	UsageTypeUpdate = UsageType("update")
	UsageTypeOpen   = UsageType("open")
)

type (
	SoiData struct {
		Name          string     `json:"name"`
		Title         string     `json:"title"`
		URI           string     `json:"uri"`
		Tags          []string   `json:"tags"`
		Rate          float32    `json:"rate"`
		NumViews      int32      `json:"num_views"`     // ページを開いた回数
		NumReads      float32    `json:"num_reads"`     // ページを実際に読んだ回数(2回半なら0.5)
		Comprehension int32      `json:"comprehension"` // 理解度 1 - 100
		CreatedAt     time.Time  `json:"created_at"`
		UsageLogs     []UsageLog `json:"usage_log"`
	}

	UsageType string

	UsageLog struct {
		Type   UsageType `json:"type"`
		UsedAt time.Time `json:"used_at"`
	}

	SoiVirtual struct {
		*SoiData
		Path string `json:"path"`
	}

	SoiVirtualBucket struct {
		Sois []*SoiVirtual `json:"sois"`
	}

	// TODO implement later
	Client struct {
		WorkingDir string
	}
)

func (s SoiVirtualBucket) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
