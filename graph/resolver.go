package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"hackernews-api/services/links"
	"hackernews-api/services/note"
	"hackernews-api/services/user_config"
	"hackernews-api/services/users"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	users.IUserService
	links.ILinkService
	note.INoteService
	user_config.IUserConfigService
}
