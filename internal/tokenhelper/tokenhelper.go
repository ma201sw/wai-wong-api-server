package tokenhelper

import (
	"context"
	"fmt"
	"time"

	"go-wai-wong/common"
	"go-wai-wong/internal/constant"
	"go-wai-wong/internal/golib"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type Service interface {
	GenToken(ctx context.Context, username string) (string, error)
	VerifyToken(ctx context.Context, tokenStr string) (string, error)
}

type tokenHelperImpl struct{}

// verify interface compliance
var _ Service = (*tokenHelperImpl)(nil)

func New() tokenHelperImpl {
	return tokenHelperImpl{}
}

func (c tokenHelperImpl) GenToken(ctx context.Context, username string) (string, error) {
	var goLibSrv golib.Service

	if err := golib.FromContextAs(ctx, &goLibSrv); err != nil {
		return "", fmt.Errorf("golib from context as err: %w", err)
	}

	claims := goLibSrv.StandardClaims(
		username,
		time.Now().Add(viper.GetDuration(constant.TokenExpiresIn)).Unix(),
		viper.GetString(constant.TokenAudience),
	)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret, err := goLibSrv.StdEncodingDecodeString(viper.GetString(constant.TokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to decode secret: %w", err)
	}

	signedString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("token signed string error: %w", err)
	}

	return signedString, nil
}

func (c tokenHelperImpl) VerifyToken(ctx context.Context, tokenStr string) (string, error) {
	var golibSrv golib.Service

	if err := golib.FromContextAs(ctx, &golibSrv); err != nil {
		return "", fmt.Errorf("get token err: %w", err)
	}

	claims := jwt.StandardClaims{}
	if _, err := golibSrv.ParseWithClaims(tokenStr, &claims, func(tok *jwt.Token) (interface{}, error) {
		secret, err := golibSrv.StdEncodingDecodeString(viper.GetString(constant.TokenSecret))
		if err != nil {
			return "", fmt.Errorf("failed to decode secret: %w", err)
		}

		return secret, nil
	}); err != nil {
		return "", fmt.Errorf("jwt parse with claims error: %w", err)
	}

	if !claims.VerifyAudience(viper.GetString(constant.TokenAudience), true) {
		return "", common.AudienceError(claims.Audience)
	}

	return claims.Subject, nil
}
