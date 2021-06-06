package soi

import "time"

type UsageLog struct {
	Type   UsageType `json:"type"`
	UsedAt time.Time `json:"used_at"`
}

type UsageType string

var (
	UsageTypeUpdate = UsageType("update")
	UsageTypeOpen   = UsageType("open")
)
