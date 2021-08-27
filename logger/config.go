// Copyright © 2021 Alibaba Group Holding Ltd.

package logger

type Config struct {
	DebugMode bool
}

var loggerConfig Config

func InitLogger(config Config) {
	loggerConfig = config
}
