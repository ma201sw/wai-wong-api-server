package golib

import "github.com/golang-jwt/jwt"

type JwtClient interface {
	ParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)
	StandardClaims(subject string, expiresAt int64, audience string) jwt.StandardClaims
}

func (c goLibImpl) ParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, keyFunc)
}

func (c goLibImpl) StandardClaims(subject string, expiresAt int64, audience string) jwt.StandardClaims {
	return jwt.StandardClaims{
		Subject:   subject,
		ExpiresAt: expiresAt,
		Audience:  audience,
	}
}
