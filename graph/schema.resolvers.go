package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"hackernews-api/graph/generated"
	"hackernews-api/graph/model"
	"hackernews-api/internal/pkg/jwt"
	"hackernews-api/services/auth"
	"hackernews-api/services/links"
	"hackernews-api/services/note"
	"hackernews-api/services/users"
	"math/rand"
	"strconv"
)

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Link{}, fmt.Errorf("access denied")
	}

	var link links.Link
	link.Title = input.Title
	link.Address = input.Address
	link.User = user
	linkID := r.Resolver.ILinkService.Save(link)

	newModel := &model.Link{
		ID:      strconv.FormatInt(linkID, 10),
		Title:   link.Title,
		Address: link.Address,
	}

	//add new chanel in observer
	for _, observer := range addLinkObserver {
		observer <- newModel
	}

	return newModel, nil
}

func (r *mutationResolver) CreateNote(ctx context.Context, input model.NewNote) (*model.Note, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Note{}, fmt.Errorf("access denied")
	}

	var note note.Note
	note.Title = input.Title
	note.Content = input.Content
	note.User = user
	noteID := r.Resolver.INoteService.Save(note)

	newModel := &model.Note{
		ID:      strconv.FormatInt(noteID, 10),
		Title:   note.Title,
		Content: note.Content,
	}

	//add new chanel in observer
	for _, observer := range addNoteObserver {
		observer <- newModel
	}

	return newModel, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	r.IUserService.Create(ctx, user)
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	correct := r.IUserService.Authenticate(ctx, user)
	if !correct {
		// 1
		return "", &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) CreateUserConfig(ctx context.Context, input model.NewUserConfig) (int, error) {
	configId := r.Resolver.IUserConfigService.Save(ctx, input.IsDarkMode, input.IsListMode, input.FkUser)
	return int(configId), nil
}

func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	var resultLinks []*model.Link
	var dbLinks []links.Link
	dbLinks = r.Resolver.ILinkService.GetAll()
	for _, link := range dbLinks {
		grahpqlUser := &model.User{
			ID:   link.User.ID,
			Name: link.User.Username,
		}
		resultLinks = append(resultLinks, &model.Link{ID: link.ID, Title: link.Title, Address: link.Address, User: grahpqlUser})
	}
	return resultLinks, nil
}

func (r *queryResolver) Notes(ctx context.Context) ([]*model.Note, error) {
	var resultNotes []*model.Note
	var dbNotes []note.Note
	dbNotes = r.Resolver.INoteService.GetAll()
	for _, note := range dbNotes {
		grahpqlUser := &model.User{
			ID:   note.User.ID,
			Name: note.User.Username,
		}
		resultNotes = append(resultNotes, &model.Note{ID: note.ID, Title: note.Title, Content: note.Content, User: grahpqlUser})
	}
	return resultNotes, nil
}

func (r *queryResolver) UserConfig(ctx context.Context) (*model.UserConfig, error) {
	uc := r.Resolver.IUserConfigService.GetConfig(ctx)

	if uc == nil {
		return &model.UserConfig{}, nil
	}

	dbUC := &model.UserConfig{
		ID:         uc.ID,
		IsListMode: uc.IsListMode,
		IsDarkMode: uc.IsDarkMode,
		User: &model.User{
			ID:   uc.User.ID,
			Name: uc.User.Username,
		},
	}

	return dbUC, nil
}

func (r *subscriptionResolver) SubscriptionLinkAdded(ctx context.Context) (<-chan *model.Link, error) {
	id := randString(8)
	fmt.Println("Random id: ", id)
	events := make(chan *model.Link, 1)

	go func() {
		<-ctx.Done()
		delete(addLinkObserver, id)
	}()

	addLinkObserver[id] = events
	return events, nil
}

func (r *subscriptionResolver) SubscriptionNoteAdded(ctx context.Context) (<-chan *model.Note, error) {
	id := randString(8)
	fmt.Println("Random id: ", id)
	events := make(chan *model.Note, 1)

	go func() {
		<-ctx.Done()
		delete(addNoteObserver, id)
	}()

	addNoteObserver[id] = events
	return events, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
var addLinkObserver map[string]chan *model.Link
var addNoteObserver map[string]chan *model.Note

func init() {
	addLinkObserver = map[string]chan *model.Link{}
	addNoteObserver = map[string]chan *model.Note{}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
