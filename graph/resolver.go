package graph

import "hackernews-api/internal/users"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	users.IUserService
}
