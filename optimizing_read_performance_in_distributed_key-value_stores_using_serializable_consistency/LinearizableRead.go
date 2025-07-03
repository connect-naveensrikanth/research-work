package main

import (
	"context"
	"fmt"
	"time"
	"go.etcd.io/etcd/client/v3"
)

func main() {
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	defer cli.Close()
	ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel1()
	linearResp, _ := cli.Get(ctx1, "my-key")
	fmt.Println("Linearizable Read:")
	for _, kv := range linearResp.Kvs {
		fmt.Printf("%s : %s\n", kv.Key, kv.Value)
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()
	serialResp, _ := cli.Get(ctx2, "my-key", clientv3.WithSerializable())
	fmt.Println("Serializable Read:")
	for _, kv := range serialResp.Kvs {
		fmt.Printf("%s : %s\n", kv.Key, kv.Value)
	}
}
