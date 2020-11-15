package users

import (
	"context"
	_ "database/sql"
	"github.com/google/wire"
	"hackernews-api/entities"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"hackernews-api/repositories"
	"hackernews-api/services/auth"
	"log"
)

type IUserService interface {
	Create(ctx context.Context, user entities.User)
	GetUserIdByName(ctx context.Context, name string) (int, error)
}

type UserService struct {
	DbProvider *database.DbProvider

	UserRepository repositories.IUserRepository
	AuthService    auth.IAuthService
}

var NewUserService = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)))

func (us *UserService) Create(ctx context.Context, user entities.User) {
	hashedPassword, err := us.AuthService.HashPassword(user.Password)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = hashedPassword

	_, err = us.UserRepository.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}

//GetUserIdByUsername check if a user exists in database by given username
func (us *UserService) GetUserIdByName(ctx context.Context, name string) (int, error) {
	id, err := us.UserRepository.GetUserIdByName(ctx, name)
	if err != nil {
		log.Fatal(err)
	}

	return id, err
}
