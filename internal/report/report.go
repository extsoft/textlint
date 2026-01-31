package report

import (
	"encoding/json"
	"fmt"
	"io"
	"textlint/internal/rules"
)

type Format string

const (
	FormatCompact Format = "compact"
	FormatTAP     Format = "tap"
	FormatJSON    Format = "json"
)

func Write(w io.Writer, format Format, issues []rules.Issue, filesScanned int) error {
	switch format {
	case FormatCompact:
		return writeCompact(w, issues, filesScanned)
	case FormatTAP:
		return writeTAP(w, issues)
	case FormatJSON:
		return writeJSON(w, issues, filesScanned)
	default:
		return fmt.Errorf("unknown format %q", format)
	}
}

func writeCompact(w io.Writer, issues []rules.Issue, filesScanned int) error {
	for _, i := range issues {
		_, err := fmt.Fprintf(w, "%s:%d:%d: %s: %s\n", i.File, i.Line, i.Column, i.RuleID, i.Message)
		if err != nil {
			return err
		}
	}
	if filesScanned >= 0 {
		_, err := fmt.Fprintf(w, "%d file(s) scanned, %d issue(s).\n", filesScanned, len(issues))
		return err
	}
	files := fileSet(issues)
	_, err := fmt.Fprintf(w, "%d file(s), %d issue(s).\n", len(files), len(issues))
	return err
}

func fileSet(issues []rules.Issue) map[string]bool {
	m := make(map[string]bool)
	for _, i := range issues {
		m[i.File] = true
	}
	return m
}

func writeTAP(w io.Writer, issues []rules.Issue) error {
	n := len(issues)
	if n == 0 {
		_, err := fmt.Fprintf(w, "1..1\nok 1 - no issues\n")
		return err
	}
	if _, err := fmt.Fprintf(w, "1..%d\n", n); err != nil {
		return err
	}
	for i, iss := range issues {
		_, err := fmt.Fprintf(w, "not ok %d - %s:%d:%d %s %s\n", i+1, iss.File, iss.Line, iss.Column, iss.RuleID, iss.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, issues []rules.Issue, filesScanned int) error {
	files := make(map[string][]jsonIssue)
	for _, i := range issues {
		files[i.File] = append(files[i.File], jsonIssue{
			Line:    i.Line,
			Column:  i.Column,
			Rule:    i.RuleID,
			Message: i.Message,
		})
	}
	sum := jsonSummary{Files: len(files), Issues: len(issues)}
	if filesScanned >= 0 {
		n := filesScanned
		sum.FilesScanned = &n
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jsonReport{
		Files:   files,
		Summary: sum,
	})
}

type jsonIssue struct {
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

type jsonSummary struct {
	Files         int  `json:"files"`
	Issues        int  `json:"issues"`
	FilesScanned  *int `json:"files_scanned,omitempty"`
}

type jsonReport struct {
	Files   map[string][]jsonIssue `json:"files"`
	Summary jsonSummary            `json:"summary"`
}
