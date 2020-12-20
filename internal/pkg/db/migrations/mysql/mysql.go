package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/sirupsen/logrus"
	"log"
)

type IDbConfig interface {
	GetDbHost() string
	GetDbPort() string
	GetDbUsername() string
	GetDbPassword() string
	GetDbDatabase() string
}

type IDbProvider interface {
	Migrate()
}

type DbProvider struct {
	Db *sql.DB
}

func InitDB(cfg IDbConfig) *DbProvider {
	var dbCon DbProvider
	connectionStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.GetDbUsername(), cfg.GetDbPassword(), cfg.GetDbHost(), cfg.GetDbPort(), cfg.GetDbDatabase())
	db, err := sql.Open("mysql", connectionStr)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	fmt.Println("DB connection success!!!")
	dbCon.Db = db
	return &dbCon
}

func (cp DbProvider) Migrate() {
	if err := cp.Db.Ping(); err != nil {
		logrus.WithError(err).Error(err)
		return
	}
	driver, _ := mysql.WithInstance(cp.Db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/db/migrations/mysql",
		"mysql",
		driver,
	)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.WithError(err).Error(err)
		return
	}
}
