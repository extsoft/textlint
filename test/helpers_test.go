package main_test

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func buildBinary(t *testing.T) string {
	exe := filepath.Join(t.TempDir(), "prosefmt")
	if err := exec.Command("go", "build", "-o", exe, "..").Run(); err != nil {
		t.Fatalf("build: %v", err)
	}
	return exe
}
