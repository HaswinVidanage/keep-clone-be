package user_config

import (
	"context"
	"github.com/google/wire"
	"hackernews-api/entities"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"hackernews-api/repositories"
	"hackernews-api/services/auth"
	"log"
)

type IUserConfigService interface {
	GetConfig(ctx context.Context) (entities.UserConfig, error)
	Save(ctx context.Context, uc entities.CreateUserConfig) (*int, error)
}

type UserConfigService struct {
	DbProvider           *database.DbProvider
	UserConfigRepository repositories.UserConfigRepository
}

var NewUserConfigService = wire.NewSet(
	wire.Struct(new(UserConfigService), "*"),
	wire.Bind(new(IUserConfigService), new(*UserConfigService)))

func (ucs UserConfigService) GetConfig(ctx context.Context) (entities.UserConfig, error) {
	userCtx := auth.ForContext(ctx)

	if userCtx == nil {
		log.Fatal("userCtx is nil")
		// todo throw unauthorised error
	}

	uc, err := ucs.UserConfigRepository.FindUserConfigByUserID(ctx, userCtx.ID)

	if err != nil {
		return entities.UserConfig{}, err
	}

	return uc, nil
}

func (ucs UserConfigService) Save(ctx context.Context, uc entities.CreateUserConfig) (*int, error) {

	userCtx := auth.ForContext(ctx)
	if userCtx == nil {
		log.Fatal("userCtx is nil")
		// todo throw unauthorised error
	}

	lastId, err := ucs.UserConfigRepository.InsertUserConfig(ctx, userCtx.ID, uc)

	if err != nil {
		return nil, err
	}
	return &lastId, nil
}
