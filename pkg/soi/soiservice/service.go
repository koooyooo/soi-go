/*
Package soiservice offers service functions
*/
package soiservice

import (
	"github.com/koooyooo/soi-go/pkg/soi"
	"github.com/koooyooo/soi-go/pkg/soi/soiregistry"
)

// SoiService offers soi.go services
type SoiService interface {
	// Add adds soi.go.
	Add(*soi.Soi) (*soi.Soi, error)
	// Search searches sois by its namePart.
	Search(namePart string) ([]*soi.Soi, error)
	// SearchByTag searches sois by its tag.
	SearchByTag(tag string) ([]*soi.Soi, error)
	// Get finds a soi.go by its name.
	Get(name string) (*soi.Soi, bool, error)
	// Remove removes a soi.go by its name.
	Remove(name string) error
	// Tag tags a soi.go
	Tag(name string, tags []string) (*soi.Soi, bool, error)
}

type ServiceType int

const (
	ServiceTypePlain ServiceType = iota
)

func NewSoiService(st ServiceType, registry soiregistry.Registry) (SoiService, bool) {
	switch st {
	case ServiceTypePlain:
		return plainSoiService{
			Registry: registry,
		}, true
	}
	return nil, false
}
