package models

import (
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

type Credentials struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

// Create a struct to hold each client's rate limiter
type Client struct {
	Limiter *rate.Limiter
}
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Book struct {
	gorm.Model
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}
