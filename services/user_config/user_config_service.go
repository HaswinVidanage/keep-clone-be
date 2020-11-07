package user_config

import (
	"context"
	"github.com/google/wire"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"hackernews-api/services/auth"
	"hackernews-api/services/users"
	"log"
)

type UserConfig struct {
	ID         string
	IsDarkMode bool
	IsListMode bool
	User       *users.User
}

type IUserConfigService interface {
	GetConfig(ctx context.Context) *UserConfig
	Save(ctx context.Context, isDarkMode bool, isListMode bool, fkUser int) int64
}

type UserConfigService struct {
	DbProvider *database.DbProvider
}

var NewUserConfigService = wire.NewSet(
	wire.Struct(new(UserConfigService), "*"),
	wire.Bind(new(IUserConfigService), new(*UserConfigService)))

func (ucs UserConfigService) GetConfig(ctx context.Context) *UserConfig {
	userCtx := auth.ForContext(ctx)

	if userCtx == nil {
		log.Fatal("userCtx is nil")
	}

	stmt, err := ucs.DbProvider.Db.Prepare("select uc.id, uc.isDarkMode, uc.isListMode, uc.fk_user, u.Username from user_config uc inner join Users u on uc.fk_user = u.ID where uc.fk_user = ?")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()
	rows, err := stmt.Query(userCtx.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var userConfigs []UserConfig
	var username string
	var id string
	for rows.Next() {
		var userConfig UserConfig
		err := rows.Scan(&userConfig.ID, &userConfig.IsDarkMode, &userConfig.IsListMode, &id, &username)
		if err != nil {
			log.Fatal(err)
		}

		userConfig.User = &users.User{
			ID:       id,
			Username: username,
		}

		userConfigs = append(userConfigs, userConfig)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	count := len(userConfigs)
	if count > 1 {
		log.Fatal("user config returned more than 1 record for fk_user : " + id)
	} else if count == 0 {
		return nil
	}

	return &userConfigs[0]
}

func (ucs UserConfigService) Save(ctx context.Context, isDarkMode bool, isListMode bool, fkUser int) int64 {
	stmt, err := ucs.DbProvider.Db.Prepare("INSERT INTO user_config (isDarkMode, isListMode, fk_user) VALUES (?,?,?)")
	if err != nil {
		log.Fatal(err)
	}

	//  execution of our sql statement.
	res, err := stmt.Exec(isDarkMode, isListMode, fkUser)
	if err != nil {
		log.Fatal(err)
	}

	// retrieving Id of inserted Link.
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}

	log.Print("Row inserted!")
	return id
}
