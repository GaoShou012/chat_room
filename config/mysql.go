package config

import "wchatv1/utils"

var (
	_           utils.MicroConfig = &mysqlConfig{}
	MysqlConfig                   = &mysqlConfig{}
)

type mysqlConfig struct {
	Addr     string
	Port     int
	User     string
	Password string
	Database string
}

func (m mysqlConfig) ConfigKey() string {
	return "mysql"
}
