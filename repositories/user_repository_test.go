package repositories_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elgris/sqrl"
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

func (s *UserReposirotyTestSuite) Test_UpdateUserByFields() {
	id := 1
	name := "john"
	email := "haswin@gmail.com"
	u := entities.UpdateUser{
		ID:    &id,
		Name:  &name,
		Email: &email,
	}

	updateMap := map[string]interface{}{}
	updateMap["`name`"] = u.Name
	updateMap["`email`"] = u.Email

	qb := sqrl.Update("`user`").SetMap(updateMap).Where(sqrl.Eq{"`id`": u.ID})
	query, _, err := qb.ToSql()
	s.NoError(err)
	s.Mock.ExpectExec(query).WithArgs(u.Email, u.Name, u.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	// now we execute our method
	lastId, err := s.Resolver.IUserRepository.UpdateUserByFields(s.Ctx, u)
	s.Nil(err)
	s.Equal(1, lastId)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)
}

func (s *UserReposirotyTestSuite) Test_FindUserByID() {
	user := entities.User{
		ID:    1,
		Name:  "john",
		Email: "haswin@gmail.com",
	}

	query := "select u.id, u.name, u.email from user u WHERE id = ?"
	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, user.Name, user.Email)

	s.Mock.ExpectPrepare(query)
	s.Mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	// now we execute our method
	user, err := s.Resolver.IUserRepository.FindUserByID(s.Ctx, 1)
	s.Nil(err)
	s.Equal(user.ID, user.ID)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)
}

func (s *UserReposirotyTestSuite) Test_FindUserByEmail() {
	user := entities.User{
		ID:    1,
		Name:  "john",
		Email: "haswin@gmail.com",
	}

	query := "select u.id, u.name, u.email from user u WHERE email = ?"
	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, user.Name, user.Email)

	s.Mock.ExpectPrepare(query)
	s.Mock.ExpectQuery(query).WithArgs(user.Email).WillReturnRows(rows)

	// now we execute our method
	user, err := s.Resolver.IUserRepository.FindUserByEmail(s.Ctx, user.Email)
	s.Nil(err)
	s.Equal(user.ID, user.ID)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)
}

//
