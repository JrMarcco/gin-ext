package jwt

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type JClaims[T any] struct {
	Data T `json:"data"`
	jwt.RegisteredClaims
}

type Manager[T any] interface {
	GenerateToken(data T) (string, error)
	VerifyToken(token string) (*JClaims[T], error)
}
