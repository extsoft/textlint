package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"textlint/internal/rules"
)

func TestWrite_Compact(t *testing.T) {
	issues := []rules.Issue{
		{File: "a.txt", Line: 1, Column: 5, RuleID: "TL010", Message: "no trailing spaces"},
		{File: "a.txt", Line: 2, Column: 1, RuleID: "TL001", Message: "file must end with exactly one newline"},
	}
	var buf bytes.Buffer
	if err := Write(&buf, FormatCompact, issues, 10); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "a.txt:1:5: TL010: no trailing spaces") {
		t.Errorf("expected compact line for first issue, got %q", out)
	}
	if !strings.Contains(out, "10 file(s) scanned, 2 issue(s).") {
		t.Errorf("expected summary with scanned count, got %q", out)
	}
}

func TestWrite_Compact_TwoFiles(t *testing.T) {
	issues := []rules.Issue{
		{File: "a.txt", Line: 1, Column: 1, RuleID: "TL001", Message: "x"},
		{File: "b.txt", Line: 1, Column: 1, RuleID: "TL001", Message: "y"},
	}
	var buf bytes.Buffer
	if err := Write(&buf, FormatCompact, issues, 6); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "6 file(s) scanned, 2 issue(s).") {
		t.Errorf("expected 6 file(s) scanned, 2 issue(s). in summary, got %q", buf.String())
	}
}

func TestWrite_TAP_ZeroIssues(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, FormatTAP, nil, -1); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "1..1\n") || !strings.Contains(out, "ok 1") {
		t.Errorf("expected TAP 1..1 and ok 1, got %q", out)
	}
}

func TestWrite_TAP_WithIssues(t *testing.T) {
	issues := []rules.Issue{
		{File: "f", Line: 1, Column: 1, RuleID: "TL001", Message: "no newline"},
	}
	var buf bytes.Buffer
	if err := Write(&buf, FormatTAP, issues, -1); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "1..1\n") || !strings.Contains(out, "not ok 1") {
		t.Errorf("expected TAP 1..1 and not ok 1, got %q", out)
	}
}

func TestWrite_JSON(t *testing.T) {
	issues := []rules.Issue{
		{File: "a.txt", Line: 2, Column: 3, RuleID: "TL010", Message: "trailing"},
	}
	var buf bytes.Buffer
	if err := Write(&buf, FormatJSON, issues, 5); err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		Files   map[string][]struct {
			Line    int    `json:"line"`
			Column  int    `json:"column"`
			Rule    string `json:"rule"`
			Message string `json:"message"`
		} `json:"files"`
		Summary struct {
			Files         int  `json:"files"`
			Issues        int  `json:"issues"`
			FilesScanned  *int `json:"files_scanned"`
		} `json:"summary"`
	}
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded.Summary.Files != 1 || decoded.Summary.Issues != 1 {
		t.Errorf("expected summary files=1 issues=1, got %+v", decoded.Summary)
	}
	if decoded.Summary.FilesScanned == nil || *decoded.Summary.FilesScanned != 5 {
		t.Errorf("expected summary files_scanned=5, got %v", decoded.Summary.FilesScanned)
	}
	list, ok := decoded.Files["a.txt"]
	if !ok || len(list) != 1 || list[0].Line != 2 || list[0].Column != 3 || list[0].Rule != "TL010" {
		t.Errorf("expected files[a.txt] with one issue, got %+v", decoded.Files)
	}
}
