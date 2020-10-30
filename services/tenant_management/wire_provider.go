package tenant_management

import "github.com/google/wire"

var Provider = wire.NewSet(
	wire.Struct(new(HttpService), "*"),
)
