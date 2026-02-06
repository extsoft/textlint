package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegration_Write_FixesFiles(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.txt")
	originalContent := []byte("hello   \nworld\t\n\n")
	if err := os.WriteFile(bad, originalContent, 0644); err != nil {
		t.Fatal(err)
	}
	exe := buildBinary(t)
	cmdFix := exec.Command(exe, "--write", bad)
	cmdFix.Dir = dir
	outFix, err := cmdFix.CombinedOutput()
	if err != nil {
		t.Fatalf("write: %v\n%s", err, outFix)
	}
	after, err := os.ReadFile(bad)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte("hello\nworld\n")
	if string(after) != string(expected) {
		t.Errorf("expected file content %q, got %q", expected, after)
	}
	if !strings.Contains(string(outFix), "Wrote") {
		t.Errorf("expected output to contain 'Wrote', got %s", outFix)
	}
}

func TestIntegration_Write_ExitCodeZero(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.txt")
	if err := os.WriteFile(bad, []byte("x  \n"), 0644); err != nil {
		t.Fatal(err)
	}
	exe := buildBinary(t)
	cmd := exec.Command(exe, "--write", bad)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Errorf("expected exit 0 for write, got error: %v", err)
	}
	if cmd.ProcessState.ExitCode() != 0 {
		t.Errorf("expected exit 0 for write, got %d", cmd.ProcessState.ExitCode())
	}
}
