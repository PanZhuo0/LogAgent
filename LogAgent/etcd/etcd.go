package etcd

import (
	"GO/LogAgent/taillog"
	"context"
	"encoding/json"
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
		//如果是空的配置,continue
		if v.Value == nil {
			fmt.Println("空的配置!")
			continue
		}
		err := json.Unmarshal(v.Value, &taillog.Mgr.LogEntries)
		if err != nil {
			fmt.Println("Json unmarshal failed,err:", err)
			panic(err)
		}
	}
	cancel()
}

// WatchConf 监视Etcd是否发生变化
func WatchConf(key string) {
	watchChan := client.Watch(context.Background(), key)
	for resp := range watchChan {
		for _, v := range resp.Events {
			//直接把新的修改发送给taillog中的Manager处理就行,这里不做处理
			//	这里的修改数据被断言是json类型LogEntry对象配置
			taillog.Mgr.NewConfCh <- v.Kv.Value
		}
	}
}
