package room

import "github.com/google/wire"

var Provider = wire.NewSet(
	wire.Struct(new(TenantRoomUser), "*"),
	wire.Struct(new(TenantRoomUserMap), "*"),
	wire.Struct(new(TenantRoomUserSortedSet), "*"),
)
