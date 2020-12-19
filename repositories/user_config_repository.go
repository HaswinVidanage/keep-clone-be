package repositories

import (
	"context"
	"database/sql"
	"github.com/elgris/sqrl"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"hackernews-api/entities"
	"hackernews-api/internal/pkg/db/migrations/mysql"
)

type IUserConfigRepository interface {
	InsertUserConfig(ctx context.Context, fkUser int, userConfig entities.CreateUserConfig) (int, error)
	UpdateUserConfigByFields(ctx context.Context, user entities.UpdateUserConfig) (int, error)
	FindUserConfigByUserID(ctx context.Context, fkUser int) (entities.UserConfig, error)
}

type UserConfigRepository struct {
	DbProvider     *database.DbProvider
	UserRepository IUserRepository
}

var NewUserConfigRepository = wire.NewSet(
	wire.Struct(new(UserConfigRepository), "*"),
	wire.Bind(new(IUserConfigRepository), new(*UserConfigRepository)))

func (ucr *UserConfigRepository) InsertUserConfig(ctx context.Context, fkUser int, userConfig entities.CreateUserConfig) (int, error) {
	stmt, err := ucr.DbProvider.Db.Prepare("INSERT INTO user_config (isDarkMode, isListMode, fk_user) VALUES (?,?,?)")
	if err != nil {
		logrus.WithError(err).Warn(err)
		return 0, err
	}

	res, err := stmt.Exec(userConfig.IsDarkMode, userConfig.IsListMode, fkUser)
	if err != nil {
		logrus.WithError(err).Warn(err)
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		logrus.WithError(err).Warn(err)
		return 0, err
	}

	return int(lastId), nil
}

func (ucr *UserConfigRepository) UpdateUserConfigByFields(ctx context.Context, n entities.UpdateUserConfig) (int, error) {
	updateMap := map[string]interface{}{}
	if n.IsListMode != nil {
		updateMap["`isListMode`"] = *n.IsListMode
	}
	if n.IsDarkMode != nil {
		updateMap["`isDarkMode`"] = *n.IsDarkMode
	}

	qb := sqrl.Update("`user_config`").SetMap(updateMap).Where(sqrl.Eq{"`id`": n.ID})
	// run query
	result, err := qb.RunWith(ucr.DbProvider.Db).Exec()
	if err != nil {
		logrus.WithError(err).Warn(err)
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		logrus.WithError(err).Warn(err)
		return 0, err
	}

	return int(lastId), nil
}

func (ucr *UserConfigRepository) FindUserConfigByUserID(ctx context.Context, fkUser int) (entities.UserConfig, error) {
	stmt, err := ucr.DbProvider.Db.Prepare("select n.id, n.isDarkMode, n.isListMode from user_config n where n.fk_user = ?")
	if err != nil {
		logrus.WithError(err).Warn(err)
		return entities.UserConfig{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(fkUser)
	var uc entities.UserConfig
	err = row.Scan(&uc.ID, &uc.IsDarkMode, &uc.IsListMode)
	if err != nil {
		if err != sql.ErrNoRows {
			// todo handle error
			//logrus.WithError(err).Warn(err)
		}
		logrus.WithError(err).Warn(err)
		return entities.UserConfig{}, err
	}

	// find user
	user, err := ucr.UserRepository.FindUserByID(ctx, fkUser)
	if err != nil {
		if err != sql.ErrNoRows {
			// todo handle error
		}
		logrus.WithError(err).Warn(err)
		return entities.UserConfig{}, err
	}
	uc.User = &user

	return uc, nil
}
