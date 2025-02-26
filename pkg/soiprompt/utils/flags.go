package utils

import "strings"

// StringArray is for flag package
//
// see: flags_test.go for example
//
//	var tags StringArray
//
//	func main() {
//	    flag.Var(&tags, "t", "tag of the uri")
//	    fmt.Printf("tags=%v", tags)
//	}
type StringArray []string

func (s *StringArray) String() string {
	return strings.Join(*s, " ")
}

func (s *StringArray) Set(v string) error {
	*s = append(*s, v)
	return nil
}
