package rules

const (
	TL010ID     = "TL010"
	TL010Msg    = "no trailing spaces at end of line"
)

func CheckTL010(file string, content []byte) []Issue {
	var issues []Issue
	lines := splitLines(content)
	for lineNum, raw := range lines {
		line := lineNum + 1
		contentPart, _ := stripLineEnding(raw)
		trailingStart := trailingSpaceStart(contentPart)
		if trailingStart >= 0 {
			issues = append(issues, Issue{
				File:    file,
				Line:    line,
				Column:  trailingStart + 1,
				RuleID:  TL010ID,
				Message: TL010Msg,
			})
		}
	}
	return issues
}

func stripLineEnding(raw []byte) (content []byte, ending []byte) {
	if len(raw) >= 2 && raw[len(raw)-2] == '\r' && raw[len(raw)-1] == '\n' {
		return raw[:len(raw)-2], raw[len(raw)-2:]
	}
	if len(raw) >= 1 && raw[len(raw)-1] == '\n' {
		return raw[:len(raw)-1], raw[len(raw)-1:]
	}
	return raw, nil
}

func trailingSpaceStart(content []byte) int {
	i := len(content)
	for i > 0 && (content[i-1] == ' ' || content[i-1] == '\t') {
		i--
	}
	if i < len(content) {
		return i
	}
	return -1
}

func FixTL010(content []byte) []byte {
	lines := splitLines(content)
	var out []byte
	for _, raw := range lines {
		contentPart, ending := stripLineEnding(raw)
		trimmed := trimTrailingSpaces(contentPart)
		out = append(out, trimmed...)
		out = append(out, ending...)
	}
	return out
}

func trimTrailingSpaces(b []byte) []byte {
	i := len(b)
	for i > 0 && (b[i-1] == ' ' || b[i-1] == '\t') {
		i--
	}
	return b[:i]
}
