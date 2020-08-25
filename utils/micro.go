package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	micro "github.com/micro/go-micro/v2"
	microconfig "github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source"
	microsourceetcd "github.com/micro/go-micro/v2/config/source/etcd"
	"github.com/micro/go-micro/v2/registry"
)

var Micro microI

type MicroConfig interface {
	ConfigKey() string
}
type MicroServiceConfig interface {
	MicroConfig
	ServiceName() string
}

type microI struct {
	Service micro.Service
	Config  microconfig.Config
	Source  source.Source

	etcd struct {
		addr     string
		registry registry.Registry
	}
}

func (m *microI) InitV2() {
	service := micro.NewService()
	service.Init()
	m.Service = service

	conf, err := microconfig.NewConfig()
	if err != nil {
		panic(err)
	}
	m.Config = conf
	m.Source = microsourceetcd.NewSource(microsourceetcd.WithAddress(service.Options().Registry.Options().Addrs[0]))
}

func (m *microI) Init(service micro.Service) {
	m.Service = service
	m.Config = service.Options().Config
	m.Source = microsourceetcd.NewSource(microsourceetcd.WithAddress(service.Options().Registry.Options().Addrs[0]))
	//fmt.Println("init etcd ", etcdAddr)
	//
	//m.etcd.addr = etcdAddr
	//m.etcd.registry = etcd.NewRegistry(
	//	registry.Addrs(m.etcd.addr),
	//)
	//
	//conf, err := microconfig.NewConfig()
	//if err != nil {
	//	panic(err)
	//}
	//m.Config = conf
	//
	//m.Source = microsourceetcd.NewSource(microsourceetcd.WithAddress(m.etcd.addr))
}

func (m *microI) InitV1() {
	fmt.Println("Init Micro")
	conf, err := microconfig.NewConfig()
	if err != nil {
		panic(err)
	}
	m.Config = conf
	m.Source = microsourceetcd.NewSource()
}

func (m *microI) GetEtcdRegistry() registry.Registry {
	return m.etcd.registry
}

func (m *microI) LoadSource() {
	fmt.Println("Load Source")
	if err := m.Config.Load(m.Source); err != nil {
		panic(err)
	}
}

func (m *microI) LoadConfig(c MicroConfig) error {
	fmt.Println("Loading Config", c.ConfigKey())
	val := m.Config.Get("micro", "config", c.ConfigKey())
	if bytes.Equal(val.Bytes(), []byte("null")) {
		return fmt.Errorf("%s is null", c.ConfigKey())
	}
	if err := json.Unmarshal(val.Bytes(), c); err != nil {
		return err
	}
	return nil
}
func (m *microI) LoadConfigMust(c MicroConfig){
	fmt.Println("Loading Config",c.ConfigKey())
	val := m.Config.Get("micro","config",c.ConfigKey())
	if bytes.Equal(val.Bytes(),[]byte("null")) {
		panic(fmt.Errorf("%s is null", c.ConfigKey()))
	}
	if err := json.Unmarshal(val.Bytes(),c); err != nil {
		panic(err)
	}
}
