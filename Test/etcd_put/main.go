package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net"
	"time"
)

func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	//获取数据
	key := fmt.Sprintf("/logAgent/%s/config", GetOutBindIP())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//除了个小错误，这里的topic 写成了topic:
	//value := `[{"path":"d:/mysql.log","topic:":"mysql_log"},{"path":"d:/redis.log","topic:":"redis_log"},{"path":"d:/nginx.log","topic:":"nginx_log"}]`
	//value := `[{"path":"d:/mysql.log","topic":"mysql_log"},{"path":"d:/redis.log","topic":"redis_log"},{"path":"d:/nginx.log","topic":"nginx_log"}]`
	//value := `[{"path":"d:/mysql.log","topic":"mysql_log"},{"path":"d:/redis.log","topic":"redis_log"}]`
	//value := `[{"path":"d:/mysql.log","topic":"mysql_log"}]`
	value := ``
	client.Put(ctx, key, value)
	cancel()
}

func GetOutBindIP() (ip string) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = localAddr.IP.String()
	fmt.Println("IP:", ip)
	return
}
