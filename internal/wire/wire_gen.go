// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package wire

import (
	"github.com/google/wire"
	"hackernews-api/graph"
	"hackernews-api/internal/config"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"hackernews-api/services/auth"
	"hackernews-api/services/note"
	"hackernews-api/services/user_config"
	"hackernews-api/services/users"
)

// Injectors from wire.go:

func GetApp() (*App, error) {
	configConfig := config.GetCfg()
	dbProvider := database.InitDB(configConfig)
	userService := &users.UserService{
		DbProvider: dbProvider,
	}
	noteService := &note.NoteService{
		DbProvider: dbProvider,
	}
	userConfigService := &user_config.UserConfigService{
		DbProvider: dbProvider,
	}
	resolver := &graph.Resolver{
		IUserService:       userService,
		INoteService:       noteService,
		IUserConfigService: userConfigService,
	}
	app := &App{
		Resolver:   resolver,
		DbProvider: dbProvider,
	}
	return app, nil
}

// wire.go:

type App struct {
	Resolver   *graph.Resolver
	DbProvider *database.DbProvider
}

var dbSet = wire.NewSet(database.InitDB, wire.Bind(new(database.IDbProvider), new(*database.DbProvider)))

var configSet = wire.NewSet(config.GetCfg, wire.Bind(new(database.IDbConfig), new(*config.Config)))

var serviceSet = wire.NewSet(users.NewUserService, auth.NewAuthService, note.NewNoteService, user_config.NewUserConfigService)
