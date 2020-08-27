package config

import (
	"im/utils"
	"time"
)

var (
	_              utils.MicroConfig = &frontierConfig{}
	FrontierConfig                   = &frontierConfig{}
)

type frontierConfig struct {
	Debug            bool
	HeartbeatTimeout int64
	WriterTimeout    time.Duration
	ReaderTimeout    time.Duration
	AcceptProcNum    int
}

func (c *frontierConfig) ConfigKey() string {
	return "frontier"
}
