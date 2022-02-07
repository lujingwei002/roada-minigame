package tlogd

import (
	"github.com/roada-go/roada"
)

type Config struct {
	Dir       string
	Prefix    string
	BackupDir string
	LineLimit int64 //最大行数
	TimeLimit int64 //最大时间,分钟
	TcpAddr   string
	Console   bool
}

func initDefaultConfig(config *Config) {
	if config.LineLimit == 0 {
		config.LineLimit = 100000
	}
	if config.TimeLimit == 0 {
		config.TimeLimit = 30
	}
}

func Register(road *roada.Road, config *Config) error {
	initDefaultConfig(config)
	if err := newTLogService(road, config); err != nil {
		return err
	}
	return nil
}
