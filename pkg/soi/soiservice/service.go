/*
Package soiservice offers service functions
*/
package soiservice

import (
	"github.com/koooyooo/soi-go/pkg/soi"
	"github.com/koooyooo/soi-go/pkg/soi/soiregistry"
)

// SoiService offers soi services
type SoiService interface {
	// Add adds soi.
	Add(*soi.Soi) (*soi.Soi, error)
	// Search searches sois by its namePart.
	Search(namePart string) ([]*soi.Soi, error)
	// SearchByTag searches sois by its tag.
	SearchByTag(tag string) ([]*soi.Soi, error)
	// Get finds a soi by its name.
	Get(name string) (*soi.Soi, bool, error)
	// Remove removes a soi by its name.
	Remove(name string) error
	// Tag tags a soi
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
