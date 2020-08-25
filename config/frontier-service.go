package config

import (
	"wchatv1/utils"
)

var (
	_                     utils.MicroServiceConfig = &frontierServiceConfig{}
	FrontierServiceConfig                          = &frontierServiceConfig{}
)

type frontierServiceConfig struct{}

func (c *frontierServiceConfig) ConfigKey() string {
	return "frontier-service"
}

func (c *frontierServiceConfig) ServiceName() string {
	return "micro.service.frontier-service"
}
