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
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

type DataBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawData struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type WithdrawalsData struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type DataForProcessing struct {
	Order  string
	UserID string
	Status string
}

type ResponceAccural struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
