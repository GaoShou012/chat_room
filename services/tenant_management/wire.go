//+build wireinject

package tenant_management

import (
	"github.com/go-redis/redis/v7"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	controller_tenant_management "wchatv1/controller/tenant_management"
	"wchatv1/src/room"
)

func NewHttpService(
	db *gorm.DB,
	redisClient *redis.ClusterClient,
) *HttpService {
	wire.Build(
		Provider,
		controller_tenant_management.Provider,
		room.Provider,
	)
	return nil
}
