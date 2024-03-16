package config

type LogTransfer struct {
	Kafka `ini:"kafka"`
	ES    `ini:"es"`
}

type Kafka struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type ES struct {
	Address string `ini:"address"`
}
