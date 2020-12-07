package entities

type UserConfig struct {
	ID         int
	IsDarkMode bool
	IsListMode bool
	User       *User
}

type CreateUserConfig struct {
	IsDarkMode bool
	IsListMode bool
}

type UpdateUserConfig struct {
	ID         *int
	IsDarkMode *bool
	IsListMode *bool
}
