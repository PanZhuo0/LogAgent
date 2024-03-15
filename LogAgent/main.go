package main

import (
	"GO/LogAgent/config"
	"GO/LogAgent/etcd"
	"GO/LogAgent/kafka"
	"fmt"
)

/*
项目背景:日志信息太多，无法在服务器上查看
解决方案:

	把机器上的日志实时收集，统一存储到中心系统
	对日志简历索引，便于查找
	提供一个界面友好的web页面实现日志展示与检索

使用到的工具:ES,KAFKA,Etcd,KIBANA
已有方案:ELK存在的问题,手动配置,无法做到定制化服务
优化:使用Etcd来实现热更新,自己自定义一个logAgent
*/
func main() {
	//	1.加载Ini配置文件
	cfg := config.Init()
	fmt.Println(cfg.Etcd)
	fmt.Println(cfg.Kafka)
	//	2.初始化Etcd和Kafka
	err := etcd.Init(cfg.Etcd.Address)
	if err != nil {
		fmt.Println("Prepare Etcd failed,err:", err)
		return
	}

	err = kafka.Init(cfg.Kafka.Address)
	if err != nil {
		fmt.Println("Init Kafka failed,err:", err)
		return
	}
}
