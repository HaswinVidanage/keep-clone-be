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
	GetUserByID(ctx context.Context, id int) (entities.User, error)
}

type UserService struct {
	DbProvider           *database.DbProvider
	UserConfigRepository repositories.IUserConfigRepository
	UserRepository       repositories.IUserRepository
	AuthService          auth.IAuthService
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

	fkUser, err := us.UserRepository.InsertUser(ctx, user)
	if err != nil {
		logrus.WithError(err).Warn(err)
		return "", err
	}

	// insert user config
	_, err = us.UserConfigRepository.InsertUserConfig(ctx, fkUser, entities.CreateUserConfig{
		IsDarkMode: false,
		IsListMode: false,
	})

	if err != nil {
		logrus.WithError(err).Warn(err)
		return "", err
	}

	token, err := jwt.GenerateToken(ctx, fkUser, user.Email)
	if err != nil {
		logrus.WithError(err).Warn(err)
		return "", err
	}
	return token, nil
}

func (us *UserService) GetUserByID(ctx context.Context, id int) (entities.User, error) {
	user, err := us.UserRepository.FindUserByID(ctx, id)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

//GetUserIdByUsername check if a user exists in database by given username
func (us *UserService) GetUserIdByEmail(ctx context.Context, email string) (int, error) {
	id, err := us.UserRepository.GetUserIdByEmail(ctx, email)
	if err != nil {
		log.Fatal(err)
	}

	return id, err
}
