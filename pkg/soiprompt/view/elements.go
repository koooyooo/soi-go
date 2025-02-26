package view

import "strings"

var SingleOptions = []SingleOption{
	{"c", "use chrome"},
	{"f", "use firefox"},
	{"s", "use safari"},
	{"v", "order by view num"},
	{"u", "order by used date"},
	{"r", "order by register date"},
}

type SingleOption struct {
	Name  string
	Usage string
}

func isHash(s string) bool {
	return strings.HasPrefix(s, "|") && strings.HasSuffix(s, "|")
}

func isOption(s string) bool {
	return strings.HasPrefix(s, "-")
}
