package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"hackernews-api/graph/generated"
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
	user_config.IUserConfigService
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
