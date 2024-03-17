package main

import (
	"GO/LogTransfer/config"
	"GO/LogTransfer/es"
	"GO/LogTransfer/kafka"
	"fmt"
	"gopkg.in/ini.v1"
)

// logTransfer
// 将日志数据从kafka取出并发往ES
func main() {
	//	0.加载配置文件
	var cfg config.LogTransfer
	err := ini.MapTo(&cfg, "./config/config.ini")
	if err != nil {
		fmt.Println("Init config failed,err:", err)
		return
	}
	fmt.Println(cfg)
	//	1.初始化
	err = kafka.Init(cfg.Kafka.Address)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = es.Init(cfg.ES.Address, cfg.ES.ChanSize, cfg.ES.Nums)
	if err != nil {
		fmt.Println(err)
		return
	}
	//	2.从kafka中取出数据
	//	3.发往ES kafka.Run中包含了把数据发送给ES
	kafka.Run(cfg.Topic)
	select {}
}
