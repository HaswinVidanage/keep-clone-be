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
	// UpdateUserByFields
	UpdateUser(ctx context.Context, u entities.User) (int, error)
	// DeleteUser
	// DeleteUserByID
	// FindAllUser
	FindUserByID(ctx context.Context, id int) (entities.User, error)
	// UsersByEmail
	UpdateUserBaseQuery(ctx context.Context, u entities.User) *sq.UpdateBuilder
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

func (ur *UserRepository) UpdateUser(ctx context.Context, u entities.User) (int, error) {
	qb := ur.UpdateUserBaseQuery(ctx, u)
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
	statement, err := ur.DbProvider.Db.Prepare("select * from user WHERE id = ?")
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

// Get queries
func (ur *UserRepository) UpdateUserBaseQuery(ctx context.Context, u entities.User) *sq.UpdateBuilder {
	qb := sq.Update("`user`").SetMap(map[string]interface{}{
		"`name`":  u.Name,
		"`email`": u.Email,
	}).Where(sq.Eq{"`id`": u.ID})
	return qb
}
