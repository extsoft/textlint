package log

import (
	"bytes"
	"testing"
)

func TestLogf_RespectsLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	SetOutput(buf)
	defer SetOutput(nil)
	SetLevel(Quiet)
	Logf(Verbose, "should not appear\n")
	if buf.Len() != 0 {
		t.Errorf("Quiet: expected no output for Verbose, got %q", buf.String())
	}
	buf.Reset()
	SetLevel(Verbose)
	Logf(Verbose, "visible\n")
	if !bytes.Contains(buf.Bytes(), []byte("visible")) {
		t.Errorf("Verbose: expected visible, got %q", buf.String())
	}
	buf.Reset()
	Logf(Debug, "debug should not appear\n")
	if buf.Len() != 0 {
		t.Errorf("Verbose: Debug messages should not appear, got %q", buf.String())
	}
	SetLevel(Debug)
	buf.Reset()
	Logf(Debug, "debug visible\n")
	if !bytes.Contains(buf.Bytes(), []byte("debug visible")) {
		t.Errorf("Debug: expected debug visible, got %q", buf.String())
	}
	SetLevel(Normal)
}
