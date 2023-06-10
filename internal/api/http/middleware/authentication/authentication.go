package authentication

import (
	"context"
	"errors"
	"fmt"
	"github.com/GZ91/bonussystem/internal/app/logger"
	"github.com/GZ91/bonussystem/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"net/http"
)

type Configer interface {
	GetSecretKey() string
}

type NodeAuthentication struct {
	conf Configer
}

func New(conf Configer) *NodeAuthentication {
	return &NodeAuthentication{conf: conf}
}

func (Node *NodeAuthentication) Authentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := ""
		cookie, err := r.Cookie("Authorization")
		if err != nil && err != http.ErrNoCookie {
			logger.Log.Error("authorization", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.URL.String() == "/api/user/login" || r.URL.String() == "/api/user/register" {
			h.ServeHTTP(w, r)
			return
		}
		if errors.Is(err, http.ErrNoCookie) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var ok bool
		userID, ok, err = validGetAuthentication(Node.conf.GetSecretKey(), cookie.Value)
		if err != nil {
			logger.Log.Error("authorization", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var userIDCTX models.CtxString = "userID"
		r = r.WithContext(context.WithValue(r.Context(), userIDCTX, userID))
		h.ServeHTTP(w, r)
	})
}

func validGetAuthentication(SecretKey, tokenString string) (string, bool, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				strErr := fmt.Sprintf("unexpected signing method: %v", t.Header["alg"])
				logger.Log.Error(strErr)
				return nil, fmt.Errorf(strErr)
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		return "", false, err
	}

	if !token.Valid {
		return "", false, nil
	}

	return claims.UserID, true, nil
}
