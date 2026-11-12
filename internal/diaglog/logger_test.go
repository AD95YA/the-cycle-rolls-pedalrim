package diaglog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestLoggerRotationAndSizeCapping(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test_pedalrim.log")
	maxBytes := int64(4 * 1024) // 4 KB for quick rotation testing

	logger, err := NewLogger(logPath, maxBytes, 16)
	if err != nil {
		t.Fatalf("failed to create test logger: %v", err)
	}
	defer logger.Close()

	// Write enough entries to force multiple rotations
	largePayload := "This is a diagnostic message simulating cycling telemetry and physics state."
	for i := 0; i < 200; i++ {
		logger.Info(SubsystemPhysics, fmt.Sprintf("[%03d] %s", i, largePayload))
	}

	// Verify main file exists and is within bounds
	mainInfo, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("expected main log file to exist: %v", err)
	}
	if mainInfo.Size() > maxBytes {
		t.Errorf("main log file size %d exceeds maxBytes %d", mainInfo.Size(), maxBytes)
	}

	// Verify rotated backup file exists and is within bounds
	rotatedPath := logPath + ".1"
	rotInfo, err := os.Stat(rotatedPath)
	if err != nil {
		t.Fatalf("expected rotated log file to exist: %v", err)
	}
	if rotInfo.Size() > maxBytes {
		t.Errorf("rotated log file size %d exceeds maxBytes %d", rotInfo.Size(), maxBytes)
	}

	totalDisk := mainInfo.Size() + rotInfo.Size()
	if totalDisk > 2*maxBytes {
		t.Errorf("total disk footprint %d exceeds 2*maxBytes %d", totalDisk, 2*maxBytes)
	}
}

func TestLoggerRateLimiting(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "rate_limit.log")

	logger, err := NewLogger(logPath, 64*1024, 32)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Close()

	// Log the same warning 10 times in rapid succession
	for i := 0; i < 10; i++ {
		logger.Warn(SubsystemPhysics, "chain slip detected on high gear")
	}

	entries := logger.GetRecentEntries()
	warnCount := 0
	for _, entry := range entries {
		if entry.Message == "chain slip detected on high gear" {
			warnCount++
		}
	}

	if warnCount != 1 {
		t.Errorf("expected exactly 1 warning due to rate-limiting, got %d", warnCount)
	}
}

func TestLoggerRingBufferOrder(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "ring_order.log")

	capacity := int32(5)
	logger, err := NewLogger(logPath, 64*1024, capacity)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Close()

	for i := 1; i <= 8; i++ {
		logger.Info(SubsystemCore, fmt.Sprintf("Event %d", i))
	}

	entries := logger.GetRecentEntries()
	if len(entries) != int(capacity) {
		t.Fatalf("expected %d entries, got %d", capacity, len(entries))
	}

	// Oldest visible should be Event 4, newest Event 8
	expected := []string{"Event 4", "Event 5", "Event 6", "Event 7", "Event 8"}
	for i, entry := range entries {
		if entry.Message != expected[i] {
			t.Errorf("entry [%d] expected %q, got %q", i, expected[i], entry.Message)
		}
	}
}

func TestLoggerConcurrentWrites(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "concurrent.log")

	logger, err := NewLogger(logPath, 32*1024, 64)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Close()

	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for step := 0; step < 25; step++ {
				logger.Info(SubsystemWorld, fmt.Sprintf("worker %d step %d", id, step))
				time.Sleep(10 * time.Microsecond)
			}
		}(worker)
	}

	wg.Wait()
}

