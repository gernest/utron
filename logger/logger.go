package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var logThis = NewDefaultLogger(os.Stdout)

// Logger is an interface for utron logger
type Logger interface {
	Info(v ...interface{})
	Errors(fv ...interface{})
	Warn(v ...interface{})
	Success(v ...interface{})
}

// DefaultLogger is the default logger
type DefaultLogger struct {
	*log.Logger
}

// NewDefaultLogger returns a default logger writing to out
func NewDefaultLogger(out io.Writer) Logger {
	d := &DefaultLogger{}
	d.Logger = log.New(out, "", log.LstdFlags)
	return d
}

// Info logs info messages
func (d *DefaultLogger) Info(v ...interface{}) {
	d.Println(fmt.Sprintf(">>INFO>> %s", fmt.Sprint(v...)))
}

// Errors log error messages
func (d *DefaultLogger) Errors(v ...interface{}) {
	d.Println(fmt.Sprintf(">>ERR>> %s", fmt.Sprint(v...)))
}

// Warn logs warning messages
func (d *DefaultLogger) Warn(v ...interface{}) {
	d.Println(fmt.Sprintf(">>WARN>> %s", fmt.Sprint(v...)))
}

// Success logs success messages
func (d *DefaultLogger) Success(v ...interface{}) {
	d.Println(fmt.Sprintf(">>SUCC>> %s", fmt.Sprint(v...)))
}
