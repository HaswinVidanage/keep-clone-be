//+build wireinject

package wire

import (
	"github.com/google/wire"
	"hackernews-api/internal/config"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"hackernews-api/repositories"
	"hackernews-api/services/auth"
	"hackernews-api/services/note"
	"hackernews-api/services/user_config"
	"hackernews-api/services/users"
)

type App struct {
	DbProvider        *database.DbProvider
	UserService       *users.UserService
	NoteService       *note.NoteService
	AuthService       *auth.AuthService
	UserConfigService *user_config.UserConfigService
}

var dbSet = wire.NewSet(
	database.InitDB,
	wire.Bind(new(database.IDbProvider), new(*database.DbProvider)),
)

var configSet = wire.NewSet(
	config.GetCfg,
	wire.Bind(new(database.IDbConfig), new(*config.Config)),
)

var repositorySet = wire.NewSet(
	repositories.NewUserRepository,
)

var serviceSet = wire.NewSet(
	users.NewUserService,
	auth.NewAuthService,
	note.NewNoteService,
	user_config.NewUserConfigService,
)

func GetApp() (*App, error) {
	panic(wire.Build(
		configSet,
		dbSet,
		repositorySet,
		serviceSet,
		wire.Struct(new(App), "*"),
	))

	return &App{}, nil
}
