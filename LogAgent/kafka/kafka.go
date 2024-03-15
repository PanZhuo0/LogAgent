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

func SendToKafka(data, topic string) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	//	发送
	pid, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Send msg failed,err:", err)
		return
	}
	fmt.Printf("pid:%v offset:%v topic:%s msg:%s\n", pid, offset, topic, data)
}
