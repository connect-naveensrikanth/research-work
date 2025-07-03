package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

type WriteRequest struct {
	Key   string
	Value string
}
type KeyValueStore struct {
	data map[string]string
	mu   sync.Mutex
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{data: make(map[string]string)}
}
func (kv *KeyValueStore) BatchWrite(batch []WriteRequest) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	for _, req := range batch {
		kv.data[req.Key] = req.Value
	}
}
func (kv *KeyValueStore) Get(key string) (string, bool) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	val, ok := kv.data[key]
	return val, ok
}

type Batcher struct {
	store        *KeyValueStore
	input        chan WriteRequest
	batchSize    int
	flushPeriod  time.Duration
	batchCount   int
	totalWrites  int
	stopChan     chan struct{}
	shutdownDone chan struct{}
}

func NewBatcher(store *KeyValueStore, size int, period time.Duration) *Batcher {
	return &Batcher{
		store:        store,
		input:        make(chan WriteRequest, 1000),
		batchSize:    size,
		flushPeriod:  period,
		stopChan:     make(chan struct{}),
		shutdownDone: make(chan struct{}),
	}
}
func (b *Batcher) Start() {
	go func() {
		batch := make([]WriteRequest, 0, b.batchSize)
		timer := time.NewTimer(b.flushPeriod)
		for {
			select {
			case req := <-b.input:
				batch = append(batch, req)
				if len(batch) >= b.batchSize {
					b.flush(batch)
					batch = batch[:0]
					if !timer.Stop() {
						<-timer.C
					}
					timer.Reset(b.flushPeriod)
				}
			case <-timer.C:
				if len(batch) > 0 {
					b.flush(batch)
					batch = batch[:0]
				}
				timer.Reset(b.flushPeriod)
			case <-b.stopChan:
				if len(batch) > 0 {
					b.flush(batch)
				}
				close(b.shutdownDone)
				return
			}
		}
	}()
}
func (b *Batcher) flush(batch []WriteRequest) {
	b.store.BatchWrite(batch)
	b.batchCount++
	b.totalWrites += len(batch)
	fmt.Printf("Flushed batch of %d writes. Total writes: %d\n", len(batch), b.totalWrites)
}
func (b *Batcher) Stop() {
	close(b.stopChan)
	<-b.shutdownDone
}
func main() {
	store := NewKeyValueStore()
	batcher := NewBatcher(store, 10, 3*time.Second)
	batcher.Start()
	go func() {
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("key%d", i)
			val := fmt.Sprintf("val%d", i)
			batcher.input <- WriteRequest{Key: key, Value: val}
			time.Sleep(150 * time.Millisecond)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	fmt.Println("\nShutting down...")
	batcher.Stop()
	fmt.Printf("Final Stats — Batches: %d, Total Writes: %d\n", batcher.batchCount, batcher.totalWrites)
}
