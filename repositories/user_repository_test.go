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

	user := entities.CreateUser{
		Name:     "Haswin",
		Email:    "haswind@hotmail.com",
		Password: "123",
	}

	query := "insert into user( name, email, password) values(?,?,?)"
	s.Mock.ExpectPrepare(query)
	s.Mock.ExpectExec(query).WithArgs(user.Name, user.Email, user.Password).WillReturnResult(sqlmock.NewResult(1, 1))

	// now we execute our method
	lastId, err := s.Resolver.IUserRepository.InsertUser(s.Ctx, user)
	s.Nil(err)
	s.Equal(1, lastId)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)

}
