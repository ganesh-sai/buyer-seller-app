// Package logging - responsible for logging related activities
package logging

import (
	"fmt"
	"io"
	"log"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

var logger *Logger
var LogPrefixes = map[int]string{
	DEBUG: "DEBUG",
	INFO:  "INFO ",
	WARN:  "WARN ",
	ERROR: "ERROR",
}

//Logger ...

type Logger struct {
	Config
}

// GetLogger - gives a logger instance to log
func GetLogger() *Logger {
	return logger
}

// Config - Input for initializing a logger
type Config struct {
	LogLevel int       // The LogLevel you want to log
	Output   io.Writer // Where exactly you want to log
	Prefix   string    // Any prefix that you want your log to have
}

// Debugf - Logs formatted string with DEBUG severity

func (l *Logger) Debugf(format string, n ...interface{}) {
	l.Logf(DEBUG, format, n...)
}

// Infof - Logs formatted string with INFO severity
func (l *Logger) Infof(format string, n ...interface{}) {
	l.Logf(INFO, format, n...)
}

// Warnf - Logs formatted string with WARN severity
func (l *Logger) Warnf(format string, n ...interface{}) {
	l.Logf(WARN, format, n...)
}

// Errorf - Logs formatted string with ERROR severity
func (l *Logger) Errorf(format string, n ...interface{}) {
	l.Logf(ERROR, format, n...)
}

// Logf - logs the formatted message
func (l *Logger) Logf(level int, s string, n ...interface{}) {
	if level >= l.LogLevel {
		log.Println(l.LogPrefix(level), fmt.Sprintf(s, n...))
	}
}

// Debug - Logs with Debug Severity
func (l *Logger) Debug(n ...interface{}) {
	l.Log(DEBUG, n...)
}

// Info - Logs with Info severity
func (l *Logger) Info(n ...interface{}) {
	l.Log(INFO, n...)
}

// Warn - Logs with Warn severity
func (l *Logger) Warn(n ...interface{}) {
	l.Log(WARN, n...)
}

// Error - Logs with ERROR severity
func (l *Logger) Error(n ...interface{}) {
	l.Log(ERROR, n...)
}

// LogPrefix - Appends any prefix that you wanted to have, like app name, component name
func (l *Logger) LogPrefix(i int) (s string) {
	if l.Prefix != "" {
		s = s + " [" + l.Prefix + "]"
	}
	s = s + " " + l.LogLevelPrefix(i)
	return
}

// LogLevelPrefix - Adds Logging Level before the actual log
func (l *Logger) LogLevelPrefix(level int) (s string) {
	prefix := LogPrefixes[level]
	return prefix
}

func (l *Logger) Log(level int, n ...interface{}) {
	if level >= l.LogLevel {
		all := append([]interface{}{l.LogPrefix(level) + ":"}, n...)
		log.Println(all...)
	}
}

// Init Creates a logger
func Init(config Config) {
	logger = &Logger{config}
}
