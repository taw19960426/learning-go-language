package main

import (
	"fmt"
	"logagent/conf"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/taillog"
	"logagent/tools"
	"strings"
	"sync"
	"time"

	"gopkg.in/ini.v1"
)

var config = new(conf.Config)

// logAgent 入口程序

func main() {
	// 0. 加载配置文件
	err := ini.MapTo(config, "./conf/config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		return
	}
	// 1. 初始化kafka连接
	err = kafka.Init(strings.Split(config.Kafka.Address, ";"), config.Kafka.ChanMaxSize)
	if err != nil {
		fmt.Printf("init kafka failed, err:%v\n", err)
		return
	}
	fmt.Println("init kafka success.")

	// 2. 初始化etcd
	err = etcd.Init(config.Etcd.Address, time.Duration(config.Etcd.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("init etcd failed,err:%v\n", err)
		return
	}
	fmt.Println("init etcd success.")
	// 实现每个logagent都拉取自己独有的配置，所以要以自己的IP地址实现热加载
	ip, err := tools.GetOurboundIP()
	if err != nil {
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(config.Etcd.Key, ip)
	// 2.1 从etcd中获取日志收集项的配置信息
	logEntryConf, err := etcd.GetConf(etcdConfKey)
	if err != nil {
		fmt.Printf("etcd.GetConf failed, err:%v\n", err)
		return
	}
	fmt.Printf("get conf from etcd success, %v\n", logEntryConf)

	// 2.2 派一个哨兵 一直监视着 zhangyafei这个key的变化（新增 删除 修改））
	for index, value := range logEntryConf {
		fmt.Printf("index:%v value:%v\n", index, value)
	}
	// 3. 收集日志发往kafka
	taillog.Init(logEntryConf)

	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(etcdConfKey, taillog.NewConfChan()) // 哨兵发现最新的配置信息会通知上面的通道
	wg.Wait()
}
