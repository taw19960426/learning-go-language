package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// 获取本地对外IP
func GetOurboundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.IP.String(), ":")[0]
	return
}

func main() {
	// etcd client put/get demo
	// use etcd/clientv3
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	value := `[{"path":"F:/go/goCode/src/github.com/xiaoTang/goStudy06/06logAgentEtcd/logs/nginx.log","topic":"web_log"},{"path":"F:/go/goCode/src/github.com/xiaoTang/goStudy06/06logAgentEtcd/logs/redis.log","topic":"redis_log"},{"path":"F:/go/goCode/src/github.com/xiaoTang/goStudy06/06logAgentEtcd/logs/mysql.log","topic":"mysql_log"}]`
	//value := `[{"path":"f:/tmp/nginx.log","topic":"web_log"},{"path":"f:/tmp/redis.log","topic":"redis_log"}]`
	//_, err = cli.Put(ctx, "zhangyafei", "dsb")

	//初始化key
	ip, err := GetOurboundIP()
	if err != nil {
		panic(err)
	}
	log_conf_key := fmt.Sprintf("/logagent/%s/collect_config", ip)
	_, err = cli.Put(ctx, log_conf_key, value)

	//_, err = cli.Put(ctx, "/logagent/collect_config", value)
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)

	resp, err := cli.Get(ctx, log_conf_key)
	//resp, err := cli.Get(ctx, "/logagent/collect_config")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}
