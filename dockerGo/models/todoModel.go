package models

import (
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type Credentials struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	gorm.Model
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Book struct {
	gorm.Model
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}
