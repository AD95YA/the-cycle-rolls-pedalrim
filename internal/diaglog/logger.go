package diaglog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"pedalrim/internal/config"
)

// LogLevel represents the severity of a diagnostic log entry.
type LogLevel uint8

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Subsystem represents the functional area logging the message.
type Subsystem string

const (
	SubsystemCore    Subsystem = "CORE"
	SubsystemWorld   Subsystem = "WORLD"
	SubsystemPhysics Subsystem = "PHYSICS"
	SubsystemEconomy Subsystem = "ECONOMY"
	SubsystemQuest   Subsystem = "QUEST"
	SubsystemInput   Subsystem = "INPUT"
	SubsystemUI      Subsystem = "UI"
)

// RingBufferEntry holds an in-memory cached diagnostic message for HUD display.
// Fields ordered largest to smallest for cache alignment.
type RingBufferEntry struct {
	Timestamp time.Time
	Message   string
	Subsystem Subsystem
	Level     LogLevel
}

// Logger provides bounded, rotating file diagnostics and an in-memory ring buffer.
// Struct fields ordered from largest to smallest.
type Logger struct {
	// 8-byte pointers, slices, and maps
	recentEntries     []RingBufferEntry
	rateLimitTracker  map[string]time.Time
	currentFile       *os.File
	mutex             sync.Mutex
	currentFileBytes  int64
	maxFileSizeBytes  int64
	rateLimitDuration time.Duration

	// 4-byte integers
	ringBufferCapacity int32
	ringBufferHead     int32
	ringBufferCount    int32

	// Slices & Strings
	logFilePath       string
	rotatedFilePath   string

	// 1-byte / bool
	isClosed          bool
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// InitGlobalLogger initializes the default singleton logger.
func InitGlobalLogger(logPath string) (*Logger, error) {
	var initErr error
	once.Do(func() {
		defaultLogger, initErr = NewLogger(logPath, config.MaxLogSizeBytes, 64)
	})
	return defaultLogger, initErr
}

// GetLogger returns the global singleton logger.
func GetLogger() *Logger {
	if defaultLogger == nil {
		l, _ := InitGlobalLogger(config.LogFilePath)
		return l
	}
	return defaultLogger
}

// NewLogger constructs a bounded rotating logger.
func NewLogger(filePath string, maxSizeBytes int64, ringCapacity int32) (*Logger, error) {
	if maxSizeBytes <= 0 {
		maxSizeBytes = config.MaxLogSizeBytes
	}
	if ringCapacity <= 0 {
		ringCapacity = 64
	}

	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create log directory %q: %w", dir, err)
		}
	}

	logger := &Logger{
		logFilePath:        filePath,
		rotatedFilePath:    filePath + ".1",
		maxFileSizeBytes:   maxSizeBytes,
		recentEntries:      make([]RingBufferEntry, ringCapacity),
		ringBufferCapacity: ringCapacity,
		rateLimitTracker:   make(map[string]time.Time),
		rateLimitDuration:  time.Second, // Suppress identical logs within 1 second
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *Logger) openCurrentFile() error {
	fileInfo, err := os.Stat(l.logFilePath)
	if err == nil {
		l.currentFileBytes = fileInfo.Size()
		if l.currentFileBytes >= l.maxFileSizeBytes {
			if rotateErr := l.rotateFiles(); rotateErr != nil {
				return rotateErr
			}
		}
	} else if os.IsNotExist(err) {
		l.currentFileBytes = 0
	}

	f, err := os.OpenFile(l.logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open log file %q: %w", l.logFilePath, err)
	}
	l.currentFile = f
	return nil
}

func (l *Logger) rotateFiles() error {
	if l.currentFile != nil {
		_ = l.currentFile.Close()
		l.currentFile = nil
	}

	_ = os.Remove(l.rotatedFilePath)

	if _, err := os.Stat(l.logFilePath); err == nil {
		if err := os.Rename(l.logFilePath, l.rotatedFilePath); err != nil {
			return fmt.Errorf("failed to rotate log file: %w", err)
		}
	}

	l.currentFileBytes = 0
	return nil
}

// Log writes a structured diagnostic message to the bounded log file and ring buffer.
func (l *Logger) Log(level LogLevel, subsystem Subsystem, message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.isClosed {
		return
	}

	now := time.Now()

	// Rate-limit identical error/warning messages
	if level >= LevelWarn {
		cacheKey := string(subsystem) + ":" + message
		lastSeen, exists := l.rateLimitTracker[cacheKey]
		if exists && now.Sub(lastSeen) < l.rateLimitDuration {
			return
		}
		l.rateLimitTracker[cacheKey] = now
	}

	// Update in-memory ring buffer for debug HUD
	entry := RingBufferEntry{
		Timestamp: now,
		Level:     level,
		Subsystem: subsystem,
		Message:   message,
	}
	l.recentEntries[l.ringBufferHead] = entry
	l.ringBufferHead = (l.ringBufferHead + 1) % l.ringBufferCapacity
	if l.ringBufferCount < l.ringBufferCapacity {
		l.ringBufferCount++
	}

	// Format entry: [2026-11-12T12:00:00Z] [LEVEL] [SUBSYSTEM] message\n
	formatted := fmt.Sprintf("[%s] [%s] [%s] %s\n",
		now.UTC().Format(time.RFC3339),
		level.String(),
		string(subsystem),
		message,
	)

	entryBytes := int64(len(formatted))

	if l.currentFileBytes+entryBytes > l.maxFileSizeBytes {
		if err := l.rotateFiles(); err != nil {
			return
		}
		f, err := os.OpenFile(l.logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return
		}
		l.currentFile = f
	}

	if l.currentFile != nil {
		n, err := l.currentFile.WriteString(formatted)
		if err == nil {
			l.currentFileBytes += int64(n)
		}
	}
}

// Debug logs at LevelDebug.
func (l *Logger) Debug(subsystem Subsystem, message string) {
	l.Log(LevelDebug, subsystem, message)
}

// Info logs at LevelInfo.
func (l *Logger) Info(subsystem Subsystem, message string) {
	l.Log(LevelInfo, subsystem, message)
}

// Warn logs at LevelWarn.
func (l *Logger) Warn(subsystem Subsystem, message string) {
	l.Log(LevelWarn, subsystem, message)
}

// Error logs at LevelError.
func (l *Logger) Error(subsystem Subsystem, message string) {
	l.Log(LevelError, subsystem, message)
}

// GetRecentEntries returns a copy of recently recorded log entries in chronological order.
func (l *Logger) GetRecentEntries() []RingBufferEntry {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	result := make([]RingBufferEntry, l.ringBufferCount)
	if l.ringBufferCount == 0 {
		return result
	}

	start := (l.ringBufferHead - l.ringBufferCount + l.ringBufferCapacity) % l.ringBufferCapacity
	for i := int32(0); i < l.ringBufferCount; i++ {
		idx := (start + i) % l.ringBufferCapacity
		result[i] = l.recentEntries[idx]
	}
	return result
}

// Close closes the underlying log file.
func (l *Logger) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.isClosed = true
	if l.currentFile != nil {
		err := l.currentFile.Close()
		l.currentFile = nil
		return err
	}
	return nil
}

