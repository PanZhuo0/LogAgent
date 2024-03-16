package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

// kafka 消费者 实例
func main() {
	topic := "mysql_log"
	partition := 0
	conn, err := kafka.DialLeader(context.Background(),
		"tcp", "localhost:9092",
		topic, partition)
	if err != nil {
		fmt.Println("Failed to dial leader:", err)
		return
	}

	//	设置读取超时时间
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	//	读取一批消息,得到的batch是一系列消息的迭代器
	batch := conn.ReadBatch(10e3, 1e6) //最小10KB， 最大1M

	//	遍历读取信息
	b := make([]byte, 10e3) //每次最多10KB数据
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

	//	关闭batch
	err = batch.Close()
	if err != nil {
		fmt.Println("Failed to close batch:", err)
	}
}
