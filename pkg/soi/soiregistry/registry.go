/*
Package soiregistry offers registry service.
*/
package soiregistry

import (
	"github.com/koooyooo/soi-go/pkg/soi"
)

type Registry interface {
	Register(s *soi.Soi) error
	FindByName(n string) ([]*soi.Soi, error)
	FindByNamePart(namePart string) ([]*soi.Soi, error)
	FindByTag(tag string) ([]*soi.Soi, error)
	Remove(name string) error
}

type RegistryType int

const (
	RegistryTypeLocal RegistryType = iota
)

func NewRegistry(t RegistryType) (Registry, bool) {
	switch t {
	case RegistryTypeLocal:
		return &singleRegistry{}, true
	}
	return nil, false
}
