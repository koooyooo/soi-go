package cache

import "soi-go/pkg/model"

type Cache struct {
	ListSoiCache []*model.SoiData
	DigPathCache []string
}

func (c *Cache) Clear() {
	c.ListSoiCache = nil
	c.DigPathCache = nil
}
