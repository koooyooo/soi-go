package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Normal Test for parseTitle
func TestParseTitleByURL(t *testing.T) {
	var buff bytes.Buffer
	buff.WriteString(`<html><head><meta charSet="utf-8"/><title>MyTitle</title></head><body></body></html>`)
	title, ok, err := parseTitle(&buff)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "MyTitle", title)
}
