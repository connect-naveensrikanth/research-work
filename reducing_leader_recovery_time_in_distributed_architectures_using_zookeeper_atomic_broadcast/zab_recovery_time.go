package main

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	Term    int
	Index   int
	Command string
}

type Replica struct {
	id          int
	log         []Message
	currentTerm int
	votedFor    int
	mu          sync.Mutex
}

type ZAB struct {
	replicas []*Replica
	leader   *Replica
	mu       sync.Mutex
}

func NewZAB(replicasCount int) *ZAB {
	zab := &ZAB{}
	for i := 0; i < replicasCount; i++ {
		zab.replicas = append(zab.replicas, &Replica{id: i})
	}
	return zab
}

func (z *ZAB) ElectLeader() {
	z.mu.Lock()
	defer z.mu.Unlock()

	for _, replica := range z.replicas {
		if replica.currentTerm == 0 {
			replica.currentTerm = 1
			replica.votedFor = replica.id
			z.leader = replica
			fmt.Printf("Replica %d becomes the leader\n", replica.id)
			break
		}
	}
}

func (z *ZAB) AppendLogEntry(command string) {
	z.mu.Lock()
	defer z.mu.Unlock()

	if z.leader == nil {
		fmt.Println("No leader available")
		return
	}

	logEntry := Message{
		Term:    z.leader.currentTerm,
		Index:   len(z.leader.log) + 1,
		Command: command,
	}
	z.leader.log = append(z.leader.log, logEntry)
	fmt.Printf("Replica %d added log entry: %v\n", z.leader.id, logEntry)
}

func (z *ZAB) CommitLog() {
	z.mu.Lock()
	defer z.mu.Unlock()

	if z.leader == nil {
		fmt.Println("No leader to commit log")
		return
	}

	fmt.Printf("Leader %d committing logs\n", z.leader.id)
}
func (z *ZAB) Run() {
	for {
		time.Sleep(2 * time.Second)
		z.ElectLeader()
		z.AppendLogEntry("Sample Command")
		z.CommitLog()
	}
}
func main() {
	zab := NewZAB(5)
	go zab.Run()
	select {}
}
