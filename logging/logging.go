// Package logging defines the minimal interface that loggers must support to be used by pack.
package logging

import (
	"io"
	"io/ioutil"

	"github.com/YousefHaggyHeroku/pack/internal/style"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// Logger defines behavior required by a logging package used by pack libraries
type Logger interface {
	Debug(msg string)
	Debugf(fmt string, v ...interface{})

	Info(msg string)
	Infof(fmt string, v ...interface{})

	Warn(msg string)
	Warnf(fmt string, v ...interface{})

	Error(msg string)
	Errorf(fmt string, v ...interface{})

	Writer() io.Writer

	IsVerbose() bool
}

// WithSelectableWriter is an optional interface for loggers that want to support a separate writer per log level.
type WithSelectableWriter interface {
	WriterForLevel(level Level) io.Writer
}

// GetWriterForLevel retrieves the appropriate Writer for the log level provided.
//
// See WithSelectableWriter
func GetWriterForLevel(logger Logger, level Level) io.Writer {
	if er, ok := logger.(WithSelectableWriter); ok {
		return er.WriterForLevel(level)
	}

	return logger.Writer()
}

// IsQuiet defines whether a pack logger is set to quiet mode
func IsQuiet(logger Logger) bool {
	if writer := GetWriterForLevel(logger, InfoLevel); writer == ioutil.Discard {
		return true
	}

	return false
}

// Tip logs a tip.
func Tip(l Logger, format string, v ...interface{}) {
	l.Infof(style.Tip("Tip: ")+format, v...)
}
