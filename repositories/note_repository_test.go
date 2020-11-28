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

type NoteRepositoryTestSuite struct {
	*test.TestApp
	suite.Suite
}

func NewNoteRepositoryTestSuite() *NoteRepositoryTestSuite {
	app, err := wire.GetTestApp()
	if err != nil {
		panic(err)
	}
	return &NoteRepositoryTestSuite{TestApp: app}
}

func TestNoteRepository(t *testing.T) {
	var testSuite *NoteRepositoryTestSuite = NewNoteRepositoryTestSuite()
	suite.Run(t, testSuite)
}

func (s *NoteRepositoryTestSuite) Test_InsertNote() {
	note := entities.CreateNote{
		Title:   "Test title",
		Content: "Test content",
		User: &entities.User{
			ID: 1,
		},
	}

	query := "INSERT INTO note(title, content, fk_user) VALUES(?,?, ?)"
	s.Mock.ExpectPrepare(query)
	s.Mock.ExpectExec(query).WithArgs(note.Title, note.Content, note.User.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	// now we execute our method
	lastId, err := s.Resolver.INoteRepository.InsertNote(s.Ctx, note)
	s.Nil(err)
	s.Equal(1, lastId)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)

}

func (s *NoteRepositoryTestSuite) Test_FindNoteByID() {
	n := entities.Note{
		Title:   "Test title",
		Content: "Test content",
		User: &entities.User{
			ID:    1,
			Name:  "John",
			Email: "haswin@gmail.com",
		},
	}

	noteQuery := "select n.id, n.title, n.content, n.fk_user from note n WHERE n.id = ?"
	noteRows := sqlmock.NewRows([]string{"id", "title", "content", "fk_user"}).
		AddRow(1, n.Title, n.Content, n.User.ID)

	s.Mock.ExpectPrepare(noteQuery)
	s.Mock.ExpectQuery(noteQuery).WithArgs(1).WillReturnRows(noteRows)

	userQuery := "select u.id, u.name, u.email from user u WHERE id = ?"
	userRows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(n.User.ID, n.User.Name, n.User.Email)

	s.Mock.ExpectPrepare(userQuery)
	s.Mock.ExpectQuery(userQuery).WithArgs(1).WillReturnRows(userRows)

	// now we execute our method
	notes, err := s.Resolver.INoteRepository.FindNoteByID(s.Ctx, n.User.ID)
	s.Nil(err)
	s.Equal(1, notes.User.ID)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)
}

func (s *NoteRepositoryTestSuite) Test_FindNotesByUserID() {
	n := entities.Note{
		Title:   "Test title",
		Content: "Test content",
		User: &entities.User{
			ID:    1,
			Name:  "John",
			Email: "haswin@gmail.com",
		},
	}

	noteQuery := "select n.id, n.title, n.content from note n where n.fk_user = ?"
	noteRows := sqlmock.NewRows([]string{"id", "title", "content"}).
		AddRow(1, n.Title, n.Content)

	userQuery := "select u.id, u.name, u.email from user u WHERE id = ?"
	userRows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(n.User.ID, n.User.Name, n.User.Email)

	s.Mock.ExpectPrepare(userQuery)
	s.Mock.ExpectQuery(userQuery).WithArgs(1).WillReturnRows(userRows)

	s.Mock.ExpectPrepare(noteQuery)
	s.Mock.ExpectQuery(noteQuery).WithArgs(1).WillReturnRows(noteRows)

	// now we execute our method
	notes, err := s.Resolver.INoteRepository.FindNotesByUserID(s.Ctx, n.User.ID)
	s.Nil(err)
	s.Equal(1, notes[0].ID)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)
}

func (s *NoteRepositoryTestSuite) Test_UpdateNoteByFields() {
	id := 1
	title := "Test title"
	content := "Test content"

	n := entities.UpdateNote{
		ID:      &id,
		Title:   &title,
		Content: &content,
	}

	updateMap := map[string]interface{}{}
	updateMap["`title`"] = title
	updateMap["`content`"] = content

	// 'UPDATE `note` SET `content` = ?, `title` = ? WHERE `id` = ?'

	qb := sqrl.Update("`note`").SetMap(updateMap).Where(sqrl.Eq{"`id`": id})
	query, _, err := qb.ToSql()
	s.NoError(err)
	s.Mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))

	// now we execute our method
	lastId, err := s.Resolver.INoteRepository.UpdateNoteByFields(s.Ctx, n)
	s.Nil(err)
	s.Equal(1, lastId)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)
}

func (s *NoteRepositoryTestSuite) Test_DeleteNoteByID() {
	note := entities.CreateNote{
		ID:      1,
		Title:   "Test title",
		Content: "Test content",
		User: &entities.User{
			ID: 2,
		},
	}

	query := "delete from note where id = ? and fk_user = ?"
	s.Mock.ExpectPrepare(query)
	s.Mock.ExpectExec(query).WithArgs(note.ID, note.User.ID).WillReturnResult(sqlmock.NewResult(3, 1))

	// now we execute our method
	lastId, err := s.Resolver.INoteRepository.DeleteNoteByID(s.Ctx, note.ID, note.User.ID)
	s.Nil(err)
	s.Equal(3, lastId)

	// we make sure that all expectations were met
	err = s.Mock.ExpectationsWereMet()
	s.Nil(err)

}
