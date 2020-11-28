package entities

type Note struct {
	ID      int
	Title   string
	Content string
	User    *User
}
