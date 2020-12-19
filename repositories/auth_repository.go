package repositories

import (
	"context"
	"database/sql"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"log"
)

type IAuthRepository interface {
	GetHashedPasswordByUserEmail(context.Context, string) (string, error)
	GetUserIdByEmail(context.Context, string) (int, error)
}

type AuthRepository struct {
	DbProvider *database.DbProvider
}

var NewAuthRepository = wire.NewSet(
	wire.Struct(new(AuthRepository), "*"),
	wire.Bind(new(IAuthRepository), new(*AuthRepository)),
)

func (ar *AuthRepository) GetHashedPasswordByUserEmail(ctx context.Context, email string) (string, error) {
	statement, err := ar.DbProvider.Db.Prepare("select password from user WHERE email = ?")
	if err != nil {
		log.Print(err)
		return "", err
	}
	row := statement.QueryRow(email)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		logrus.WithError(err).Warn(err)
		if err == sql.ErrNoRows {
			return "", err
		} else {
			return "", err
		}
	}

	return hashedPassword, nil
}

func (ar *AuthRepository) GetUserIdByEmail(ctx context.Context, email string) (int, error) {
	statement, err := ar.DbProvider.Db.Prepare("select id from user WHERE email = ?")
	if err != nil {
		log.Print(err)
		return 0, err
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
