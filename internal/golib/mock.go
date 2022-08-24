package golib

import (
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/golang-jwt/jwt"
)

type GoLibImplMock struct {
	StdEncodingDecodeStringFn func(s string) ([]byte, error)
	ParseWithClaimsFn         func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)
	StandardClaimsFn          func(subject string, expiresAt int64, audience string) jwt.StandardClaims
	CopyFn                    func(dst io.Writer, src io.Reader) (written int64, err error)
	UnmarshalFn               func(data []byte, v interface{}) error
	MarshalFn                 func(v interface{}) ([]byte, error)
}

func (c *GoLibImplMock) StdEncodingDecodeString(s string) ([]byte, error) {
	if c != nil && c.StdEncodingDecodeStringFn != nil {
		return c.StdEncodingDecodeStringFn(s)
	}

	return base64.StdEncoding.DecodeString(s)
}

func (c *GoLibImplMock) ParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	if c != nil && c.ParseWithClaimsFn != nil {
		return c.ParseWithClaimsFn(tokenString, claims, keyFunc)
	}

	return jwt.ParseWithClaims(tokenString, claims, keyFunc)
}

func (c *GoLibImplMock) StandardClaims(subject string, expiresAt int64, audience string) jwt.StandardClaims {
	if c != nil && c.StandardClaimsFn != nil {
		return c.StandardClaimsFn(subject, expiresAt, audience)
	}

	goLibSrv := New()

	return goLibSrv.StandardClaims(subject, expiresAt, audience)
}

func (c *GoLibImplMock) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	if c != nil && c.CopyFn != nil {
		return c.CopyFn(dst, src)
	}

	return io.Copy(dst, src)
}

func (c *GoLibImplMock) Unmarshal(data []byte, v interface{}) error {
	if c != nil && c.UnmarshalFn != nil {
		return c.UnmarshalFn(data, v)
	}

	return json.Unmarshal(data, v)
}

func (c *GoLibImplMock) Marshal(v interface{}) ([]byte, error) {
	if c != nil && c.MarshalFn != nil {
		return c.MarshalFn(v)
	}

	return json.Marshal(v)
}
