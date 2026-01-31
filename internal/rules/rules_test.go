package rules

import (
	"bytes"
	"testing"
)

func TestCheckTL001_NoNewline(t *testing.T) {
	content := []byte("hello")
	issues := CheckTL001("f", content)
	if len(issues) != 1 || issues[0].RuleID != TL001ID || issues[0].Message != TL001NoEnd {
		t.Errorf("expected one TL001 no-end issue, got %v", issues)
	}
}

func TestCheckTL001_OneNewlineLF(t *testing.T) {
	content := []byte("hello\n")
	issues := CheckTL001("f", content)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %v", issues)
	}
}

func TestCheckTL001_OneNewlineCRLF(t *testing.T) {
	content := []byte("hello\r\n")
	issues := CheckTL001("f", content)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %v", issues)
	}
}

func TestCheckTL001_MultipleNewlinesLF(t *testing.T) {
	content := []byte("hello\n\n")
	issues := CheckTL001("f", content)
	if len(issues) != 1 || issues[0].Message != TL001Multi {
		t.Errorf("expected one TL001 multi issue, got %v", issues)
	}
}

func TestCheckTL001_MultipleNewlinesCRLF(t *testing.T) {
	content := []byte("a\r\n\r\n")
	issues := CheckTL001("f", content)
	if len(issues) != 1 || issues[0].Message != TL001Multi {
		t.Errorf("expected one TL001 multi issue, got %v", issues)
	}
}

func TestCheckTL001_EmptyFile(t *testing.T) {
	issues := CheckTL001("f", nil)
	if len(issues) != 1 || issues[0].Message != TL001NoEnd {
		t.Errorf("expected one TL001 no-end for empty file, got %v", issues)
	}
}

func TestFixTL001_NoNewline(t *testing.T) {
	content := []byte("hello")
	out := FixTL001(content)
	if !bytes.HasSuffix(out, []byte("\n")) || len(out) != 6 {
		t.Errorf("expected hello\\n, got %q", out)
	}
}

func TestFixTL001_OneNewline(t *testing.T) {
	content := []byte("hello\n")
	out := FixTL001(content)
	if !bytes.Equal(out, content) {
		t.Errorf("expected unchanged, got %q", out)
	}
}

func TestFixTL001_MultipleNewlines(t *testing.T) {
	content := []byte("hello\n\n\n")
	out := FixTL001(content)
	if !bytes.Equal(out, []byte("hello\n")) {
		t.Errorf("expected hello\\n, got %q", out)
	}
}

func TestFixTL001_CRLF(t *testing.T) {
	content := []byte("a\r\n\r\n")
	out := FixTL001(content)
	if !bytes.Equal(out, []byte("a\r\n")) {
		t.Errorf("expected a\\r\\n, got %q", out)
	}
}

func TestCheckTL010_NoTrailing(t *testing.T) {
	content := []byte("hello\n")
	issues := CheckTL010("f", content)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %v", issues)
	}
}

func TestCheckTL010_TrailingSpace(t *testing.T) {
	content := []byte("hello   \n")
	issues := CheckTL010("f", content)
	if len(issues) != 1 || issues[0].Line != 1 || issues[0].Column != 6 || issues[0].RuleID != TL010ID {
		t.Errorf("expected one TL010 at line 1 col 6, got %v", issues)
	}
}

func TestCheckTL010_TrailingTab(t *testing.T) {
	content := []byte("x\t\n")
	issues := CheckTL010("f", content)
	if len(issues) != 1 || issues[0].Column != 2 {
		t.Errorf("expected one TL010 col 2, got %v", issues)
	}
}

func TestCheckTL010_MultipleLines(t *testing.T) {
	content := []byte("a\nb  \nc\n")
	issues := CheckTL010("f", content)
	if len(issues) != 1 || issues[0].Line != 2 || issues[0].Column != 2 {
		t.Errorf("expected one TL010 at line 2 col 2, got %v", issues)
	}
}

func TestCheckTL010_CRLF(t *testing.T) {
	content := []byte("hi  \r\n")
	issues := CheckTL010("f", content)
	if len(issues) != 1 || issues[0].Column != 3 {
		t.Errorf("expected one TL010 col 3, got %v", issues)
	}
}

func TestFixTL010(t *testing.T) {
	content := []byte("a  \nb\t\t\n")
	out := FixTL010(content)
	expected := []byte("a\nb\n")
	if !bytes.Equal(out, expected) {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestFixTL010_PreservesCRLF(t *testing.T) {
	content := []byte("x  \r\n")
	out := FixTL010(content)
	if !bytes.Equal(out, []byte("x\r\n")) {
		t.Errorf("expected x\\r\\n, got %q", out)
	}
}

func TestCheck_Combined(t *testing.T) {
	content := []byte("a  \nb")
	issues := Check("f", content)
	if len(issues) < 2 {
		t.Errorf("expected at least TL010 and TL001, got %v", issues)
	}
}

func TestFix_Combined(t *testing.T) {
	content := []byte("x  \n\n\n")
	out := Fix(content)
	expected := []byte("x\n")
	if !bytes.Equal(out, expected) {
		t.Errorf("expected %q, got %q", expected, out)
	}
	issues := Check("f", out)
	if len(issues) != 0 {
		t.Errorf("fixed content should have no issues, got %v", issues)
	}
}
