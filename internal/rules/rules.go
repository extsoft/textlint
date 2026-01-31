package rules

import "os"

func Check(file string, content []byte) []Issue {
	var issues []Issue
	issues = append(issues, CheckTL010(file, content)...)
	issues = append(issues, CheckTL001(file, content)...)
	return issues
}

func Fix(content []byte) []byte {
	out := FixTL010(content)
	out = FixTL001(out)
	return out
}

func CheckFile(path string) ([]Issue, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Check(path, content), nil
}
