package utils

import (
	"flag"
)

func ExampleStringArray_Set() {
	var tags StringArray

	flag.Var(&tags, "t", "tag of the uri")
}
