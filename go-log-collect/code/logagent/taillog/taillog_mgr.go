package taillog

import (
	"fmt"
	"logagent/etcd"
	"time"
)

var taskMrg *TailLogMgr

type TailLogMgr struct {
	logEntry []*etcd.LogEntry
	taskMap map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func Init(logEntryConf []*etcd.LogEntry)  {
	taskMrg = &TailLogMgr{
		logEntry: logEntryConf,
		taskMap: make(map[string]*TailTask, 16),
		newConfChan: make(chan  []*etcd.LogEntry),  // 无缓冲区的通道
	}
	for _, logEntry := range logEntryConf{
		// 3.1 循环每一个日志收集项，创建TailObj
		// logEntry.Path  要收集的全日志文件的路径
		// 初始化的时候齐了多少个tailTask 都要记下来，为了后续判断方便
		tailObj := NewTailTask(logEntry.Path, logEntry.Topic)
		mk := fmt.Sprintf("%s_%s", logEntry.Path, logEntry.Topic)
		taskMrg.taskMap[mk] = tailObj
	}
	go taskMrg.run()
}

// 监听自己的newConfChan,有了新的配合过来之后就做对应的处理

func (t *TailLogMgr) run()  {
	for  {
		select {
		case newConf := <- t.newConfChan:
			// 1. 配置新增
			for _, conf := range newConf {
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := t.taskMap[mk]
				if ok {
					// 原来就有，不需要操作
					continue
				}else {
					// 新增的
					tailObj := NewTailTask(conf.Path, conf.Topic)
					t.taskMap[mk] = tailObj
				}
			}
			// 找出原来t.logEntry有，但是newConf中没有的，删掉
			for _, c1 := range t.logEntry{  // 循环原来的配置
				isDelete := true
				for _, c2 := range newConf{  // 取出新的配置
					if c2.Path == c1.Path && c2.Topic == c1.Topic {
						isDelete = false
						continue
					}
				}
				if isDelete {
					// 把c1对应的这个tailObj给停掉
					mk := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					// t.taskNap[mk] ==> tailObj
					t.taskMap[mk].cancelFunc()
				}
			}
			// 2. 配置删除
			// 3. 配置变更
			fmt.Println("新的配置来了！", newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}

// 一个函数，向外暴露taskMgr的newConfChan
func NewConfChan() chan <-[]*etcd.LogEntry {
	return taskMrg.newConfChan
}