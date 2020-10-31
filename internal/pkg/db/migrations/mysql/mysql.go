package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	mysql2 "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hackernews-api/graph/model"
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
	Db *gorm.DB
}

func InitDB(cfg IDbConfig) *DbProvider {
	var dbCon DbProvider
	fmt.Println("DB HOST WIRRED : (IF EMPTY DON'T GIVE UP) :", cfg.GetDbPort())
	dsn := "sa:qweqwe@tcp(localhost:3305)/hackernews_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	sqlDB, err := db.DB()
	if err = sqlDB.Ping(); err != nil {
		log.Panic(err)
	}
	err = db.AutoMigrate(&model.User{}, &model.Link{})
	if err != nil {
		log.Panic(err)
	}
	dbCon.Db = db
	return &dbCon
}

func (cp DbProvider) Migrate() {
	sqlDB, err := cp.Db.DB()
	if err = sqlDB.Ping(); err != nil {
		log.Panic(err)
	}
	driver, _ := mysql2.WithInstance(sqlDB, &mysql2.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/db/migrations/mysql",
		"mysql",
		driver,
	)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
