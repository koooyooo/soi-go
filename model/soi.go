package model

import "strings"

type Sois struct {
	Sois []Soi `json:"sois"`
}

func (s *Sois) Add(soi Soi) {
	s.Sois = append(s.Sois, soi)
}

func (s *Sois) Remove(name string) {
	var newSois []Soi
	for _, v := range s.Sois {
		if v.Name == name {
			continue
		}
		newSois = append(newSois, v)
	}
	s.Sois = newSois
}

func (s *Sois) Contains(name string) bool {
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

func FilterByFunc(sois []Soi, f func(Soi) bool) []Soi {
	var result []Soi
	for _, v := range sois {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
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

func FilterByNamePart(sois []Soi, namePart string) []Soi {
	byNamePart := func(s Soi) bool {
		return strings.Contains(s.Name, namePart)
	}
	return FilterByFunc(sois, byNamePart)
}
