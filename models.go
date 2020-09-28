package main

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
	Admin    bool   `json:"admin" bson:"admin"`
}

type article struct {
	Date  string `json:"date" bson:"date"`
	Title string `json:"title" bson:"title"`
	Text  string `json:"text" bson:"text"`
}

var mySigningKey = []byte("secret")

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

// HashPassword is used to hash password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash checks if the hashed password is valid
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
