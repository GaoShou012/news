package config

import (
	"im/utils"
)

var (
	_                  utils.MicroConfig = &kafkaClusterConfig{}
	KafkaClusterConfig                   = &kafkaClusterConfig{}
)

type kafkaClusterConfig struct {
	Addr []string
}

func (c *kafkaClusterConfig) ConfigKey() string {
	return "kafka-cluster"
}
