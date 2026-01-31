package log

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Level int

const (
	Quiet Level = iota
	Normal
	Verbose
	Debug
)

var (
	mu            sync.RWMutex
	currentLevel Level = Normal
	output       io.Writer = os.Stderr
)

func SetLevel(l Level) {
	mu.Lock()
	defer mu.Unlock()
	currentLevel = l
}

func GetLevel() Level {
	mu.RLock()
	defer mu.RUnlock()
	return currentLevel
}

func SetOutput(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	if w == nil {
		output = os.Stderr
		return
	}
	output = w
}

func out() io.Writer {
	mu.RLock()
	defer mu.RUnlock()
	return output
}

func Logf(l Level, format string, args ...interface{}) {
	if GetLevel() < l {
		return
	}
	fmt.Fprintf(out(), format, args...)
}

func Log(l Level, msg string) {
	if GetLevel() < l {
		return
	}
	fmt.Fprint(out(), msg)
}
