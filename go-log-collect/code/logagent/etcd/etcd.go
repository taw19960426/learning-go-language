package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"strings"
	"time"
)

var (
	cli *clientv3.Client
)

type LogEntry struct {
	Path  string `json:"path"`  // 日志存放的路径
	Topic string `json:"topic"` // 日志发往kafka中的哪个Topic
}

// 初始化etcd的函数

func Init(addr string, timeout time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(addr, ";"),
		DialTimeout: timeout,
	})
	return
}

// 从etcd中获取日志收集项的配置信息
func GetConf(key string) (logEntryConf []*LogEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		err = json.Unmarshal(ev.Value, &logEntryConf)
		if err != nil {
			fmt.Printf("unmarshal etcd value failed,err:%v\n", err)
			return
		}
	}
	return
}

// etcd watch
func WatchConf(key string, newConfChan chan<- []*LogEntry) {
	rch := cli.Watch(context.Background(), key) // <-chan WatchResponse
	// 从通道尝试取值（监视的信息）
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			//	通知taillog.taskMgr
			var newConf []*LogEntry
			//1. 先判断操作的类型
			if ev.Type != clientv3.EventTypeDelete {
				// 如果不是是删除操作
				err := json.Unmarshal(ev.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("unmarshal failed, err:%v\n", err)
					continue
				}
			}

			fmt.Printf("get new conf: %v\n", newConf)
			newConfChan <- newConf
		}
	}
}
