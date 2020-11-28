//+build wireinject

package wire

import (
	"context"
	"github.com/google/wire"
	"hackernews-api/graph"
	"hackernews-api/internal/config"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"hackernews-api/repositories"
	"hackernews-api/services/auth"
	"hackernews-api/services/note"
	"hackernews-api/services/user_config"
	"hackernews-api/services/users"
	"hackernews-api/test"
)

type App struct {
	Resolver          graph.Resolver
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
	repositories.NewAuthRepository,
	repositories.NewNoteRepository,
)

var serviceSet = wire.NewSet(
	users.NewUserService,
	auth.NewAuthService,
	note.NewNoteService,
	user_config.NewUserConfigService,
)

//var resolverSet = wire.NewSet(
//	graph.Resolver{
//		IUserService:       users.NewUserService,
//		INoteService:       note.NewNoteService,
//		IUserConfigService: user_config.NewUserConfigService,
//		IAuthService:       auth.NewAuthService,
//	},
//)

func GetApp() (*App, error) {
	panic(wire.Build(
		configSet,
		dbSet,
		repositorySet,
		serviceSet,
		wire.Struct(new(graph.Resolver), "*"),
		wire.Struct(new(App), "*"),
	))

	return &App{}, nil
}

func getTestContext() context.Context {
	return context.Background()
}

func GetTestApp() (*test.TestApp, error) {
	wire.Build(
		getTestContext,
		test.InitMockDB,
		repositorySet,
		serviceSet,
		wire.Struct(new(graph.Resolver), "*"),
		test.NewTestApp,
	)
	return &test.TestApp{}, nil
}
