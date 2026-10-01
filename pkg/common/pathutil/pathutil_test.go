package pathutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/koooyooo/soi-go/pkg/model"
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

func TestResolveExistingUnderRootLeadingSlash(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "portal")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	filePath := filepath.Join(dir, "google.json")
	if err := os.WriteFile(filePath, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := ResolveExistingUnderRoot(root, "/portal/google")
	if err != nil {
		t.Fatalf("leading slash should still resolve under root: %v", err)
	}
	if got != filePath {
		t.Fatalf("got %s want %s", got, filePath)
	}
}

func TestResolveExistingUnderRootWithSois(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "himeno-kato")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	filePath := filepath.Join(dir, "ryoke-f.json")
	if err := os.WriteFile(filePath, []byte(`{"name":"ryoke-f","path":"hino-kato"}`), 0600); err != nil {
		t.Fatal(err)
	}
	sois := []*model.SoiData{{
		Name: "ryoke-f",
		Path: "himeno-kato",
		Hash: "abc",
	}}
	got, err := ResolveExistingUnderRootWithSois(root, "himeno-kato/ryoke-f", sois)
	if err != nil {
		t.Fatal(err)
	}
	if got != filePath {
		t.Fatalf("got %s want %s", got, filePath)
	}
}
