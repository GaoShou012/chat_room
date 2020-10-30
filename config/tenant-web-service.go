package config

import "wchatv1/utils"

var (
	_                utils.MicroServiceConfig = &tenantWebService{}
	TenantWebService                          = &tenantWebService{}
)

type tenantWebService struct {}

func (t tenantWebService) ConfigKey() string {
	return "tenant-web-service"
}

func (t tenantWebService) ServiceName() string {
	return "tenant-web-service"
}
