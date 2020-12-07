package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
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
	fmt.Println("DB Port:", cfg.GetDbPort())
	db, err := sql.Open("mysql", "sa:qweqwe@tcp(localhost:3306)/keep_db")
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
		log.Fatal(err)
	}
	driver, _ := mysql.WithInstance(cp.Db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/db/migrations/mysql",
		"mysql",
		driver,
	)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
