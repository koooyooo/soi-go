package pathutil

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/koooyooo/soi-go/pkg/common/file"
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

// ResolveExistingUnderRoot は list / 補完で出る論理パス（拡張子なし・空白あり）も含め、
// 実在するファイルまたはディレクトリへ解決します。
func ResolveExistingUnderRoot(root, rel string) (string, error) {
	var lastResolveErr error
	for _, candidate := range existingCandidates(rel) {
		abs, err := ResolveUnderRoot(root, candidate)
		if err != nil {
			lastResolveErr = err
			continue
		}
		if file.Exists(abs) {
			return abs, nil
		}
	}
	if lastResolveErr != nil {
		return "", lastResolveErr
	}
	return "", fmt.Errorf("no file or dir found: %s", preferRel(rel))
}

// ResolveDestUnderRoot は移動先など「まだ存在しなくてよい」パスを正規化して解決します。
// 末尾が / ならディレクトリ、それ以外は soi ファイル（.json）として扱います。
func ResolveDestUnderRoot(root, rel string) (string, error) {
	rel = strings.TrimSpace(rel)
	if strings.HasSuffix(rel, "/") {
		return ResolveUnderRoot(root, rel)
	}
	return ResolveUnderRoot(root, preferRel(rel))
}

func existingCandidates(rel string) []string {
	rel = strings.TrimSpace(rel)
	if rel == "" {
		return nil
	}
	var out []string
	add := func(s string) {
		if s == "" {
			return
		}
		for _, x := range out {
			if x == s {
				return
			}
		}
		out = append(out, s)
	}

	add(rel)
	add(strings.TrimSuffix(rel, "/"))

	trimmed := strings.TrimSuffix(rel, "/")
	dir := filepath.ToSlash(filepath.Dir(trimmed))
	base := filepath.Base(trimmed)
	if dir == "." {
		dir = ""
	}
	storable := file.ToStorableName(base)
	if dir == "" {
		add(storable)
	} else {
		add(dir + "/" + storable)
	}
	if !strings.HasSuffix(base, ".json") {
		withJSON := base + ".json"
		if dir == "" {
			add(withJSON)
		} else {
			add(dir + "/" + withJSON)
		}
	}
	return out
}

func preferRel(rel string) string {
	rel = strings.TrimSpace(rel)
	if rel == "" {
		return rel
	}
	if strings.HasSuffix(rel, "/") {
		return strings.TrimSpace(rel)
	}
	dir := filepath.ToSlash(filepath.Dir(rel))
	base := filepath.Base(rel)
	if dir == "." {
		dir = ""
	}
	name := file.ToStorableName(base)
	if dir == "" {
		return name
	}
	return dir + "/" + name
}
