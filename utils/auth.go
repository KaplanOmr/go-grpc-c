package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const jwtKey = "MY_AWESOME_KEY"

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)
	jwtKey := []byte(jwtKey)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", errors.New("CANNOT_CREATED_TOKEN")
	}

	return tokenString, nil
}

func CheckJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("TOKEN_INVALID")
	}

	return token.Claims.(*Claims).Username, nil
}
