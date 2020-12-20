package repositories

import (
	"context"
	"database/sql"
	"github.com/elgris/sqrl"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"hackernews-api/entities"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"log"
)

type INoteRepository interface {
	InsertNote(ctx context.Context, note entities.CreateNote) (int, error)
	UpdateNoteByFields(ctx context.Context, user entities.UpdateNote) (int, error)
	FindNoteByID(ctx context.Context, id int) (entities.Note, error)
	FindNotesByUserID(ctx context.Context, id int) ([]entities.Note, error)
	DeleteNoteByID(ctx context.Context, noteID int, userID int) (bool, error)
}

type NoteRepository struct {
	DbProvider     *database.DbProvider
	UserRepository IUserRepository
}

var NewNoteRepository = wire.NewSet(
	wire.Struct(new(NoteRepository), "*"),
	wire.Bind(new(INoteRepository), new(*NoteRepository)))

func (nr *NoteRepository) InsertNote(ctx context.Context, note entities.CreateNote) (int, error) {
	stmt, err := nr.DbProvider.Db.Prepare("INSERT INTO note(title, content, fk_user) VALUES(?,?, ?)")
	if err != nil {
		logrus.WithError(err).Warn(err)
		return 0, err
	}

	res, err := stmt.Exec(note.Title, note.Content, note.User.ID)
	if err != nil {
		logrus.WithError(err).Warn(err)
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		logrus.WithError(err).Warn(err)
		return 0, err
	}

	return int(lastId), nil
}

func (nr *NoteRepository) FindNoteByID(ctx context.Context, id int) (entities.Note, error) {
	statement, err := nr.DbProvider.Db.Prepare("select n.id, n.title, n.content, n.fk_user from note n WHERE n.id = ?")
	if err != nil {
		logrus.WithError(err)
		return entities.Note{}, err
	}
	row := statement.QueryRow(id)
	var note entities.Note
	var fkUser int
	err = row.Scan(&note.ID, &note.Title, &note.Content, &fkUser)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return entities.Note{}, err
	}

	// find user
	user, err := nr.UserRepository.FindUserByID(ctx, fkUser)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return entities.Note{}, err
	}
	note.User = &user
	return note, nil
}

func (nr *NoteRepository) FindNotesByUserID(ctx context.Context, id int) ([]entities.Note, error) {

	// find user
	user, err := nr.UserRepository.FindUserByID(ctx, id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return []entities.Note{}, err
	}

	stmt, err := nr.DbProvider.Db.Prepare("select n.id, n.title, n.content from note n where n.fk_user = ?")
	if err != nil {
		log.Print(err)
		return []entities.Note{}, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		log.Print(err)
		return []entities.Note{}, err
	}
	defer rows.Close()

	var notes []entities.Note
	for rows.Next() {
		var note entities.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content)
		if err != nil {
			log.Print(err)
			return []entities.Note{}, err
		}
		note.User = &user
		notes = append(notes, note)
	}
	if err = rows.Err(); err != nil {
		log.Print(err)
		return []entities.Note{}, err
	}

	return notes, nil
}

func (nr *NoteRepository) UpdateNoteByFields(ctx context.Context, n entities.UpdateNote) (int, error) {
	updateMap := map[string]interface{}{}
	if n.Title != nil {
		updateMap["`title`"] = *n.Title
	}
	if n.Content != nil {
		updateMap["`content`"] = *n.Content
	}

	qb := sqrl.Update("`note`").SetMap(updateMap).Where(sqrl.Eq{"`id`": n.ID})
	// run query
	result, err := qb.RunWith(nr.DbProvider.Db).Exec()
	if err != nil {
		logrus.WithError(err).Warn(err)
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		logrus.WithError(err).Warn(err)
		return 0, err
	}

	return int(lastId), nil
}

func (nr *NoteRepository) DeleteNoteByID(ctx context.Context, noteID int, userID int) (bool, error) {
	stmt, err := nr.DbProvider.Db.Prepare("delete from note where id = ? and fk_user = ?")
	if err != nil {
		logrus.WithError(err).Warn(err)
		return false, err
	}

	res, err := stmt.Exec(noteID, userID)
	if err != nil {
		logrus.WithError(err).Warn(err)
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logrus.WithError(err).Warn(err)
		return false, err
	}

	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
