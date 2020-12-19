package note

import (
	"context"
	"errors"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"hackernews-api/entities"
	"hackernews-api/repositories"
	"hackernews-api/services/auth"
)

type INoteService interface {
	SaveNote(ctx context.Context, note entities.Note) (int, error)
	GetAll(ctx context.Context) ([]entities.Note, error)
	DeleteNote(ctx context.Context, id int) (bool, error)
}

type NoteService struct {
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
		logrus.WithError(err)
		return 0, err
	}
	return id, nil
}

func (ns NoteService) GetAll(ctx context.Context) ([]entities.Note, error) {
	userCtx := auth.ForContext(ctx)
	if userCtx == nil {
		err := errors.New("unauthorised")
		logrus.WithError(err).Warn(err)
		return nil, err
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
		err := errors.New("unauthorised")
		logrus.WithError(err).Warn(err)
		return false, err
	}
	isDeleted, err := ns.NoteRepository.DeleteNoteByID(ctx, id, userCtx.ID)
	if err != nil {
		logrus.WithError(err).Warn(err)
		return false, err
	}

	return isDeleted, nil
}
