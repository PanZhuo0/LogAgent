package etcd

import (
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
