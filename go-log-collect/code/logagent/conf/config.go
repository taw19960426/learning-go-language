package conf

type Config struct {
	Kafka Kafka `ini:"kafka"`
	Etcd  Etcd  `ini:"etcd"`
}

type Kafka struct {
	Address string `ini:"address"`
	ChanMaxSize int `ini:"chan_max_zise"`
}

type Etcd struct {
	Address string `ini:"address"`
	Key     string `ini:"collect_log_key"`
	Timeout int    `ini:"timeout"`
}

//type TailLog struct {
//	FileName string `ini:"filename"`
//}
