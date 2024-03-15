package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
)

var producer sarama.SyncProducer

func Init(address string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	producer, err = sarama.NewSyncProducer([]string{address}, nil)
	if err != nil {
		return
	}
	fmt.Println("Init Kafka Succeed!")
	return
}
