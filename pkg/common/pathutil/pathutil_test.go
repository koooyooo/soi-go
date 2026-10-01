package pathutil

import (
	"os"
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

func TestResolveExistingUnderRootLogicalName(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "portal")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	filePath := filepath.Join(dir, "google.json")
	if err := os.WriteFile(filePath, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := ResolveExistingUnderRoot(root, "portal/google")
	if err != nil {
		t.Fatalf("logical path should resolve: %v", err)
	}
	if got != filePath {
		t.Fatalf("got %s want %s", got, filePath)
	}

	got, err = ResolveExistingUnderRoot(root, "portal/google.json")
	if err != nil {
		t.Fatalf("explicit json should resolve: %v", err)
	}
	if got != filePath {
		t.Fatalf("got %s want %s", got, filePath)
	}
}

func TestResolveExistingUnderRootStorableName(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "notes")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	filePath := filepath.Join(dir, "Round_Trip.json")
	if err := os.WriteFile(filePath, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := ResolveExistingUnderRoot(root, "notes/Round Trip")
	if err != nil {
		t.Fatalf("spaced logical name should resolve: %v", err)
	}
	if got != filePath {
		t.Fatalf("got %s want %s", got, filePath)
	}
}

func TestResolveDestUnderRoot(t *testing.T) {
	root := t.TempDir()
	got, err := ResolveDestUnderRoot(root, "archive/example")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, "archive", "example.json")
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}
