package soi

import "time"

var (
	UsageTypeUpdated = UsageType("updated")
	UsageTypeUsed    = UsageType("used")
)

type (
	SoiData struct {
		Name      string     `json:"name"`
		Title     string     `json:"title"`
		URI       string     `json:"uri"`
		Tags      []string   `json:"tags"`
		Created   time.Time  `json:"created_at"`
		UsageLogs []UsageLog `json:"usage_log"`
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
