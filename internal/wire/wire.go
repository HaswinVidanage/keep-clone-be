//+build wireinject

package wire

import (
	"github.com/google/wire"
	"hackernews-api/internal/config"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"hackernews-api/services/auth"
	"hackernews-api/services/links"
	"hackernews-api/services/users"
)

type App struct {
	DbProvider     *database.DbProvider
	UserService    *users.UserService
	LinkService    *links.LinkService
	NewAuthService *auth.AuthService
}

var dbSet = wire.NewSet(
	database.InitDB,
	wire.Bind(new(database.IDbProvider), new(*database.DbProvider)),
)

var configSet = wire.NewSet(
	config.GetCfg,
	wire.Bind(new(database.IDbConfig), new(*config.Config)),
)

var serviceSet = wire.NewSet(
	users.NewUserService,
	links.NewLinkService,
	auth.NewAuthService,
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