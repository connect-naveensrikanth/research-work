package main

import (
	"fmt"
	"os"
	"sync"
	"time"
	"math/rand"
)

type LogEntry struct {
	ID      int
	Message string
}

type Node struct {
	ID   int
	Logs []LogEntry
}

func (n *Node) WriteLog(entry LogEntry) {
	n.Logs = append(n.Logs, entry)
}

func replicateLogs(nodes []Node, entry LogEntry, wg *sync.WaitGroup) {
	for i := range nodes {
		go func(node *Node) {
			defer wg.Done()
			node.WriteLog(entry)
		}(&nodes[i])
	}
}

func simulateNetworkDelay() {
	delay := rand.Intn(100)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

func measureDiskThroughput() int {
	start := time.Now()
	file, _ := os.Create("testfile.txt")
	defer file.Close()
	for i := 0; i < 1000; i++ {
		file.WriteString(fmt.Sprintf("Line %d\n", i))
	}
	elapsed := time.Since(start)
	throughput := int(float64(1000) / elapsed.Seconds())
	return throughput
}

func logReplication(nodes []Node, entry LogEntry, wg *sync.WaitGroup) {
	simulateNetworkDelay()
	replicateLogs(nodes, entry, wg)
}

func main() {
	nodes := []Node{{ID: 1}, {ID: 2}, {ID: 3}}
	entry := LogEntry{ID: 1, Message: "Initial Configuration"}
	var wg sync.WaitGroup

	wg.Add(len(nodes))
	logReplication(nodes, entry, &wg)
	wg.Wait()

	throughput := measureDiskThroughput()
	fmt.Println("Disk throughput after initial replication:", throughput, "MB/s")

	nodes2 := []Node{{ID: 4}, {ID: 5}, {ID: 6}}
	entry2 := LogEntry{ID: 2, Message: "Updated Configuration"}
	var wg2 sync.WaitGroup

	wg2.Add(len(nodes2))
	logReplication(nodes2, entry2, &wg2)
	wg2.Wait()

	throughput2 := measureDiskThroughput()
	fmt.Println("Disk throughput after second replication:", throughput2, "MB/s")

	nodes3 := []Node{{ID: 7}, {ID: 8}, {ID: 9}}
	entry3 := LogEntry{ID: 3, Message: "Service Restart"}
	var wg3 sync.WaitGroup

	wg3.Add(len(nodes3))
	logReplication(nodes3, entry3, &wg3)
	wg3.Wait()

	throughput3 := measureDiskThroughput()
	fmt.Println("Disk throughput after third replication:", throughput3, "MB/s")

	nodes4 := []Node{{ID: 10}, {ID: 11}, {ID: 12}}
	entry4 := LogEntry{ID: 4, Message: "Backup Completed"}
	var wg4 sync.WaitGroup

	wg4.Add(len(nodes4))
	logReplication(nodes4, entry4, &wg4)
	wg4.Wait()

	throughput4 := measureDiskThroughput()
	fmt.Println("Disk throughput after fourth replication:", throughput4, "MB/s")
}
