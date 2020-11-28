package entities

type UserConfig struct {
	ID         int
	IsDarkMode bool
	IsListMode bool
	User       *User
}
