package soi

import "encoding/json"

type (
	SoiWithPath struct {
		SoiData *SoiData
		Path    string
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
