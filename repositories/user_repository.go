package repositories

import (
	"context"
	"database/sql"
	sq "github.com/elgris/sqrl"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"hackernews-api/entities"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"log"
)

type IUserRepository interface {
	InsertUser(context.Context, entities.CreateUser) (int, error)
	GetUserIdByEmail(context.Context, string) (int, error)
	UpdateUserByFields(ctx context.Context, user entities.UpdateUser) (int, error)
	FindUserByID(ctx context.Context, id int) (entities.User, error)
	FindUserByEmail(ctx context.Context, email string) (entities.User, error)
	// DeleteUser
	// DeleteUserByID
	// FindAllUser
}

type UserRepository struct {
	DbProvider *database.DbProvider
}

var NewUserRepository = wire.NewSet(
	wire.Struct(new(UserRepository), "*"),
	wire.Bind(new(IUserRepository), new(*UserRepository)))

func (ur *UserRepository) InsertUser(ctx context.Context, user entities.CreateUser) (int, error) {
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

	return int(lastId), nil
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

func (ur *UserRepository) UpdateUserByFields(ctx context.Context, u entities.UpdateUser) (int, error) {
	updateMap := map[string]interface{}{}
	if u.Name != nil {
		updateMap["`name`"] = *u.Name
	}
	if u.Email != nil {
		updateMap["`email`"] = *u.Email
	}

	qb := sq.Update("`user`").SetMap(updateMap).Where(sq.Eq{"`id`": u.ID})
	// run query
	result, err := qb.RunWith(ur.DbProvider.Db).Exec()
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

func (ur *UserRepository) FindUserByID(ctx context.Context, id int) (entities.User, error) {
	statement, err := ur.DbProvider.Db.Prepare("select u.id, u.name, u.email from user u WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(id)
	var user entities.User
	err = row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return entities.User{}, err
	}
	return user, nil
}

func (ur *UserRepository) FindUserByEmail(ctx context.Context, email string) (entities.User, error) {
	statement, err := ur.DbProvider.Db.Prepare("select u.id, u.name, u.email from user u WHERE email = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(email)
	var user entities.User
	err = row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return entities.User{}, err
	}
	return user, nil
}
