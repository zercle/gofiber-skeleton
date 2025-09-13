package providers

import (
	"go.uber.org/fx"
)

// Module exports dependency injection modules.
var Module = fx.Options(
	InfrastructureModule,
	AuthModule,
	PostModule,
)
