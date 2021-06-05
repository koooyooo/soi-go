package utils

import (
	"strings"

	"github.com/c-bata/go-prompt"
)

func IsOptionWord(d prompt.Document) bool {
	return strings.HasPrefix(d.GetWordBeforeCursor(), "-")
}
