package model

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type MyCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}
