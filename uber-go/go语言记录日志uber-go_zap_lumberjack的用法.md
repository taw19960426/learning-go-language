## 0 需求分析

 - 日志切割能够根据文件大小、时间或间隔等来切割日志文件；
 - 支持不同的日志级别，例如 DEBUG ， INFO ， WARN ， ERROR 等；
 - 能够打印基本信息，如调用文件、函数名和行号，日志时间等；
 - 根据时间或者天数来保存日志信息

## 1 环境安装

```go
go get -u go.uber.org/zap
go get -v github.com/uber-go/atomic
go get -v github.com/uber-go/multierr
go get -uv github.com/natefinch/lumberjack
```
如果安装失败，就将我GitHub上的安装包**go.uber.org**解压后拷贝到 **GOPATH/src** 下，整个测试用例也在github上
[GitHub地址](https://github.com/taw19960426/learning-go-language/tree/main/uber-go)

```cpp
https://github.com/taw19960426/learning-go-language/tree/main/uber-go
```

```go
在这里插入代码片
```
2 参考博客
[这里面写得很详细](https://blog.csdn.net/wohu1104/article/details/107326794)
```go
https://blog.csdn.net/wohu1104/article/details/107326794
```

## 3 结果展示
![在这里插入图片描述](https://img-blog.csdnimg.cn/bb488a9e053e4e61b7b0ff589375fa1b.png?x-oss-process=image/watermark,type_d3F5LXplbmhlaQ,shadow_50,text_Q1NETiBA5ZSQ57u05bq3,size_20,color_FFFFFF,t_70,g_se,x_16)

## 4 源代码

main.go

```go
package main

import (
	"time"

	"go.uber.org/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// logpath 日志文件路径
// loglevel 日志级别
func InitLogger(logpath string, loglevel string) {
	// 日志分割
	hook := lumberjack.Logger{
		Filename:   logpath, // 日志文件路径，默认 os.TempDir()
		MaxSize:    1,       // 每个日志文件保存1M，默认 100M
		MaxBackups: 30,      // 保留30个备份，默认不限
		MaxAge:     7,       // 保留7天，默认不限
		Compress:   true,    // 是否压缩，默认不压缩
	}
	write := zapcore.AddSync(&hook)
	// 设置日志级别
	// debug 可以打印出 info debug warn
	// info  级别可以打印 warn info
	// warn  只能打印 warn
	// debug->info->warn->error
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	core := zapcore.NewCore(
		// zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewJSONEncoder(encoderConfig),
		// zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&write)), // 打印到控制台和文件
		write,
		level,
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段,如：添加一个服务器名称
	filed := zap.Fields(zap.String("serviceName", "192.168.1.199"))
	// 构造日志 如果不需要一些参数可以删除
	logger = zap.New(core, caller, development, filed)
	//logger = zap.New(core, development)
	logger.Info("DefaultLogger init success")
}

func main() {
	// 历史记录日志名字为：my.log，服务重新启动，日志会追加，不会删除
	InitLogger("./logs/my.log", "debug")
	// 强结构形式
	logger.Info("test",
		zap.String("string", "xiaotang"),
		zap.Int("int", 3),
		zap.Duration("time", time.Second),
	)

	// // 必须 key-value 结构形式 性能下降一点
	// logger.Sugar().Infow("test-",
	// 	"string", "kk",
	// 	"int", 1,
	// 	"time", time.Second,
	// )

	logger.Error("test02",
		zap.String("string", "x666g"),
		zap.Int("int", 4),
		zap.Duration("time", time.Second),
	)

	for {
		logger.Error("test02",
			zap.String("string", "x666g"),
			zap.Int("int", 4),
			zap.Duration("time", time.Second),
		)
	}
}

```

