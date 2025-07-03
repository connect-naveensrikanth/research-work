package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	startTime := time.Now()
	numReads := 1000

	for i := 0; i < numReads; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := cli.Get(ctx, "my-key", clientv3.WithSerializable())
		if err != nil {
			log.Println("Error during read operation:", err)
			continue
		}
	}

	elapsedTime := time.Since(startTime)
	qps := float64(numReads) / elapsedTime.Seconds()

	fmt.Printf("Total Reads: %d\n", numReads)
	fmt.Printf("Elapsed Time: %s\n", elapsedTime)
	fmt.Printf("QPS (Queries Per Second): %.2f\n", qps)
}
