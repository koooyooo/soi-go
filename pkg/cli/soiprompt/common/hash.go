package common

import (
	"crypto/sha1"
	"fmt"
)

func Hash(s string) string {
	sh := sha1.New()
	sh.Write([]byte(s))
	b := sh.Sum(nil)
	return fmt.Sprintf("%x", b)
}
