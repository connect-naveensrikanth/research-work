package main

import (
	"fmt"
	"os"
	"time"
	"strconv"
	"math/rand"
	"strings"
)

type WAL struct {
	logFile *os.File
}

func NewWAL(logFilePath string) (*WAL, error) {
	file, err := os.Create(logFilePath)
	if err != nil {
		return nil, err
	}
	return &WAL{logFile: file}, nil
}

func (wal *WAL) WriteLog(entry string) error {
	_, err := wal.logFile.WriteString(entry + "\n")
	return err
}

func (wal *WAL) Close() error {
	return wal.logFile.Close()
}

func generateLogEntry(id int, entrySize int) string {
	timestamp := time.Now().UnixNano()
	entryData := fmt.Sprintf("LogEntryID-%d-Timestamp-%d-Data-%s", id, timestamp, string(make([]byte, entrySize)))
	return entryData
}

func simulateWriteLoad(wal *WAL, numWrites int, entrySize int) (int64, error) {
	start := time.Now()
	for i := 0; i < numWrites; i++ {
		entry := generateLogEntry(i, entrySize)
		if err := wal.WriteLog(entry); err != nil {
			return 0, err
		}
	}
	duration := time.Since(start)
	return duration.Milliseconds(), nil
}

func simulateMultipleFiles(logPath string, numFiles int, numWrites int, entrySize int) (int64, error) {
	totalDuration := int64(0)
	totalSize := int64(0)

	for i := 0; i < numFiles; i++ {
		filePath := fmt.Sprintf("%s_%d.log", logPath, i)
		wal, err := NewWAL(filePath)
		if err != nil {
			return 0, err
		}
		defer wal.Close()

		duration, err := simulateWriteLoad(wal, numWrites, entrySize)
		if err != nil {
			return 0, err
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return 0, err
		}
		fileSize := fileInfo.Size()

		totalDuration += duration
		totalSize += fileSize
	}

	return totalDuration, totalSize
}

func simulateRandomLoad(wal *WAL, numWrites int, minEntrySize int, maxEntrySize int) (int64, error) {
	start := time.Now()
	for i := 0; i < numWrites; i++ {
		entrySize := rand.Intn(maxEntrySize-minEntrySize) + minEntrySize
		entry := generateLogEntry(i, entrySize)
		if err := wal.WriteLog(entry); err != nil {
			return 0, err
		}
	}
	duration := time.Since(start)
	return duration.Milliseconds(), nil
}

func main() {
	logPath := "wal_log"
	numFiles := 5
	numWrites := 10000
	entrySize := 128

	totalDuration, totalSize, err := simulateMultipleFiles(logPath, numFiles, numWrites, entrySize)
	if err != nil {
		fmt.Println("Error during write load simulation:", err)
		return
	}

	throughput := float64(totalSize) / float64(totalDuration) * 1000
	fmt.Printf("Total Disk throughput across %d files: %.2f MB/s\n", numFiles, throughput/1024/1024)
	fmt.Printf("Total log entries written: %d\n", numFiles*numWrites)
	fmt.Printf("Total size written: %d bytes\n", totalSize)

	// Simulate a random load
	rand.Seed(time.Now().UnixNano())
	wal, err := NewWAL("random_wal.log")
	if err != nil {
		fmt.Println("Error creating WAL:", err)
		return
	}
	defer wal.Close()

	randDuration, randSize, err := simulateRandomLoad(wal, numWrites, 64, 1024)
	if err != nil {
		fmt.Println("Error during random load simulation:", err)
		return
	}

	randThroughput := float64(randSize) / float64(randDuration) * 1000
	fmt.Printf("Random Disk throughput: %.2f MB/s\n", randThroughput/1024/1024)
	fmt.Printf("Total random log entries written: %d\n", numWrites)
	fmt.Printf("Total random size written: %d bytes\n", randSize)
}
