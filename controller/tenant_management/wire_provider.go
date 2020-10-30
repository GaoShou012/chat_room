package tenant_api

import "github.com/google/wire"

var Provider = wire.NewSet(
	wire.Struct(new(Auth), "*"),
	wire.Struct(new(Users), "*"),
)
