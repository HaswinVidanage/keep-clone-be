package users

import (
	"context"
	_ "database/sql"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"hackernews-api/entities"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"hackernews-api/internal/pkg/jwt"
	"hackernews-api/repositories"
	"hackernews-api/services/auth"
	"log"
)

type IUserService interface {
	CreateUser(ctx context.Context, user entities.CreateUser) (string, error)
	GetUserIdByEmail(ctx context.Context, name string) (int, error)
}

type UserService struct {
	DbProvider *database.DbProvider

	UserRepository repositories.IUserRepository
	AuthService    auth.IAuthService
}

var NewUserService = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)))

func (us *UserService) CreateUser(ctx context.Context, user entities.CreateUser) (string, error) {
	hashedPassword, err := us.AuthService.HashPassword(user.Password)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = hashedPassword

	lastId, err := us.UserRepository.InsertUser(ctx, user)
	if err != nil {
		logrus.WithError(err).Warn(err)
		return "", err
	}

	token, err := jwt.GenerateToken(ctx, lastId, user.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}

//GetUserIdByUsername check if a user exists in database by given username
func (us *UserService) GetUserIdByEmail(ctx context.Context, email string) (int, error) {
	id, err := us.UserRepository.GetUserIdByEmail(ctx, email)
	if err != nil {
		log.Fatal(err)
	}

	return id, err
}
