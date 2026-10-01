package pathutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/koooyooo/soi-go/pkg/common/file"
	"github.com/koooyooo/soi-go/pkg/model"
)

// ResolveUnderRoot は root 配下の相対パスを絶対パスに解決し、root 外への脱出を防ぎます。
func ResolveUnderRoot(root, rel string) (string, error) {
	rel = strings.TrimSpace(rel)
	rel = strings.TrimPrefix(rel, "./")
	if rel == "" || rel == "." || rel == "/" {
		return "", fmt.Errorf("refusing to operate on bucket root")
	}
	// filepath.Join は絶対パス要素があると直前を捨てるため、先頭の / を落とす
	rel = strings.TrimPrefix(rel, "/")

	cleanRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	target := filepath.Join(cleanRoot, filepath.FromSlash(rel))
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

// ResolveExistingUnderRootWithSois はファイル名解決に失敗したとき、ロード済み soi の
// Path/Name（list 表示と同じ論理パス）から実ファイルを探します。
func ResolveExistingUnderRootWithSois(root, rel string, sois []*model.SoiData) (string, error) {
	if abs, err := ResolveExistingUnderRoot(root, rel); err == nil {
		return abs, nil
	}
	want := normalizeLogical(rel)
	for _, s := range sois {
		logical := normalizeLogical(filepath.ToSlash(filepath.Join(s.Path, s.Name)))
		if logical != want && logical+".json" != want {
			continue
		}
		if abs, ok := AbsPathForSoi(root, s); ok {
			return abs, nil
		}
	}
	return "", fmt.Errorf("no file or dir found: %s", preferRel(rel))
}

// AbsPathForSoi はバケットルート上の実ファイルパスを返します。
func AbsPathForSoi(bucketRoot string, s *model.SoiData) (string, bool) {
	candidates := []string{
		filepath.Join(bucketRoot, filepath.FromSlash(s.Path), file.ToStorableName(s.Name)),
		filepath.Join(bucketRoot, filepath.FromSlash(s.Path), s.Name+".json"),
		filepath.Join(bucketRoot, filepath.FromSlash(s.Path), s.Name),
	}
	for _, c := range candidates {
		if file.Exists(c) {
			return c, true
		}
	}
	dir := filepath.Join(bucketRoot, filepath.FromSlash(s.Path))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", false
	}
	want := strings.TrimSuffix(file.ToStorableName(s.Name), ".json")
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		base := strings.TrimSuffix(e.Name(), ".json")
		if base == want || base == s.Name {
			return filepath.Join(dir, e.Name()), true
		}
	}
	return "", false
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

func normalizeLogical(rel string) string {
	rel = strings.TrimSpace(rel)
	rel = strings.TrimPrefix(rel, "./")
	rel = strings.TrimPrefix(rel, "/")
	rel = strings.TrimSuffix(rel, "/")
	rel = strings.TrimSuffix(rel, ".json")
	return filepath.ToSlash(rel)
}

func existingCandidates(rel string) []string {
	rel = strings.TrimSpace(rel)
	rel = strings.TrimPrefix(rel, "./")
	if rel == "" {
		return nil
	}
	var out []string
	add := func(s string) {
		if s == "" {
			return
		}
		s = strings.TrimPrefix(s, "/")
		for _, x := range out {
			if x == s {
				return
			}
		}
		out = append(out, s)
	}

	add(rel)
	add(strings.TrimSuffix(rel, "/"))

	trimmed := strings.TrimSuffix(strings.TrimPrefix(rel, "/"), "/")
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
	rel = strings.TrimPrefix(rel, "./")
	rel = strings.TrimPrefix(rel, "/")
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
