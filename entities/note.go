package entities

type Note struct {
	ID      int
	Title   string
	Content string
	User    *User
}

type CreateNote struct {
	ID      int
	Title   string
	Content string
	User    *User
}

type UpdateNote struct {
	ID      *int
	Title   *string
	Content *string
}
