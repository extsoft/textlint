package fix

import (
	"os"
	"path/filepath"
	"testing"
	"prosefmt/internal/rules"
)

func TestApply(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.txt")
	content := []byte("hello   \nworld\t\t\n\n\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}
	if err := Apply(path); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte("hello\nworld\n")
	if string(after) != string(expected) {
		t.Errorf("expected %q, got %q", expected, after)
	}
	issues, err := rules.CheckFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 0 {
		t.Errorf("fixed file should have no issues, got %v", issues)
	}
}
