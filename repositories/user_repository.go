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
	GetUserIdByName(context.Context, string) (int, error)
}

type UserRepository struct {
	DbProvider *database.DbProvider
}

var NewUserRepository = wire.NewSet(
	wire.Struct(new(UserRepository), "*"),
	wire.Bind(new(IUserRepository), new(*UserRepository)))

func (ur *UserRepository) InsertUser(ctx context.Context, user entities.User) (int64, error) {
	statement, err := ur.DbProvider.Db.Prepare("insert into user(name,password) values(?,?)")
	print(statement)
	if err != nil {
		log.Fatal(err)
	}

	result, err := statement.Exec(user.Name, user.Password)
	if err != nil {
		logrus.WithError(err).Fatal(err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		logrus.WithError(err).Fatal(err)
	}

	return lastId, nil
}

func (ur *UserRepository) GetUserIdByName(ctx context.Context, name string) (int, error) {
	statement, err := ur.DbProvider.Db.Prepare("select id from user WHERE name = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(name)

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
