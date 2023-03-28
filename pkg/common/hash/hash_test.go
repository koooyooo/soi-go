package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha1(t *testing.T) {
	s, err := Sha1("Hello")
	assert.NoError(t, err)
	assert.Equal(t, "f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0", s)
}

func TestSha256(t *testing.T) {
	s := Sha256("hello world")
	assert.Equal(t, "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", s)
}
