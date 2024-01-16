package goes

import (
	"fmt"
	"log"
)

type Logger interface {
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	Fatal(msg string, keysAndValues ...any)
}

func NewLogger(base any) Logger {
	logger, ok := base.(Logger)
	if ok {
		return logger
	}

	if log, ok := base.(*log.Logger); ok {
		return &stdLogger{logger: log}
	}

	return &stdLogger{logger: log.Default()}
}

type stdLogger struct {
	logger *log.Logger
}

func (l *stdLogger) Info(msg string, keysAndValues ...any) {
	l.logger.Printf(fmt.Sprint("[INFO] %v", msg), zipKeysAndValues(keysAndValues)...)
}
func (l *stdLogger) Warn(msg string, keysAndValues ...any) {
	l.logger.Printf(fmt.Sprint("[WARN] %v", msg), zipKeysAndValues(keysAndValues)...)
}
func (l *stdLogger) Error(msg string, keysAndValues ...any) {
	l.logger.Printf(fmt.Sprint("[ERROR] %v", msg), zipKeysAndValues(keysAndValues)...)
}
func (l *stdLogger) Fatal(msg string, keysAndValues ...any) {
	l.logger.Fatalf(fmt.Sprint("[FATAL] %v", msg), zipKeysAndValues(keysAndValues)...)
}

func zipKeysAndValues(keysAndValues []any) []any {
	var result []any
	length := len(keysAndValues)
	for i := 0; i < length; i += 2 {
		if i == length-1 {
			// If this is the last element and there's no pair, append "missing" as the value.
			result = append(result, fmt.Sprint("%v=missing", keysAndValues[i])+"=missing")
		} else {
			result = append(result, fmt.Sprint("%v=%v", keysAndValues[i], keysAndValues[i+1]))
		}
	}
	return result
}
