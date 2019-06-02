package model

import (
	"fmt"
	"strings"
)

type SoiCup struct {
	Sois []Soi `json:"sois"`
}

func (s *SoiCup) Contains(name string) bool {
	for _, v := range s.Sois {
		if v.Name == name {
			return true
		}
	}
	return false
}

type Soi struct {
	Name string   `json:"name"`
	Uri  string   `json:"link"`
	Tags []string `json:"tags"`
}

func (s Soi) String() string {
	return fmt.Sprintf("%-15s %-45s %v \n", s.Name, s.Uri, s.Tags)
}

func FilterByFunc(sois []Soi, f func(Soi) bool) []Soi {
	var result []Soi
	for _, v := range sois {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func FilterByName(sois []Soi, name string) []Soi {
	byName := func(s Soi) bool {
		return s.Name == name
	}
	return FilterByFunc(sois, byName)
}

func FilterByExcludeName(sois []Soi, name string) []Soi {
	byExcludeName := func(s Soi) bool {
		return s.Name != name
	}
	return FilterByFunc(sois, byExcludeName)
}

func FilterByNamePart(sois []Soi, namePart string) []Soi {
	byNamePart := func(s Soi) bool {
		return strings.Contains(s.Name, namePart)
	}
	return FilterByFunc(sois, byNamePart)
}

func FilterByTags(sois []Soi, tags []string) []Soi {
	byTags := func(s Soi) bool {
		for _, t := range tags {
			findTag := false
			for _, vt := range s.Tags {
				if t == vt {
					findTag = true
				}
			}
			if !findTag {
				return false
			}
		}
		return true
	}
	return FilterByFunc(sois, byTags)
}
