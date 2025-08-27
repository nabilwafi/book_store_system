package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Init initializes the logger
func Init() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		os.Mkdir(logsDir, 0755)
	}

	// Create log file with current date
	currentTime := time.Now()
	logFileName := fmt.Sprintf("%s/app-%s.log", logsDir, currentTime.Format("2006-01-02"))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}

	// Set output to both file and stdout
	log.SetOutput(logFile)
	log.AddHook(&ConsoleHook{})
}

// ConsoleHook is a hook to log to console
type ConsoleHook struct{}

// Levels returns all log levels
func (hook *ConsoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire logs to console
func (hook *ConsoleHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	fmt.Print(line)
	return nil
}

// getCallerInfo gets the caller information
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "???"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// Info logs an info message
func Info(args ...interface{}) {
	if log == nil {
		Init()
	}
	log.WithField("caller", getCallerInfo()).Info(args...)
}

// Infof logs an info message with format
func Infof(format string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.WithField("caller", getCallerInfo()).Infof(format, args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	if log == nil {
		Init()
	}
	log.WithField("caller", getCallerInfo()).Error(args...)
}

// Errorf logs an error message with format
func Errorf(format string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.WithField("caller", getCallerInfo()).Errorf(format, args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	if log == nil {
		Init()
	}
	log.WithField("caller", getCallerInfo()).Warn(args...)
}

// Warnf logs a warning message with format
func Warnf(format string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.WithField("caller", getCallerInfo()).Warnf(format, args...)
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	if log == nil {
		Init()
	}
	log.WithField("caller", getCallerInfo()).Debug(args...)
}

// Debugf logs a debug message with format
func Debugf(format string, args ...interface{}) {
	if log == nil {
		Init()
	}
	log.WithField("caller", getCallerInfo()).Debugf(format, args...)
}