package rules

const (
	TL001ID    = "TL001"
	TL001NoEnd = "file must end with exactly one newline"
	TL001Multi = "file must end with exactly one newline (multiple newlines at end)"
)

func CheckTL001(file string, content []byte) []Issue {
	var issues []Issue
	if len(content) == 0 {
		return issues
	}
	le := detectLineEnding(content)
	if le == lineEndCRLF {
		if len(content) < 2 || content[len(content)-2] != '\r' || content[len(content)-1] != '\n' {
			issues = append(issues, Issue{File: file, Line: 1, Column: 1, RuleID: TL001ID, Message: TL001NoEnd})
			return issues
		}
		i := len(content) - 2
		for i >= 2 && content[i-2] == '\r' && content[i-1] == '\n' {
			i -= 2
		}
		if i != len(content)-2 {
			issues = append(issues, Issue{File: file, Line: 1, Column: 1, RuleID: TL001ID, Message: TL001Multi})
		}
		return issues
	}
	if content[len(content)-1] != '\n' {
		issues = append(issues, Issue{File: file, Line: 1, Column: 1, RuleID: TL001ID, Message: TL001NoEnd})
		return issues
	}
	i := len(content) - 1
	for i > 0 && content[i-1] == '\n' {
		i--
	}
	if i != len(content)-1 {
		issues = append(issues, Issue{File: file, Line: 1, Column: 1, RuleID: TL001ID, Message: TL001Multi})
	}
	return issues
}

func FixTL001(content []byte) []byte {
	if len(content) == 0 {
		return content
	}
	le := detectLineEnding(content)
	if le == lineEndCRLF {
		end := len(content)
		for end >= 2 && content[end-2] == '\r' && content[end-1] == '\n' {
			end -= 2
		}
		return append(content[:end], []byte(lineEndCRLF)...)
	}
	end := len(content)
	for end > 0 && content[end-1] == '\n' {
		end--
	}
	return append(content[:end], '\n')
}
