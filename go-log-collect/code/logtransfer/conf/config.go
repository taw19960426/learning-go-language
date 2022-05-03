package conf

//LogTransfer 全局配置
type Logtransfer struct {
	KafkaCfg `ini:"kafka"`
	ESCfg    `ini:"es"`
}

//Kafka...
type KafkaCfg struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

//ESCfg
type ESCfg struct {
	Address  string `ini:"address"`
	ChanSize int    `ini:"size"`
	Worker   int    `ini:"worker"`
}
