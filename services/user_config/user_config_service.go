package user_config

import (
	"context"
	"errors"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"hackernews-api/entities"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"hackernews-api/repositories"
	"hackernews-api/services/auth"
)

type IUserConfigService interface {
	GetConfig(ctx context.Context) (entities.UserConfig, error)
	Save(ctx context.Context, uc entities.CreateUserConfig) (*int, error)
	Update(ctx context.Context, uc entities.UpdateUserConfig) (*int, error)
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
		logrus.Warn("user context is nil")
		return entities.UserConfig{}, errors.New("unauthorised")
	}

	uc, err := ucs.UserConfigRepository.FindUserConfigByUserID(ctx, userCtx.ID)

	if err != nil {
		logrus.WithError(err)
		return entities.UserConfig{}, err
	}

	return uc, nil
}

func (ucs UserConfigService) Save(ctx context.Context, uc entities.CreateUserConfig) (*int, error) {
	userCtx := auth.ForContext(ctx)
	if userCtx == nil {
		logrus.Warn("userCtx is nil")
		return nil, errors.New("unauthorised")
	}

	lastId, err := ucs.UserConfigRepository.InsertUserConfig(ctx, userCtx.ID, uc)

	if err != nil {
		return nil, err
	}
	return &lastId, nil
}

func (ucs UserConfigService) Update(ctx context.Context, uc entities.UpdateUserConfig) (*int, error) {
	userCtx := auth.ForContext(ctx)
	if userCtx == nil {
		logrus.Warn("userCtx is nil")
		return nil, errors.New("unauthorised")
	}

	// todo add fkUser and authorise all updates
	lastId, err := ucs.UserConfigRepository.UpdateUserConfigByFields(ctx, uc)

	if err != nil {
		return nil, err
	}
	return &lastId, nil
}
