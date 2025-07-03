package main

import (
	"fmt"
	"sync"
	"time"
)

type LatencyMetrics struct {
	count int
	total time.Duration
	min   time.Duration
	max   time.Duration
	mu    sync.Mutex
}

func NewLatencyMetrics() *LatencyMetrics {
	return &LatencyMetrics{
		min: time.Duration(1<<63 - 1),
	}
}

func (m *LatencyMetrics) Record(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.count++
	m.total += latency
	if latency < m.min {
		m.min = latency
	}
	if latency > m.max {
		m.max = latency
	}
}

func (m *LatencyMetrics) Print() {
	m.mu.Lock()
	defer m.mu.Unlock()
	avg := time.Duration(0)
	if m.count > 0 {
		avg = m.total / time.Duration(m.count)
	}
	fmt.Printf("Count: %d\nMin: %v\nMax: %v\nAvg: %v\n", m.count, m.min, m.max, avg)
}

func simulateWork(m *LatencyMetrics) {
	for i := 0; i < 50; i++ {
		start := time.Now()
		time.Sleep(time.Duration(10+5*i%20) * time.Millisecond)
		latency := time.Since(start)
		m.Record(latency)
	}
}

func main() {
	metrics := NewLatencyMetrics()
	simulateWork(metrics)
	metrics.Print()
}
