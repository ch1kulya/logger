package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogFunctions(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(LevelDebug)
	defer SetLevel(LevelDebug)

	Debug("debug message %d", 1)
	Info("info message %s", "test")
	Warn("warn message")
	Error("error message")

	out := buf.String()
	t.Log("\n" + out)

	if !strings.Contains(out, "DEBUG") {
		t.Error("expected DEBUG in output")
	}
	if !strings.Contains(out, "INFO") {
		t.Error("expected INFO in output")
	}
	if !strings.Contains(out, "WARN") {
		t.Error("expected WARN in output")
	}
	if !strings.Contains(out, "ERR") {
		t.Error("expected ERR in output")
	}
}

func TestLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	SetLevel(LevelWarn)
	defer SetLevel(LevelDebug)

	Debug("should not appear")
	Info("should not appear")
	Warn("should appear")
	Error("should appear")

	out := buf.String()
	t.Log("\n" + out)

	if strings.Contains(out, "DEBUG") {
		t.Error("DEBUG should be filtered")
	}
	if strings.Contains(out, "INFO") {
		t.Error("INFO should be filtered")
	}
	if !strings.Contains(out, "WARN") {
		t.Error("WARN should appear")
	}
	if !strings.Contains(out, "ERR") {
		t.Error("ERR should appear")
	}
}

func TestFatal(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)

	var exitCode int
	originalExit := exitFunc
	exitFunc = func(code int) { exitCode = code }
	defer func() { exitFunc = originalExit }()

	Fatal("fatal error %s", "test")

	out := buf.String()
	t.Log("\n" + out)

	if !strings.Contains(out, "FATAL") {
		t.Error("expected FATAL in output")
	}
	if !strings.Contains(out, "fatal error test") {
		t.Error("expected message in output")
	}
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
}
