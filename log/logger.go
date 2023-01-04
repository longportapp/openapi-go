package log

import (
	"log"
	"strings"
)

type LogLevel int

const (
	LevelDebug LogLevel = -1
	LevelInfo  LogLevel = iota
	LevelWarn
	LevelError
)

var logger Logger

func init() {
	dlogger := &DefaultLogger{}
	logger = dlogger
}

// SetLogger set another logger
func SetLogger(nlog Logger) {
	logger = nlog
}

// Logger is the interface for Logger types
type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})

	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	Debugf(fmt string, args ...interface{})
}

// DefaultLogger implement by fmt
type DefaultLogger struct {
	lvl LogLevel
}

func (l *DefaultLogger) SetLevel(lvl string) {
	lvl = strings.ToLower(lvl)

	switch lvl {
	case "info":
		l.lvl = LevelInfo
	case "debug":
		l.lvl = LevelDebug
	case "warn":
		l.lvl = LevelWarn
	case "error":
		l.lvl = LevelError
	}
}

func (l *DefaultLogger) Info(args ...interface{}) {
	if l.lvl <= LevelInfo {
		log.Println("[INFO]", args)
	}

}

func (l *DefaultLogger) Infof(msg string, args ...interface{}) {
	if l.lvl <= LevelInfo {
		log.Printf("[INFO] "+msg, args...)
	}

}

func (l *DefaultLogger) Error(args ...interface{}) {
	if l.lvl <= LevelError {
		log.Println("[ERR]", args)
	}
}

func (l *DefaultLogger) Errorf(msg string, args ...interface{}) {
	if l.lvl <= LevelError {
		log.Printf("[ERR] "+msg, args...)
	}
}

func (l *DefaultLogger) Debug(args ...interface{}) {
	if l.lvl <= LevelDebug {
		log.Println("[DEBUG]", args)
	}
}

func (l *DefaultLogger) Debugf(msg string, args ...interface{}) {
	if l.lvl <= LevelDebug {
		log.Printf("[DEBUG] "+msg, args...)
	}
}

func (l *DefaultLogger) Warn(args ...interface{}) {
	if l.lvl <= LevelWarn {
		log.Println("[WARN]", args)
	}
}

func (l *DefaultLogger) Warnf(msg string, args ...interface{}) {
	if l.lvl <= LevelWarn {
		log.Printf("[WARN] "+msg, args...)
	}
}
