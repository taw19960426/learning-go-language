package es

import (
	"context"
	"fmt"
	"strings"
	"time"

	//"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic"
)

//初始化ES，准备接收kafka那边发来的数据

type LogData struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

var (
	client *elastic.Client
	ch     chan *LogData
)

//init...
func Init(address string, chanSize int, worker int) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		return
	}
	fmt.Println("connect to es success")
	ch = make(chan *LogData, chanSize)

	//开辟 worker 个协程去读取Kafka数据发送至es
	for i := 0; i < worker; i++ {
		go SendToES()
	}

	return
}

func SendToESChan(msg *LogData) {
	ch <- msg
}

//发送数据到ES
func SendToES() {
	//链式操作
	for {
		select {
		case msg := <-ch:
			put1, err := client.Index().
				Index(msg.Topic). //Index表数据库
				Type("xxx").
				BodyJson(msg). //把一个go语言的对象转换为json格式
				Do(context.Background())
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Indexed %s to index %s,type %s\n", put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(time.Second)
		}
	}
}
