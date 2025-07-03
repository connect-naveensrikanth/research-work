package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Replica struct {
	id       int
	isLeader bool
}
type Cluster struct {
	replicas []*Replica
	leader   *Replica
	mu       sync.Mutex
}

func NewCluster(size int) *Cluster {
	cluster := &Cluster{
		replicas: make([]*Replica, size),
	}
	for i := 0; i < size; i++ {
		cluster.replicas[i] = &Replica{id: i}
	}
	return cluster
}
func (c *Cluster) ElectLeader() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, replica := range c.replicas {
		replica.isLeader = false
	}
	leader := c.replicas[rand.Intn(len(c.replicas))]
	leader.isLeader = true
	c.leader = leader
}
func (c *Cluster) FailLeader() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.leader != nil {
		c.leader.isLeader = false
		c.leader = nil
	}
}
func (c *Cluster) RecoveryTime() time.Duration {
	start := time.Now()
	c.FailLeader()
	c.ElectLeader()
	return time.Since(start)
}
func main() {
	rand.Seed(time.Now().UnixNano())
	cluster := NewCluster(5)
	cluster.ElectLeader()
	fmt.Printf("Initial leader: Replica %d\n", cluster.leader.id)
	recoveryTime := cluster.RecoveryTime()
	fmt.Printf("Leader failure recovery time: %v\n", recoveryTime)
}
