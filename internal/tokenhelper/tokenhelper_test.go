package tokenhelper

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-wai-wong/internal/config"
	"go-wai-wong/internal/constant"
	"go-wai-wong/internal/golib"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func Test_GenToken(t *testing.T) {
	t.Parallel()

	config.LoadConfig()

	type args struct {
		username string
	}

	tests := []struct {
		name         string
		c            tokenHelperImpl
		args         args
		goLibMock    func(t *testing.T) *golib.GoLibImplMock
		tokenSuccess bool
		wantErr      bool
	}{
		{
			name: "genToken-successfulReturn",
			c:    tokenHelperImpl{},
			args: args{username: "testUsername"},
			goLibMock: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
			tokenSuccess: true,
			wantErr:      false,
		},
		{
			name: "genToken-decodingSecretErr",
			c:    tokenHelperImpl{},
			args: args{username: "testUsername"},
			goLibMock: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{
					StdEncodingDecodeStringFn: func(s string) ([]byte, error) {
						return []byte{}, fmt.Errorf("test error")
					},
				}
			},
			tokenSuccess: false,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			ctx = golib.WithGoLib(ctx, tt.goLibMock(t))

			got, err := tt.c.GenToken(ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Fatalf("client.GenToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got == "" && tt.tokenSuccess {
				t.Fatalf("client.GenToken() = %v, want token successful gen", got)
			}
		})
	}
}

func Test_VerifyToken(t *testing.T) {
	t.Parallel()

	config.LoadConfig()

	type args struct {
		username string
	}

	tests := []struct {
		name                    string
		args                    args
		c                       tokenHelperImpl
		goLibMockForVerifyToken func(t *testing.T) *golib.GoLibImplMock
		goLibMockForGenToken    func(t *testing.T) *golib.GoLibImplMock
		tokenHelperMock         func(t *testing.T) *TokenClientImplMock
		want                    string
		wantErr                 bool
	}{
		{
			name: "verifyToken-successfulVerify",
			c:    tokenHelperImpl{},
			goLibMockForVerifyToken: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
			goLibMockForGenToken: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
			tokenHelperMock: func(t *testing.T) *TokenClientImplMock {
				t.Helper()

				return &TokenClientImplMock{}
			},
			want:    "testUsername",
			wantErr: false,
			args:    args{username: "testUsername"},
		},
		{
			name: "verifyToken-failParseClaims",
			c:    tokenHelperImpl{},
			goLibMockForVerifyToken: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{
					ParseWithClaimsFn: func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
						return nil, fmt.Errorf("test error")
					},
				}
			},
			goLibMockForGenToken: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
			tokenHelperMock: func(t *testing.T) *TokenClientImplMock {
				t.Helper()

				return &TokenClientImplMock{}
			},
			want:    "",
			wantErr: true,
			args:    args{username: "testUsername"},
		},
		{
			name: "verifyToken-failStdEncodingDecodeString",
			c:    tokenHelperImpl{},
			goLibMockForVerifyToken: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{
					StdEncodingDecodeStringFn: func(s string) ([]byte, error) {
						return []byte{}, fmt.Errorf("test error")
					},
				}
			},
			goLibMockForGenToken: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
			tokenHelperMock: func(t *testing.T) *TokenClientImplMock {
				t.Helper()

				return &TokenClientImplMock{}
			},
			want:    "",
			wantErr: true,
			args:    args{username: "testUsername"},
		},
		{
			name: "verifyToken-badAudienceErr",
			c:    tokenHelperImpl{},
			goLibMockForVerifyToken: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
			goLibMockForGenToken: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{
					StandardClaimsFn: func(subject string, expiresAt int64, audience string) jwt.StandardClaims {
						return jwt.StandardClaims{
							Subject:   "testUsername",
							ExpiresAt: time.Now().Add(viper.GetDuration(constant.TokenExpiresIn)).Unix(),
							Audience:  "bad audience",
						}
					},
				}
			},
			tokenHelperMock: func(t *testing.T) *TokenClientImplMock {
				t.Helper()

				return &TokenClientImplMock{}
			},
			want:    "",
			wantErr: true,
			args:    args{username: "testUsername"},
		},
		{
			name: "verifyToken-expiredTokenErr",
			c:    tokenHelperImpl{},
			goLibMockForVerifyToken: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
			goLibMockForGenToken: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{
					StandardClaimsFn: func(subject string, expiresAt int64, audience string) jwt.StandardClaims {
						return jwt.StandardClaims{
							Subject:   "testUsername",
							ExpiresAt: time.Now().Add(-viper.GetDuration(constant.TokenExpiresIn)).Unix(),
							Audience:  viper.GetString(constant.TokenAudience),
						}
					},
				}
			},
			tokenHelperMock: func(t *testing.T) *TokenClientImplMock {
				t.Helper()

				return &TokenClientImplMock{}
			},
			want:    "",
			wantErr: true,
			args:    args{username: "testUsername"},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			ctx = golib.WithGoLib(ctx, tt.goLibMockForGenToken(t))

			tok, err := tt.tokenHelperMock(t).GenToken(ctx, tt.args.username)
			if err != nil {
				panic(err)
			}

			ctx = golib.WithGoLib(ctx, tt.goLibMockForVerifyToken(t))

			got, err := tt.c.VerifyToken(ctx, tok)
			if (err != nil) != tt.wantErr {
				t.Fatalf("tokenHelperImpl.VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("tokenHelperImpl.VerifyToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkGenToken(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()

	goLibSrv := golib.New()
	tokenHelperSrv := New()

	ctx = golib.WithGoLib(ctx, goLibSrv)

	for i := 0; i < b.N; i++ {
		_, _ = tokenHelperSrv.GenToken(ctx, "test")
	}
}
