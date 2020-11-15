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
		Name           string     `json:"name"`
		Title          string     `json:"title"`
		URI            string     `json:"uri"`
		Tags           []string   `json:"tags"`
		Rate           float32    `json:"rate"`
		NumOfTimesRead float32    `json:"num_of_times_read"`
		CreatedAt      time.Time  `json:"created_at"`
		UsageLogs      []UsageLog `json:"usage_log"`
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
