package models

import "github.com/golang-jwt/jwt/v4"

type DataRegisteration struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CtxString string

type Claims struct {
	*jwt.RegisteredClaims
	UserID string
}
