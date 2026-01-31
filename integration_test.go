package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegration_CheckAndFix(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.txt")
	if err := os.WriteFile(bad, []byte("hello   \nworld\t\n\n"), 0644); err != nil {
		t.Fatal(err)
	}
	exe := filepath.Join(t.TempDir(), "textlint")
	if err := exec.Command("go", "build", "-o", exe, ".").Run(); err != nil {
		t.Fatalf("build: %v", err)
	}
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
	cmdFix := exec.Command(exe, "--fix", bad)
	cmdFix.Dir = dir
	if outFix, err := cmdFix.CombinedOutput(); err != nil {
		t.Fatalf("fix: %v\n%s", err, outFix)
	}
	cmdCheck2 := exec.Command(exe, "--check", bad)
	cmdCheck2.Dir = dir
	out2, _ := cmdCheck2.CombinedOutput()
	if cmdCheck2.ProcessState.ExitCode() != 0 {
		t.Errorf("expected exit 0 after fix, got %d\n%s", cmdCheck2.ProcessState.ExitCode(), out2)
	}
	if !strings.Contains(string(out2), "0 issue(s)") {
		t.Errorf("expected 0 issue(s) after fix, got %s", out2)
	}
}
