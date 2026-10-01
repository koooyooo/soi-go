package pathutil

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ResolveUnderRoot は root 配下の相対パスを絶対パスに解決し、root 外への脱出を防ぎます。
func ResolveUnderRoot(root, rel string) (string, error) {
	if rel == "" || rel == "." || rel == "/" {
		return "", fmt.Errorf("refusing to operate on bucket root")
	}
	cleanRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	target := filepath.Join(cleanRoot, rel)
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	sep := string(filepath.Separator)
	if absTarget != cleanRoot && !strings.HasPrefix(absTarget, cleanRoot+sep) {
		return "", fmt.Errorf("path escapes bucket root: %s", rel)
	}
	if absTarget == cleanRoot {
		return "", fmt.Errorf("refusing to operate on bucket root")
	}
	return absTarget, nil
}
