package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
)

func Sha1(s string) string {
	sh := sha1.New()
	io.WriteString(sh, s)
	return fmt.Sprintf("%x", sh.Sum(nil))
}

func Sha256(s string) string {
	b := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", b)
}
