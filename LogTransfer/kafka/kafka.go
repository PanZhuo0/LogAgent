package kafka

import (
	"GO/LogTransfer/es"
	"fmt"
	"github.com/IBM/sarama"
)

var consumer sarama.Consumer

func Init(address string) (err error) {
	consumer, err = sarama.NewConsumer([]string{address}, nil)
	if err != nil {
		fmt.Printf("failed to start consumer,err:%v\n", err)
		return
	}
	fmt.Println("Init Kafka succeed!")
	return
}

func Run(topic string) {
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Printf("failed to get list of partition:%v\n", err)
		return
	}
	fmt.Println("分区列表:", partitionList)
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition:%d,err:%v\n", partition, err)
			return
		}
		//defer pc.AsyncClose()
		//	异步从每个分区消费数据
		//	这里存在一个问题，就是这里的goroutine(内包?)就是这里的pc会关闭,这里应该不让pc关闭,或者保证该程序不结束
		// 	这里要保证只要这个goroutine在,就不能关闭这个pc，所以defer pc.AsyncClose() 不能执行，等GC回收
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("DataFromKafka ---- Topic:%s Pid:%d Offset:%d Key:%v Value:%v\n", topic, msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				es.SendToES(topic, msg.Value)
			}
		}(pc)
	}
}
