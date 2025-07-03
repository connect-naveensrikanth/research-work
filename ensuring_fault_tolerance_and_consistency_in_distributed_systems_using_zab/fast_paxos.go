package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Proposal struct {
	Value string
	ID    int
}

type FastPaxos struct {
	quorumSize int
	accepted   map[int]Proposal
	prepared   map[int]Proposal
	mu         sync.Mutex
}

func NewFastPaxos(quorumSize int) *FastPaxos {
	return &FastPaxos{
		quorumSize: quorumSize,
		accepted:   make(map[int]Proposal),
		prepared:   make(map[int]Proposal),
	}
}

func (fp *FastPaxos) Propose(value string) int {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	id := rand.Int()
	proposal := Proposal{Value: value, ID: id}
	fp.prepared[id] = proposal
	return id
}

func (fp *FastPaxos) Prepare(id int, value string) bool {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	if _, exists := fp.prepared[id]; exists {
		proposal := Proposal{Value: value, ID: id}
		fp.accepted[id] = proposal
		return true
	}
	return false
}

func (fp *FastPaxos) Commit(id int) string {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	if proposal, exists := fp.accepted[id]; exists {
		delete(fp.accepted, id)
		return proposal.Value
	}
	return ""
}

func simulateWrite(fp *FastPaxos, value string) {
	start := time.Now()
	id := fp.Propose(value)
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	committed := fp.Prepare(id, value)
	if committed {
		result := fp.Commit(id)
		latency := time.Since(start)
		fmt.Printf("Write Value: %s, Commit Result: %s, Latency: %v\n", value, result, latency)
	} else {
		latency := time.Since(start)
		fmt.Printf("Write Value: %s, Commit Failed, Latency: %v\n", value, latency)
	}
}

func main() {
	fp := NewFastPaxos(3)
	for i := 0; i < 10; i++ {
		go simulateWrite(fp, fmt.Sprintf("Value%d", i))
	}
	time.Sleep(2 * time.Second)
}
