package kafka

import (
	"fmt"
	"logtransfer/es"

	"github.com/Shopify/sarama"
)

//LogData...
type LogData struct {
	Data string `json:"data"`
}

//初始化kafka消费者，从kafka取数据发往ES
func Init(addr []string, topic string) (err error) {
	consumer, err := sarama.NewConsumer(addr, nil)
	if err != nil {
		fmt.Printf("fail to start consumer,err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(topic) //根据topic取到所有的分区
	if err != nil {
		fmt.Println("fail to get list of partition:", err)
		return
	}
	var pc sarama.PartitionConsumer
	fmt.Println(partitionList)
	for partition := range partitionList {
		pc, err = consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v", partition, err)
			return
		}
		defer pc.AsyncClose()
		//异步从每个分区消费消息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				//直接发给ES
				var ld = es.LogData{
					Topic: topic,
					Data:  string(msg.Value),
				}
				es.SendToESChan(&ld) //函数调函数
				//优化一下，直接放到chann中
			}
		}(pc)
		select {}
	}
	return
}
