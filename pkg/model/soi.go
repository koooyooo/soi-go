package model

import (
	"path/filepath"
	"time"
)

type SoiData struct {
	ID     int
	Name   string   `json:"name"`  // ファイル名
	Title  string   `json:"title"` // <title>属性
	Hash   string   `json:"hash"`
	Path   string   `json:"path"` // (<basepath{soiroot + bucket}> + <path> + <name>.json)
	URI    string   `json:"uri"`
	Tags   []string `json:"tags"`
	KVTags []KVTag  `json:"kv_tags"`
	Rate   float32  `json:"rate"`

	OGTitle       string    `json:"og_title"`
	OGURL         string    `json:"og_url"`
	OGType        string    `json:"og_type"`
	OGDescription string    `json:"og_description"`
	OGSiteName    string    `json:"og_site_name"`
	OGImages      []OGImage `json:"og_images"`

	NumViews      int        `json:"num_views"`     // ページを開いた回数
	NumReads      float32    `json:"num_reads"`     // ページを実際に読んだ回数(2回半なら0.5)
	Comprehension int        `json:"comprehension"` // 理解度 1 - 100
	CreatedAt     time.Time  `json:"created_at"`
	UsageLogs     []UsageLog `json:"usage_log"`
}

func (s *SoiData) FilePath(basePath string) string {
	return filepath.Join(basePath, s.Path, s.Name+".json")
}

type OGImage struct {
	URL       string `json:"url"`
	SecureURL string `json:"secure_url"`
	Type      string `json:"type"` // Content-Type
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Alt       string `json:"alt"`
}

type KVTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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

func GetUsageType(s string) (*UsageType, bool) {
	switch s {
	case string(UsageTypeUpdate):
		return &UsageTypeUpdate, true
	case string(UsageTypeOpen):
		return &UsageTypeOpen, true
	default:
		return nil, false
	}
}
