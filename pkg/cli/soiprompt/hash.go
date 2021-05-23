package soiprompt

import (
	"crypto/sha1"
	"fmt"
)

func hash(s string) string {
	sh := sha1.New()
	sh.Write([]byte(s))
	b := sh.Sum(nil)
	return fmt.Sprintf("%x", b)
}
