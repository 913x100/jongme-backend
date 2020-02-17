package model

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type AuthenticateClaims struct {
	UserID string `json:"user_id"`
	PageID string `json:"page_id"`
	jwt.StandardClaims
}

type ExpirePageClaims struct {
	UserID string `json:"user_id"`
	PageID string `json:"page_id"`
	jwt.StandardClaims
}
