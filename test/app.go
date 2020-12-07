package test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/wire"
	"hackernews-api/graph"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"log"
)

type TestAppOptions struct {
	Resolver *graph.Resolver
}

type TestApp struct {
	*TestAppOptions
	Ctx  context.Context
	Mock sqlmock.Sqlmock
}

var mock sqlmock.Sqlmock

var NewTestApp = wire.NewSet(
	wire.Struct(new(TestAppOptions), "*"),
	InitTestApp,
)

func InitTestApp(ctx context.Context, options *TestAppOptions) *TestApp {
	return &TestApp{
		Ctx:            ctx,
		TestAppOptions: options,
		Mock:           mock,
	}
}

func InitMockDB() *database.DbProvider {
	var dbCon database.DbProvider
	db, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("an error was not expected when opening a stub database connection", err)
	}
	//defer db.Close()

	dbCon.Db = db
	mock = _mock
	return &dbCon
}
