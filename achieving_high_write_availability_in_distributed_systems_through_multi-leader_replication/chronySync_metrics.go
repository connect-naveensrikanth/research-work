package main

import (
	"fmt"
	"time"
)

type TimeMetrics struct {
	OperationName string
	StartTime     time.Time
	EndTime       time.Time
	Duration      time.Duration
}

func (t *TimeMetrics) Start() {
	t.StartTime = time.Now()
}

func (t *TimeMetrics) End() {
	t.EndTime = time.Now()
	t.Duration = t.EndTime.Sub(t.StartTime)
}

func (t *TimeMetrics) Report() {
	fmt.Printf("Operation: %s\n", t.OperationName)
	fmt.Printf("Start Time: %s\n", t.StartTime.Format(time.RFC3339))
	fmt.Printf("End Time: %s\n", t.EndTime.Format(time.RFC3339))
	fmt.Printf("Duration: %v\n", t.Duration)
}

func someOperation() {
	time.Sleep(2 * time.Second)
}

func anotherOperation() {
	time.Sleep(1 * time.Second)
}

func main() {
	operation1 := &TimeMetrics{OperationName: "Operation 1"}
	operation2 := &TimeMetrics{OperationName: "Operation 2"}

	operation1.Start()
	someOperation()
	operation1.End()
	operation2.Start()
	anotherOperation()
	operation2.End()
	operation1.Report()
	operation2.Report()
}
