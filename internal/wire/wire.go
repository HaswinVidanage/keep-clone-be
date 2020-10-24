//+build wireinject

package wire

import (
	"github.com/google/wire"
	"hackernews-api/internal/config"
	"hackernews-api/internal/pkg/db/migrations/mysql"
)

type App struct {
	ConnectionProvider *database.ConnectionProvider
}

var dbSet = wire.NewSet(
	database.InitDB,
	wire.Bind(new(database.IConnectionProvider), new(*database.ConnectionProvider)),
)

var configSet = wire.NewSet(
	config.GetCfg,
	wire.Bind(new(database.IDbConfig), new(*config.Config)),
)

func GetApp() (*App, error) {
	panic(wire.Build(
		configSet,
		dbSet,
		wire.Struct(new(App), "*"),
	))

	return &App{}, nil
}
