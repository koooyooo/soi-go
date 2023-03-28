package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
)

func Sha1(s string) (string, error) {
	sh := sha1.New()
	if _, err := io.WriteString(sh, s); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sh.Sum(nil)), nil
}

func Sha256(s string) string {
	b := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", b)
}
