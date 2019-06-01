package model

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
