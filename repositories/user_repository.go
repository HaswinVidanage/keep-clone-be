package repositories

import (
	"context"
	"database/sql"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"hackernews-api/entities"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"log"
)

type IUserRepository interface {
	InsertUser(context.Context, entities.User) (int64, error)
	GetUserIdByEmail(context.Context, string) (int, error)
}

type UserRepository struct {
	DbProvider *database.DbProvider
}

var NewUserRepository = wire.NewSet(
	wire.Struct(new(UserRepository), "*"),
	wire.Bind(new(IUserRepository), new(*UserRepository)))

func (ur *UserRepository) InsertUser(ctx context.Context, user entities.User) (int64, error) {
	statement, err := ur.DbProvider.Db.Prepare("insert into user( name, email, password) values(?,?,?)")
	print(statement)
	if err != nil {
		logrus.WithError(err).Warn(err)
		return 0, err
	}

	result, err := statement.Exec(user.Name, user.Email, user.Password)
	if err != nil {
		logrus.WithError(err).Warn(err)
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		logrus.WithError(err).Warn(err)
		return 0, err
	}

	return lastId, nil
}

func (ur *UserRepository) GetUserIdByEmail(ctx context.Context, email string) (int, error) {
	statement, err := ur.DbProvider.Db.Prepare("select id from user WHERE email = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(email)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return Id, nil
}
