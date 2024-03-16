package main

import (
	"fmt"
	"github.com/IBM/sarama"
)

func main() {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Printf("failed to start consumer,err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions("mysql_log")
	if err != nil {
		fmt.Printf("failed to get list of partition:%v\n", err)
		return
	}
	fmt.Println("分区列表:", partitionList)
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("mysql_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition:%d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		//	异步从每个分区消费数据
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Pid:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
			}
		}(pc)
	}

	//保证go routine不停止
	select {}
}
