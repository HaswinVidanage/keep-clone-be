package note

import (
	"github.com/google/wire"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"hackernews-api/services/users"
	"log"
)

type Note struct {
	ID      string
	Title   string
	Content string
	User    *users.User
}

type INoteService interface {
	Save(note Note) int64
	GetAll() []Note
}

type NoteService struct {
	DbProvider *database.DbProvider
}

var NewNoteService = wire.NewSet(
	wire.Struct(new(NoteService), "*"),
	wire.Bind(new(INoteService), new(*NoteService)))

func (ns NoteService) Save(note Note) int64 {
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

func (ns NoteService) GetAll() []Note {
	stmt, err := ns.DbProvider.Db.Prepare("select n.id, n.title, n.content, n.fk_user, u.name from note n inner join user u on n.fk_user = u.ID") // changed
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var notes []Note
	var name string
	var id string
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &id, &name)
		if err != nil {
			log.Fatal(err)
		}
		note.User = &users.User{
			ID:   id,
			Name: name,
		}
		notes = append(notes, note)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return notes
}
