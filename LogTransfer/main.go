package main

import (
	"GO/LogTransfer/config"
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
	//	2.从kafka中取出数据
	//	3.发往ES
}
