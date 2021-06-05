package common

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func Hash(s string) string {
	sh := sha1.New()
	io.WriteString(sh, s)
	return fmt.Sprintf("%x", sh.Sum(nil))
}
