package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"hackernews-api/repositories"
	"hackernews-api/services/auth"
	"hackernews-api/services/note"
	"hackernews-api/services/user_config"
	"hackernews-api/services/users"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	users.IUserService
	note.INoteService
	auth.IAuthService
	user_config.IUserConfigService
	repositories.IUserRepository
	repositories.INoteRepository
	repositories.IUserConfigRepository
}
