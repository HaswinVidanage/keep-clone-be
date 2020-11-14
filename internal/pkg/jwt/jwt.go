package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// secret key being used to sign tokens
var (
	SecretKey = []byte("secret")
)

// GenerateToken generates a jwt token and assign a username to it's claims and return it
func GenerateToken(name string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

var mySigningKey = []byte("test")

// ParseToken parses a jwt token and returns the username in it's claims
func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//return SecretKey, nil
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error getting token")
		}

		// todo get secret from config
		//jwtSecret, err := ts.SecretsConfig.GetJwtSecret()
		//return []byte(jwtSecret), err
		return mySigningKey, nil
	})

	if token != nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			name := claims["name"].(string)
			return name, nil
		}
		return "", errors.New("claim failed")
	}

	return "", err
}
