package pathutil

import (
	"path/filepath"
	"testing"
)

func TestResolveUnderRoot(t *testing.T) {
	root := t.TempDir()

	got, err := ResolveUnderRoot(root, "a/b.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(root, "a", "b.json")
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}

	if _, err := ResolveUnderRoot(root, "../outside"); err == nil {
		t.Fatal("expected escape error")
	}
	if _, err := ResolveUnderRoot(root, ""); err == nil {
		t.Fatal("expected empty path error")
	}
}
