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
	SaveNote(ctx context.Context, note entities.Note) (int, error)
	GetAll(ctx context.Context) ([]entities.Note, error)
	DeleteNote(ctx context.Context, id int) (bool, error)
}

type NoteService struct {
	DbProvider     *database.DbProvider
	NoteRepository repositories.INoteRepository
}

var NewNoteService = wire.NewSet(
	wire.Struct(new(NoteService), "*"),
	wire.Bind(new(INoteService), new(*NoteService)))

func (ns NoteService) SaveNote(ctx context.Context, note entities.Note) (int, error) {
	id, err := ns.NoteRepository.InsertNote(ctx, entities.CreateNote{
		Title:   note.Title,
		Content: note.Content,
		User:    note.User,
	})

	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	return id, nil
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

func (ns NoteService) DeleteNote(ctx context.Context, id int) (bool, error) {
	userCtx := auth.ForContext(ctx)
	if userCtx == nil {
		log.Fatal("unauthorised")
	}
	isDeleted, err := ns.NoteRepository.DeleteNoteByID(ctx, id, userCtx.ID)
	if err != nil {
		return false, err
	}

	return isDeleted, nil
}
