package model

import (
	"fmt"
	"strings"
)

type SoiCup struct {
	Sois []Soi `json:"sois"`
}

func (cup *SoiCup) Contains(name string) bool {
	for _, v := range cup.Sois {
		if v.Name == name {
			return true
		}
	}
	return false
}

func (cup *SoiCup) TagSet(namePart string) []string {
	var tags []string
	for _, s := range cup.Sois {
		for _, t := range s.Tags {
			found := false
			for _, a := range tags {
				if t == a {
					found = true
				}
			}
			if !found && t != "" && strings.Contains(t, namePart) {
				tags = append(tags, t)
			}
		}
	}
	return tags
}

type Soi struct {
	Name string   `json:"name"`
	Uri  string   `json:"link"`
	Tags []string `json:"tags"`
}

func (s Soi) String() string {
	return fmt.Sprintf("%-15s %-45s %v", s.Name, s.Uri, s.Tags)
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
	return filterByTags(sois, tags, func(soiTag, specifiedTag string) bool {
		return soiTag == specifiedTag
	})
}

func FilterByTagsPartial(sois []Soi, tags []string) []Soi {
	return filterByTags(sois, tags, func(soiTag, specifiedTag string) bool {
		return strings.Contains(soiTag, specifiedTag)
	})
}

func filterByTags(sois []Soi, tags []string, matcher func(string, string) bool) []Soi {
	byTags := func(s Soi) bool {
		for _, soiTag := range s.Tags {
			findTag := false
			for _, t := range tags {
				if matcher(soiTag, t) {
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
