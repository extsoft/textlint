package rules

const (
	lineEndLF   = "\n"
	lineEndCRLF = "\r\n"
)

func detectLineEnding(content []byte) string {
	for i := 0; i < len(content); i++ {
		if content[i] == '\r' && i+1 < len(content) && content[i+1] == '\n' {
			return lineEndCRLF
		}
		if content[i] == '\n' {
			return lineEndLF
		}
	}
	return lineEndLF
}

func splitLines(content []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i := 0; i <= len(content); i++ {
		if i == len(content) {
			lines = append(lines, content[start:i])
			break
		}
		if content[i] == '\n' {
			lines = append(lines, content[start:i+1])
			start = i + 1
		} else if content[i] == '\r' && i+1 < len(content) && content[i+1] == '\n' {
			lines = append(lines, content[start:i+2])
			start = i + 2
			i++
		}
	}
	return lines
}
