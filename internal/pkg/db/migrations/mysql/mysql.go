package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/google/wire"
	"log"
)

var Db *sql.DB

type IDbConfig interface {
	GetDbHost() string
	GetDbPort() string
	GetDbUsername() string
	GetDbPassword() string
	GetDbDatabase() string
}

type IConnectionProvider interface {
	Migrate()
}

type ConnectionProvider struct {
	Db *sql.DB
}

var NewConnectionProvider = wire.NewSet(
	wire.Struct(new(ConnectionProvider), "*"),
	wire.Bind(new(IConnectionProvider), new(*ConnectionProvider)))

func InitDB(cfg IDbConfig) *ConnectionProvider {
	var dbCon ConnectionProvider
	fmt.Println("DB HOST WIRRED : (IF EMPTY DON'T GIVE UP) :", cfg.GetDbPort())
	db, err := sql.Open("mysql", "sa:qweqwe@tcp(localhost:3305)/hackernews_db")
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	dbCon.Db = db
	Db = db
	return &dbCon
}

func (cp ConnectionProvider) Migrate() {
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
