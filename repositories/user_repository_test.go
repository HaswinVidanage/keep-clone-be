package repositories_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"hackernews-api/entities"
	"hackernews-api/internal/wire"
	"hackernews-api/test"
	"testing"
)

type UserReposirotyTestSuite struct {
	*test.TestApp
	suite.Suite
}

func NewUserReposirotyTestSuite() *UserReposirotyTestSuite {
	app, err := wire.GetTestApp()
	if err != nil {
		panic(err)
	}
	return &UserReposirotyTestSuite{TestApp: app}
}

func TestVisitService(t *testing.T) {
	var testSuite *UserReposirotyTestSuite = NewUserReposirotyTestSuite()
	suite.Run(t, testSuite)
}

func (s *UserReposirotyTestSuite) TestSomething() {
	//assert.True(t, true, "True is true!")

	// todo get mock and db from the suite
	//s.DbProvider.Db
	//db, mock, err := sqlmock.New()
	//if err != nil {
	//	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	//}
	//defer db.Close()

	user := entities.CreateUser{
		Name:     "Haswin",
		Email:    "haswind@hotmail.com",
		Password: "123",
	}

	//s.Mock.ExpectationsWereMet()

	// todo try out this
	s.Mock.ExpectBegin()
	s.Mock.ExpectPrepare("insert into user( name, email, password) values(?,?,?)")
	s.Mock.ExpectExec("INSERT INTO user").WithArgs(user.Name, user.Email, user.Password).WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	// now we execute our method
	token, err := s.Resolver.IUserRepository.InsertUser(s.Ctx, user)
	s.Nil(err)
	s.NotEqual(token, "")

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)

}
