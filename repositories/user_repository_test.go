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

func NewUserRepositoryTestSuite() *UserReposirotyTestSuite {
	app, err := wire.GetTestApp()
	if err != nil {
		panic(err)
	}
	return &UserReposirotyTestSuite{TestApp: app}
}

func TestUserRepository(t *testing.T) {
	var testSuite *UserReposirotyTestSuite = NewUserRepositoryTestSuite()
	suite.Run(t, testSuite)
}

func (s *UserReposirotyTestSuite) Test_InsertUser() {
	user := entities.CreateUser{
		Name:     "john",
		Email:    "haswin@gmail.com",
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

func (s *UserReposirotyTestSuite) Test_GetUserIdByEmail() {
	email := "haswin@gmail.com"
	query := "select id from user WHERE email = ?"
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)

	s.Mock.ExpectPrepare(query)
	s.Mock.ExpectQuery(query).WithArgs(email).WillReturnRows(rows)

	// now we execute our method
	lastId, err := s.Resolver.IUserRepository.GetUserIdByEmail(s.Ctx, email)
	s.Nil(err)
	s.Equal(1, lastId)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)
}

func (s *UserReposirotyTestSuite) Test_UpdateUser() {
	user := entities.User{
		ID:    "1",
		Name:  "john",
		Email: "haswin@gmail.com",
	}

	qb := s.Resolver.IUserRepository.UpdateUserBaseQuery(s.Ctx, user)
	query, _, err := qb.ToSql()
	s.NoError(err)
	s.Mock.ExpectExec(query).WithArgs(user.Email, user.Name, user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	// now we execute our method
	lastId, err := s.Resolver.IUserRepository.UpdateUser(s.Ctx, user)
	s.Nil(err)
	s.Equal(1, lastId)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)
}

func (s *UserReposirotyTestSuite) Test_FindUserByID() {
	user := entities.User{
		ID:    "1",
		Name:  "john",
		Email: "haswin@gmail.com",
	}

	query := "select * from user WHERE id = ?"
	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, user.Name, user.Email)

	s.Mock.ExpectPrepare(query)
	s.Mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	// now we execute our method
	user, err := s.Resolver.IUserRepository.FindUserByID(s.Ctx, 1)
	s.Nil(err)
	s.Equal("1", user.ID)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)
}
