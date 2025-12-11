package tunnel

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type LogEntry struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}

var (
	logMu      sync.RWMutex
	logEntries []LogEntry
	maxEntries = 200
)

func logWithLevel(level, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)

	log.Printf("[%s] %s", level, msg)

	logMu.Lock()
	defer logMu.Unlock()

	logEntries = append(logEntries, LogEntry{
		Time:    time.Now(),
		Level:   level,
		Message: msg,
	})

	if len(logEntries) > maxEntries {
		diff := len(logEntries) - maxEntries
		logEntries = logEntries[diff:]
	}
}

func Infof(format string, args ...interface{}) {
	logWithLevel("INFO", format, args...)
}

func Errorf(format string, args ...interface{}) {
	logWithLevel("ERROR", format, args...)
}

func Debugf(format string, args ...interface{}) {
	logWithLevel("DEBUG", format, args...)
}

func GetLogEntries() []LogEntry {
	logMu.RLock()
	defer logMu.RUnlock()

	out := make([]LogEntry, len(logEntries))
	copy(out, logEntries)
	return out
}
