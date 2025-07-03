package main

import (
	"fmt"
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

func NewCluster(numReplicas int) *Cluster {
	cluster := &Cluster{}
	for i := 0; i < numReplicas; i++ {
		cluster.replicas = append(cluster.replicas, &Replica{id: i})
	}
	return cluster
}
func (c *Cluster) ElectLeader() *Replica {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, replica := range c.replicas {
		if replica.isLeader {
			replica.isLeader = false
		}
	}
	newLeader := c.replicas[0]
	newLeader.isLeader = true
	c.leader = newLeader
	return c.leader
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
	startTime := time.Now()

	c.FailLeader()
	c.ElectLeader()

	return time.Since(startTime)
}
func main() {
	cluster := NewCluster(5)
	cluster.ElectLeader()
	fmt.Printf("Initial Leader: Replica %d\n", cluster.leader.id)
	recoveryTime := cluster.RecoveryTime()
	fmt.Printf("Leader failure recovery time: %v\n", recoveryTime)
}
