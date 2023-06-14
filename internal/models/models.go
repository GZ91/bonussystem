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

type DataOrder struct {
	Number     string
	Status     string
	Accural    int
	UploadedAt string
}

type DataOrderForJSON struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accural    string `json:"accural,omitempty"`
	UploadedAt string `json:"uploaded_at"`
}
