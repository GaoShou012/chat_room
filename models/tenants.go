package models

type Tenants struct {
	Model
	Enable     *bool
	TenantCode *string
	TenantKey  *string
}
