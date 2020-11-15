package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/wire"
	"golang.org/x/crypto/bcrypt"
	"hackernews-api/entities"
	"hackernews-api/internal/pkg/db/migrations/mysql"
	"hackernews-api/internal/pkg/jwt"
	"hackernews-api/repositories"
	"log"
	"net/http"
	"strconv"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

type IAuthService interface {
	AuthMiddleware() func(http.Handler) http.Handler
	HashPassword(string) (string, error)
	CheckPasswordHash(string, string) bool
	Authenticate(context.Context, string, string) bool
	Login(context.Context, string, string) (string, error)
	RefreshToken(context.Context, string) (string, error)
}

type AuthService struct {
	DbProvider     *database.DbProvider
	UserRepository repositories.IUserRepository
}

var NewAuthService = wire.NewSet(
	wire.Struct(new(AuthService), "*"),
	wire.Bind(new(IAuthService), new(*AuthService)))

func (as AuthService) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			header := r.Header.Get("Authorization")

			if header == "" {
				next.ServeHTTP(w, r)
				return
			} else {

				//validate jwt token
				tokenStr := header
				email, err := jwt.ParseToken(tokenStr)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}

				// create user and check if user exists in db
				user := entities.User{Email: email}
				id, err := as.UserRepository.GetUserIdByEmail(r.Context(), email)
				if err != nil {
					// token parsing failed, nevertheless we allow routing
					next.ServeHTTP(w, r)
					return
				}

				user.ID = strconv.Itoa(id)
				// put it in context
				ctx := context.WithValue(r.Context(), userCtxKey, &user)

				// and call the next with our new context
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			}

		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *entities.User {
	raw, _ := ctx.Value(userCtxKey).(*entities.User)
	return raw
}

//HashPassword hashes given password
func (as AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPassword hash compares raw password with it's hashed values
func (as AuthService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (as AuthService) Authenticate(ctx context.Context, email string, password string) bool {
	statement, err := as.DbProvider.Db.Prepare("select password from user WHERE email = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(email)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}

	return as.CheckPasswordHash(password, hashedPassword)
}

func (as AuthService) Login(ctx context.Context, email string, password string) (string, error) {

	correct := as.Authenticate(ctx, email, password)
	if !correct {
		return "", errors.New("wrong name or password")
	}

	id, err := as.UserRepository.GetUserIdByEmail(ctx, email)
	token, err := jwt.GenerateToken(ctx, id, email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (as AuthService) RefreshToken(ctx context.Context, token string) (string, error) {
	email, err := jwt.ParseToken(token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	id, err := as.UserRepository.GetUserIdByEmail(ctx, email)
	token, err = jwt.GenerateToken(ctx, id, email)
	if err != nil {
		return "", err
	}
	return token, nil
}
