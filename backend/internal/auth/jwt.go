package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var (
	key = []byte(os.Getenv("SECRET_KEY"))
)

type Claims struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, username string) (string, error) {
	Claims := Claims{userID, username, jwt.RegisteredClaims{
		Issuer:    "igropoisk",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	return token.SignedString(key)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
