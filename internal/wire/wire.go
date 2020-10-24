//+build wireinject

package wire

import (
	"github.com/google/wire"
	"hackernews-api/internal/config"
	"hackernews-api/internal/links"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"hackernews-api/internal/users"
)

type App struct {
	ConnectionProvider *database.ConnectionProvider
	UserService        *users.UserService
	LinkService        *links.LinkService
}

var dbSet = wire.NewSet(
	database.InitDB,
	wire.Bind(new(database.IConnectionProvider), new(*database.ConnectionProvider)),
)

var configSet = wire.NewSet(
	config.GetCfg,
	wire.Bind(new(database.IDbConfig), new(*config.Config)),
)

var serviceSet = wire.NewSet(
	users.NewUserService,
	links.NewLinkService,
)

func GetApp() (*App, error) {
	panic(wire.Build(
		configSet,
		dbSet,
		serviceSet,
		wire.Struct(new(App), "*"),
	))

	return &App{}, nil
}
