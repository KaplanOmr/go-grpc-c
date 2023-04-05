package main

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

func generateJWT(username string) (string, error) {
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
