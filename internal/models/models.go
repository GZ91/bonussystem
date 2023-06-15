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
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accural    float64 `json:"accural,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

type DataBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
