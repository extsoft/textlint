package rules

type Issue struct {
	File    string
	Line    int
	Column  int
	RuleID  string
	Message string
}
