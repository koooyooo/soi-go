package hash

import (
	"crypto/sha256"
	"fmt"
)

func Hash(s string) string {
	b := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", b)
}
