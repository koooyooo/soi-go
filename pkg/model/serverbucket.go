package model

import (
	"encoding/json"

	"soi-go/pkg/common/hash"
)

type ServerBucket struct {
	Sois []*SoiData `json:"sois"`
}

func (s *ServerBucket) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (s *ServerBucket) Hash() (string, error) {
	var hashes string
	for _, s := range s.Sois {
		hashes += s.Hash
	}
	return hash.Sha1(hashes)
}

// TODO implement
type Client struct {
	WorkingDir string
}
