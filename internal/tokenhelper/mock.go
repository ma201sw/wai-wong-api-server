package tokenhelper

import (
	"context"
)

type TokenClientImplMock struct {
	GenTokenFn    func(ctx context.Context, username string) (string, error)
	VerifyTokenFn func(ctx context.Context, tokenStr string) (string, error)
}

func (c *TokenClientImplMock) GenToken(ctx context.Context, username string) (string, error) {
	if c != nil && c.GenTokenFn != nil {
		return c.GenTokenFn(ctx, username)
	}

	tokenHelperSrv := New()

	return tokenHelperSrv.GenToken(ctx, username)
}

func (c *TokenClientImplMock) VerifyToken(ctx context.Context, tokenStr string) (string, error) {
	if c != nil && c.VerifyTokenFn != nil {
		return c.VerifyTokenFn(ctx, tokenStr)
	}

	tokenHelperSrv := New()

	return tokenHelperSrv.VerifyToken(ctx, tokenStr)
}
