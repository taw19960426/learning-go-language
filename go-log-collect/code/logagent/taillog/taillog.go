package taillog

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"logagent/kafka"
)

// 专门收集日志的模块


type TailTask struct {
	path string
	topic string
	instance *tail.Tail
	// 为了能实现退出r,run()
	ctx context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(path, topic string) (t *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	t = &TailTask{
		path: path,
		topic: topic,
		ctx: ctx,
		cancelFunc: cancel,
	}
	err := t.Init()
	if err != nil {
		fmt.Println("tail file failed, err:", err)
	}
	return
}

func (t TailTask) Init() (err error) {
	config := tail.Config{
		ReOpen:    true,						// 充新打开
		Follow:    true,						// 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},	// 从文件哪个地方开始读
		MustExist: false,							// 文件不存在不报错
		Poll:      true}
	t.instance, err = tail.TailFile(t.path, config)
	// 当goroutine执行的函数退出的时候,goriutine就退出了
	go t.run()  // 直接去采集日志发送到kafka
	return
}

func (t *TailTask) run()  {
	for  {
		select {
		case <- t.ctx.Done():
			fmt.Printf("tail task:%s_%s 结束了...\n", t.path, t.topic)
			return
		case line :=<- t.instance.Lines:   // 从TailTask的通道中一行一行的读取日志
			// 3.2 发往kafka
			fmt.Printf("get log data from %s success, log:%v\n", t.path, line.Text)
			kafka.SendToChan(t.topic, line.Text)
		}
	}
}
