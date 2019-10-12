package service

import (
	"github.com/koooyooo/soi-go/model"
)

type SoiService interface {
	Add(name, uri string, tags []string) (*model.Soi, error)
	Search(namepart string) ([]model.Soi, error)
	Get(name string) (*model.Soi, bool, error)
	Remove(name string) error
	Tag(name string, tags []string) (*model.Soi, bool, error)
}
