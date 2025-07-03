package main

import (
	"fmt"
	"math/rand"
	"time"
)

func simulateLeaseBasedLatency(clusterSize int) float64 {
	latency := 0.2 + float64(clusterSize-3)*0.1
	latency += (rand.Float64() - 0.5) * 0.1
	if latency < 0.2 {
		latency = 0.2
	}
	return latency
}

func simulateClusterFailure(clusterSize int) bool {
	failChance := rand.Float64()
	if failChance < 0.1 {
		return true
	}
	return false
}

func simulateSingleRead(clusterSize int) (float64, error) {
	if simulateClusterFailure(clusterSize) {
		return 0, fmt.Errorf("cluster failure detected at size %d", clusterSize)
	}
	latency := simulateLeaseBasedLatency(clusterSize)
	return latency, nil
}

func calculateTotalLatency(clusterSize int, numReads int) (float64, error) {
	totalLatency := 0.0
	for i := 0; i < numReads; i++ {
		latency, err := simulateSingleRead(clusterSize)
		if err != nil {
			return 0, err
		}
		totalLatency += latency
	}
	return totalLatency, nil
}

func displayLatencyReport(clusterSizes []int, numReads int) {
	fmt.Println("Cluster Size (Nodes) | Lease-Based Latency (ms) | Total Latency for Reads (ms) | Error Occurred")
	for _, size := range clusterSizes {
		latency := simulateLeaseBasedLatency(size)
		totalLatency, err := calculateTotalLatency(size, numReads)
		if err != nil {
			fmt.Printf("%d | %.2f | N/A | %s\n", size, latency, err)
		} else {
			fmt.Printf("%d | %.2f | %.2f | No\n", size, latency, totalLatency)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	clusterSizes := []int{3, 5, 7, 9, 11}
	numReads := 100
	displayLatencyReport(clusterSizes, numReads)
}
