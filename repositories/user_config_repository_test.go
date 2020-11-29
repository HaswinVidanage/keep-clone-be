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

type UserConfigRepositoryTestSuite struct {
	*test.TestApp
	suite.Suite
}

func NewUserConfigRepositoryTestSuite() *UserConfigRepositoryTestSuite {
	app, err := wire.GetTestApp()
	if err != nil {
		panic(err)
	}
	return &UserConfigRepositoryTestSuite{TestApp: app}
}

func TestUserConfigRepository(t *testing.T) {
	var testSuite *UserConfigRepositoryTestSuite = NewUserConfigRepositoryTestSuite()
	suite.Run(t, testSuite)
}

func (s *UserConfigRepositoryTestSuite) Test_InsertUserConfig() {
	uc := entities.CreateUserConfig{
		IsListMode: true,
		IsDarkMode: true,
	}
	fkUser := 1
	query := "INSERT INTO user_config (isDarkMode, isListMode, fk_user) VALUES (?,?,?)"
	s.Mock.ExpectPrepare(query)
	s.Mock.ExpectExec(query).WithArgs(uc.IsDarkMode, uc.IsListMode, fkUser).WillReturnResult(sqlmock.NewResult(1, 1))

	// now we execute our method
	lastId, err := s.Resolver.IUserConfigRepository.InsertUserConfig(s.Ctx, fkUser, uc)
	s.Nil(err)
	s.Equal(1, lastId)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)

}

func (s *UserConfigRepositoryTestSuite) Test_UpdateUserConfigByFields() {
	id := 1
	isListMode := true
	isDarkMode := true

	uc := entities.UpdateUserConfig{
		ID:         &id,
		IsListMode: &isListMode,
		IsDarkMode: &isDarkMode,
	}

	updateMap := map[string]interface{}{}
	updateMap["`isDarkMode`"] = uc.IsDarkMode
	updateMap["`isListMode`"] = uc.IsListMode

	qb := sqrl.Update("`user_config`").SetMap(updateMap).Where(sqrl.Eq{"`id`": uc.ID})
	query, _, err := qb.ToSql()
	s.NoError(err)
	s.Mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))

	// now we execute our method
	lastId, err := s.Resolver.IUserConfigRepository.UpdateUserConfigByFields(s.Ctx, uc)
	s.Nil(err)
	s.Equal(1, lastId)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)
}

func (s *UserConfigRepositoryTestSuite) Test_FindUserConfigByUserID() {
	uc := entities.UserConfig{
		ID:         1,
		IsListMode: true,
		IsDarkMode: true,
	}

	fkUser := 1
	query := "select n.id, n.isDarkMode, n.isListMode from user_config n where n.fk_user = ?"
	rows := sqlmock.NewRows([]string{"id", "isDarkMode", "isListMode"}).
		AddRow(uc.ID, uc.IsDarkMode, uc.IsListMode)

	s.Mock.ExpectPrepare(query)
	s.Mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	userQuery := "select u.id, u.name, u.email from user u WHERE id = ?"
	userRows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "john", "john@mail.com")

	s.Mock.ExpectPrepare(userQuery)
	s.Mock.ExpectQuery(userQuery).WithArgs(1).WillReturnRows(userRows)

	// now we execute our method
	userConfig, err := s.Resolver.IUserConfigRepository.FindUserConfigByUserID(s.Ctx, fkUser)
	s.Nil(err)
	s.Equal(1, userConfig.ID)
	s.Equal(uc.IsListMode, userConfig.IsListMode)
	s.Equal(uc.IsDarkMode, userConfig.IsDarkMode)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)
}
