package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegration_Check_FindsIssues(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.txt")
	if err := os.WriteFile(bad, []byte("hello   \nworld\t\n\n"), 0644); err != nil {
		t.Fatal(err)
	}
	exe := buildBinary(t)
	cmd := exec.Command(exe, "--check", bad)
	cmd.Dir = dir
	out, _ := cmd.CombinedOutput()
	if cmd.ProcessState.ExitCode() != 1 {
		t.Errorf("expected exit 1 when issues found, got %d\n%s", cmd.ProcessState.ExitCode(), out)
	}
	if !strings.Contains(string(out), "TL001") {
		t.Errorf("expected output to contain TL001, got %s", out)
	}
	if !strings.Contains(string(out), "TL010") {
		t.Errorf("expected output to contain TL010, got %s", out)
	}
}

func TestIntegration_Check_NoIssues(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "good.txt")
	if err := os.WriteFile(good, []byte("hello\nworld\n"), 0644); err != nil {
		t.Fatal(err)
	}
	exe := buildBinary(t)
	cmd := exec.Command(exe, "--check", good)
	cmd.Dir = dir
	out, _ := cmd.CombinedOutput()
	if cmd.ProcessState.ExitCode() != 0 {
		t.Errorf("expected exit 0 when no issues, got %d\n%s", cmd.ProcessState.ExitCode(), out)
	}
	if !strings.Contains(string(out), "0 issue(s)") {
		t.Errorf("expected output to contain '0 issue(s)', got %s", out)
	}
}
