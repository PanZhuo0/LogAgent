package config

import "gopkg.in/ini.v1"

type Etcd struct {
	Address string `ini:"address"`
	Key     string `ini:"key"`
}

type Kafka struct {
	Address string `ini:"address"`
}

type Config struct {
	Etcd  `ini:"etcd"`
	Kafka `ini:"kafka"`
}

func Init() (cfg Config) {
	ini.MapTo(&cfg, "./config/config.ini")
	return
}
