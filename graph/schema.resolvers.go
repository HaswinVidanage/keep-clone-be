package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"hackernews-api/entities"
	"hackernews-api/graph/generated"
	"hackernews-api/graph/model"
	"hackernews-api/services/auth"
	"math/rand"
)

func (r *mutationResolver) CreateNote(ctx context.Context, input model.NewNote) (*model.Note, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Note{}, fmt.Errorf("access denied")
	}

	var note entities.Note
	note.Title = input.Title
	note.Content = input.Content
	note.User = user

	noteID := r.Resolver.INoteService.SaveNote(note)
	userDto, err := r.Resolver.IUserService.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	newModel := &model.Note{
		ID:      int(noteID),
		Title:   note.Title,
		Content: note.Content,
		// TODO move to service and get user by id
		User: &model.User{
			ID:    userDto.ID,
			Email: userDto.Email,
			Name:  userDto.Name,
		},
	}

	//add new chanel in observer
	for _, observer := range addNoteObserver {
		observer <- newModel
	}

	return newModel, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	token, err := r.IUserService.CreateUser(ctx, entities.CreateUser{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	token, err := r.IAuthService.Login(ctx, input.Email, input.Password)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	token, err := r.IAuthService.RefreshToken(ctx, input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	return token, nil
}

func (r *mutationResolver) CreateUserConfig(ctx context.Context, input model.NewUserConfig) (int, error) {
	configId := r.Resolver.IUserConfigService.Save(ctx, input.IsDarkMode, input.IsListMode, input.FkUser)
	return int(configId), nil
}

func (r *queryResolver) Notes(ctx context.Context) ([]*model.Note, error) {
	var resultNotes []*model.Note
	var dbNotes []entities.Note
	dbNotes = r.Resolver.INoteService.GetAll()
	for _, note := range dbNotes {
		grahpqlUser := &model.User{
			ID:    note.User.ID,
			Name:  note.User.Name,
			Email: note.User.Email,
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
			ID:    uc.User.ID,
			Name:  uc.User.Name,
			Email: uc.User.Email,
		},
	}

	return dbUC, nil
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
var addNoteObserver map[string]chan *model.Note

func init() {
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
