package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	microconfig "github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source"
	microsourceetcd "github.com/micro/go-micro/v2/config/source/etcd"
	"github.com/micro/go-micro/v2/registry"
	"time"
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
	Service               micro.Service
	Config                microconfig.Config
	Source                source.Source
	DefaultCliCallOptions client.CallOption
	etcd                  struct {
		addr     string
		registry registry.Registry
	}
}

func (m *microI) Init(service micro.Service) {
	if service == nil {
		service = micro.NewService()
		service.Init()
	}
	m.Service = service

	conf, err := microconfig.NewConfig()
	if err != nil {
		panic(err)
	}
	m.Config = conf
	m.Source = microsourceetcd.NewSource(microsourceetcd.WithAddress(service.Options().Registry.Options().Addrs[0]))

	m.DefaultCliCallOptions = func(options *client.CallOptions) {
		options.RequestTimeout = time.Second * 3
		options.DialTimeout = time.Second * 3
	}
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
func (m *microI) LoadConfigMust(c MicroConfig) {
	fmt.Printf("Loading Config: %s ...",c.ConfigKey())
	val := m.Config.Get("micro", "config", c.ConfigKey())
	if bytes.Equal(val.Bytes(), []byte("null")) {
		panic(fmt.Errorf("%s is null", c.ConfigKey()))
	}
	if err := json.Unmarshal(val.Bytes(), c); err != nil {
		panic(err)
	}
	fmt.Printf("ok\n")
}
