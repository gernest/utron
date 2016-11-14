package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestDefaultLogger(t *testing.T) {

	logData := []string{
		"INFO", "SUCC", "ERR", "WARN",
	}
	buf := &bytes.Buffer{}
	l := NewDefaultLogger(buf)
	msg := "hello"
	l.Info(msg)
	l.Errors(msg)
	l.Warn(msg)
	l.Success(msg)

	out := buf.String()
	for _, v := range logData {
		if !strings.Contains(out, v) {
			t.Errorf("expected %s to contain %s", out, v)
		}
	}
}
