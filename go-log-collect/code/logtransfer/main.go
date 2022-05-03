package main

import (
	"fmt"
	"logtransfer/conf"
	"logtransfer/es"
	"logtransfer/kafka"

	"gopkg.in/ini.v1"
)

//log transfer
//将日志数据从kafka取出来发往ES

func main() {
	//0.加载配置文件
	var cfg conf.Logtransfer
	err := ini.MapTo(&cfg, "./conf/cfg.ini")
	if err != nil {
		fmt.Println("init config err:", err)
		return
	}
	fmt.Printf("cfg:%v\n", cfg)
	//1.初始化ES
	//1.1 初始化一个ES连接的client
	//1.2 对外提供y一个往ES写入数据的一个函数
	err = es.Init(cfg.ESCfg.Address, cfg.ESCfg.ChanSize, cfg.ESCfg.Worker)
	if err != nil {
		fmt.Println("init ES consumer failed,err:", err)
		return
	}
	fmt.Println("init es success.")
	//2.初始化Kafka
	//2.1 连接kafka，创建分区的消费者
	//2.2 每个分区的消费者分别取出数据，通过sendToES（）将数据发往ES
	err = kafka.Init([]string{cfg.KafkaCfg.Address}, cfg.KafkaCfg.Topic)
	if err != nil {
		fmt.Println("init kafka consumer failed,err:", err)
		return
	}
	//3.从kafka中取数据

	//4.发往ES
	select {}
}
