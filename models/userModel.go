package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// User model definition
type User struct {
	Id             int        `json:"id"`
	UserName       string     `json:"user_name"`
	Email          string     `json:"email"`
	PasswordHash   string     `json:"password_hash"`
	Address        NullString `json:"address"`
	PhoneNumber    NullString `json:"phone_number"`
	ProfilePicture NullString `json:"profile_picture"`
	InstagramLink  NullString `json:"instagram_link"`
	Description    NullString `json:"description"`
	Role           string     `json:"role"` // "buyer" or "seller"
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      NullTime   `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Id int `json:"id"`
	jwt.StandardClaims
}
