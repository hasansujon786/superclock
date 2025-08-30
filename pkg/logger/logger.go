package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// Logger struct
type Logger struct {
	dump io.Writer
}

var globalLogger *Logger

// Init initializes the global logger if DEBUG is set
func Init() {
	if _, ok := os.LookupEnv("DEBUG"); ok {
		file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err == nil {
			globalLogger = &Logger{dump: file}
		}
	}
}

// Close closes the logfile (optional, but good for cleanup)
func Close() {
	if globalLogger != nil {
		if f, ok := globalLogger.dump.(*os.File); ok {
			f.Close()
		}
	}
}

// log is an internal helper for leveled logging
func log(level string, v ...any) {
	if globalLogger == nil || globalLogger.dump == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(globalLogger.dump, "[%s] [%s] ", timestamp, level)

	// Dump all args using spew
	for _, arg := range v {
		spew.Fdump(globalLogger.dump, arg)
	}
}

// Public API
func Info(v ...any)  { log("INFO", v...) }
func Warn(v ...any)  { log("WARN", v...) }
func Error(v ...any) { log("ERROR", v...) }
func Debug(v ...any) { log("DEBUG", v...) }
