package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	s := Hash("Hello")
	assert.Equal(t, "f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0", s)
}
