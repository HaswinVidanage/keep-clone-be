package users

import (
	"context"
	"database/sql"
	_ "database/sql"
	"github.com/google/wire"
	"golang.org/x/crypto/bcrypt"
	"hackernews-api/internal/pkg/db/migrations/mysql"

	"log"
)

type IUserService interface {
	Create(ctx context.Context, user User)
	GetUserIdByUsername(ctx context.Context, username string) (int, error)
	Authenticate(ctx context.Context, user User) bool
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

type UserService struct {
	DbProvider *database.DbProvider
}

var NewUserService = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)))

func (us *UserService) Create(ctx context.Context, user User) {
	statement, err := us.DbProvider.Db.Prepare("INSERT INTO Users(Username,Password) VALUES(?,?)")
	print(statement)
	if err != nil {
		log.Fatal(err)
	}
	hashedPassword, err := HashPassword(user.Password)
	_, err = statement.Exec(user.Username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
}

//GetUserIdByUsername check if a user exists in database by given username
func (us *UserService) GetUserIdByUsername(ctx context.Context, username string) (int, error) {
	statement, err := us.DbProvider.Db.Prepare("select ID from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(username)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return Id, nil
}

func (us *UserService) Authenticate(ctx context.Context, user User) bool {
	statement, err := us.DbProvider.Db.Prepare("select Password from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(user.Username)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}

	return CheckPasswordHash(user.Password, hashedPassword)
}

//HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
