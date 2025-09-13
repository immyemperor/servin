package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging capabilities
type Logger struct {
	level    LogLevel
	logger   *log.Logger
	file     *os.File
	verbose  bool
	withFile bool
}

var defaultLogger *Logger

func init() {
	defaultLogger = &Logger{
		level:   INFO,
		logger:  log.New(os.Stderr, "", 0),
		verbose: false,
	}
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, verbose bool, logFile string) (*Logger, error) {
	l := &Logger{
		level:   level,
		verbose: verbose,
	}

	var writers []io.Writer
	writers = append(writers, os.Stderr)

	// Add file logging if specified
	if logFile != "" {
		// Ensure log directory exists
		logDir := filepath.Dir(logFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %v", err)
		}

		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
		l.file = file
		l.withFile = true
		writers = append(writers, file)
	}

	// Create multi-writer for both console and file output
	multiWriter := io.MultiWriter(writers...)
	l.logger = log.New(multiWriter, "", 0)

	return l, nil
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetVerbose enables or disables verbose logging
func (l *Logger) SetVerbose(verbose bool) {
	l.verbose = verbose
}

// Close closes the log file if it was opened
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

// formatMessage formats a log message with timestamp, level, and caller info
func (l *Logger) formatMessage(level LogLevel, msg string, withCaller bool) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var caller string
	if withCaller || l.verbose {
		_, file, line, ok := runtime.Caller(3)
		if ok {
			caller = fmt.Sprintf(" [%s:%d]", filepath.Base(file), line)
		}
	}

	return fmt.Sprintf("%s [%s]%s %s", timestamp, level.String(), caller, msg)
}

// log is the internal logging method
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	msg := fmt.Sprintf(format, args...)
	formatted := l.formatMessage(level, msg, level >= ERROR)

	l.logger.Print(formatted)

	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

// Global logging functions using the default logger
func SetLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

func SetVerbose(verbose bool) {
	defaultLogger.SetVerbose(verbose)
}

func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

// InitLogger initializes the global logger with file logging
func InitLogger(level LogLevel, verbose bool, logFile string) error {
	logger, err := NewLogger(level, verbose, logFile)
	if err != nil {
		return err
	}
	defaultLogger = logger
	return nil
}

// Close closes the global logger
func Close() {
	defaultLogger.Close()
}

// GetLogPath returns the appropriate log file path for the current platform
func GetLogPath() string {
	var logDir string

	switch runtime.GOOS {
	case "windows":
		homeDir, _ := os.UserHomeDir()
		logDir = filepath.Join(homeDir, ".servin", "logs")
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		logDir = filepath.Join(homeDir, ".servin", "logs")
	case "linux":
		logDir = "/var/log/servin"
	default:
		logDir = "/var/log/servin"
	}

	return filepath.Join(logDir, "servin.log")
}
