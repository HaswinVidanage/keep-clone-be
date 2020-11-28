package note

import (
	"context"
	"github.com/google/wire"
	"hackernews-api/entities"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"hackernews-api/repositories"
	"hackernews-api/services/auth"
	"log"
)

type INoteService interface {
	SaveNote(note entities.Note) int64
	GetAll(ctx context.Context) ([]entities.Note, error)
}

type NoteService struct {
	DbProvider     *database.DbProvider
	NoteRepository repositories.INoteRepository
}

var NewNoteService = wire.NewSet(
	wire.Struct(new(NoteService), "*"),
	wire.Bind(new(INoteService), new(*NoteService)))

func (ns NoteService) SaveNote(note entities.Note) int64 {
	stmt, err := ns.DbProvider.Db.Prepare("INSERT INTO note(title, content, fk_user) VALUES(?,?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	//  execution of our sql statement.
	res, err := stmt.Exec(note.Title, note.Content, note.User.ID)
	if err != nil {
		log.Fatal(err)
	}

	// retrieving Id of inserted Link.
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}

	log.Print("Row inserted!")
	return id
}

func (ns NoteService) GetAll(ctx context.Context) ([]entities.Note, error) {
	userCtx := auth.ForContext(ctx)
	if userCtx == nil {
		log.Fatal("unauthorised")
	}
	notes, err := ns.NoteRepository.FindNotesByUserID(ctx, userCtx.ID)
	if err != nil {
		return nil, err
	}
	return notes, nil
}
