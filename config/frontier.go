package config

import (
	"wchatv1/utils"
)

var (
	_              utils.MicroConfig = &frontierConfig{}
	FrontierConfig                   = &frontierConfig{}
)

type frontierConfig struct {
	LogLevel         int   `json:"logLevel"`
	HeartbeatTimeout int64 `json:"heartbeatTimeout"`
	WriterBufferSize int   `json:"writerBufferSize"`
	ReaderBufferSize int   `json:"readerBufferSize"`
}

func (c *frontierConfig) ConfigKey() string {
	return "frontier"
}
