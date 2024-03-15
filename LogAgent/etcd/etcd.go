package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var client *clientv3.Client

func Init(address string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{address},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		return
	}
	fmt.Println("Prepare Etcd succeed!")
	return
}

func GetConf(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	response, err := client.Get(ctx, key)
	if err != nil {
		fmt.Println("Get Configuration from Etcd failed,err:", err)
		panic(err)
	}
	for _, v := range response.Kvs {
		fmt.Println(string(v.Value))
	}
	cancel()
}
