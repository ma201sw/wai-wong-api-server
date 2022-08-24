package sumapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-wai-wong/internal/config"
	"go-wai-wong/internal/constant"
	"go-wai-wong/internal/golib"
	"go-wai-wong/internal/provider/jsonprovider"
	"go-wai-wong/internal/tokenhelper"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func Test_ValidateToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	config.LoadConfig()

	router := chi.NewRouter()
	server := httptest.NewServer(router)

	defer server.Close()

	client := &http.Client{}

	tokenHelperSrv := tokenhelper.New()
	goLibSrv := golib.New()

	router.Use(golib.Inject(goLibSrv))
	router.Use(tokenhelper.Inject(tokenHelperSrv))
	router.Use(validateToken)
	router.Route("/sumapi/v1", func(r chi.Router) {
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {})
	})

	ctx = golib.WithGoLib(ctx, goLibSrv)

	request, err := http.NewRequestWithContext(ctx, "GET", server.URL+"/sumapi/v1/test", http.NoBody)
	if err != nil {
		t.Fatalf("Could not make the request: %v", err)
	}

	tokenStr, err := tokenHelperSrv.GenToken(ctx, "testid") // generate the token
	if err != nil {
		t.Fatalf("Get token failed, %v", err)
	}

	request.Header.Add("Authorization", "Bearer "+tokenStr)

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("Request failed!, %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Request not OK: %v", response.StatusCode)
	}
}

func Test_MissingBearerToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	config.LoadConfig()

	router := chi.NewRouter()
	server := httptest.NewServer(router)

	defer server.Close()

	client := &http.Client{}

	tokenHelperSrv := tokenhelper.New()

	router.Use(tokenhelper.Inject(tokenHelperSrv))
	router.Use(validateToken)

	InstallRoutes(router)

	request, err := http.NewRequestWithContext(ctx, "GET", server.URL+"/sumapi/v1/test", http.NoBody)
	if err != nil {
		t.Fatalf("Could not make the request: %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("Could not make the request: %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Request was wrongfully authorized")
	}

	if response.Header.Get("WWW-Authenticate") != "Bearer" {
		t.Fatalf("No authentication error information provided")
	}
}

func Test_BadSecret(t *testing.T) {
	t.Parallel()

	config.LoadConfig()

	ctx := context.Background()

	router := chi.NewRouter()
	server := httptest.NewServer(router)

	defer server.Close()

	client := &http.Client{}

	tokenHelperSrv := tokenhelper.New()

	router.Use(tokenhelper.Inject(tokenHelperSrv))
	router.Use(validateToken)

	InstallRoutes(router)

	request, err := http.NewRequestWithContext(ctx, "GET", server.URL+"/sumapi/v1/test", http.NoBody)
	if err != nil {
		t.Fatalf("Could not make the request: %v", err)
	}

	claims := jwt.StandardClaims{
		Subject:   "usertest",
		ExpiresAt: time.Now().Add(viper.GetDuration(constant.TokenExpiresIn)).Unix(),
		Audience:  viper.GetString(constant.TokenAudience),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte("wrongsecret"))
	if err != nil {
		t.Fatalf("Could not create signed string %v", err)
	}

	request.Header.Add("Authorization", "Bearer "+tokenStr)

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("Could not make the request: %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Request was wrongfully authorized")
	}
}

func Test_handleAuth(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	config.LoadConfig()

	client := &http.Client{}

	type args struct {
		method string
		url    string
		body   string
	}

	tests := []struct {
		name               string
		args               args
		tokenClientMock    func(t *testing.T) *tokenhelper.TokenClientImplMock
		goLibMock          func(t *testing.T) *golib.GoLibImplMock
		expectedStatusCode int
	}{
		{
			name: "handAuthTest-noBody",
			args: args{
				method: "POST",
				url:    "/sumapi/v1/auth",
				body:   "",
			},
			expectedStatusCode: 500,
			tokenClientMock: func(t *testing.T) *tokenhelper.TokenClientImplMock {
				t.Helper()

				return &tokenhelper.TokenClientImplMock{}
			},
			goLibMock: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
		},
		{
			name: "handAuthTest-getNotAllowed",
			args: args{
				method: "GET",
				url:    "/sumapi/v1/auth",
				body:   "",
			},
			expectedStatusCode: 405,
			tokenClientMock: func(t *testing.T) *tokenhelper.TokenClientImplMock {
				t.Helper()

				return &tokenhelper.TokenClientImplMock{}
			},
			goLibMock: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
		},
		{
			name: "handAuthTest-validRequest",
			args: args{
				method: "POST",
				url:    "/sumapi/v1/auth",
				body:   `{"username": "test", "password": "test"}`,
			},
			expectedStatusCode: 200,
			tokenClientMock: func(t *testing.T) *tokenhelper.TokenClientImplMock {
				t.Helper()

				return &tokenhelper.TokenClientImplMock{}
			},
			goLibMock: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
		},
		{
			name: "handAuthTest-genTokenErr",
			args: args{
				method: "POST",
				url:    "/sumapi/v1/auth",
				body:   `{"username": "test", "password": "test"}`,
			},
			expectedStatusCode: 500,
			tokenClientMock: func(t *testing.T) *tokenhelper.TokenClientImplMock {
				t.Helper()

				return &tokenhelper.TokenClientImplMock{
					GenTokenFn: func(ctx context.Context, username string) (string, error) {
						return "", fmt.Errorf("token error")
					},
				}
			},
			goLibMock: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
		},
		{
			name: "handAuthTest-usernamePasswordEmpty",
			args: args{
				method: "POST",
				url:    "/sumapi/v1/auth",
				body:   `{"username": "", "password": ""}`,
			},
			expectedStatusCode: 403,
			tokenClientMock: func(t *testing.T) *tokenhelper.TokenClientImplMock {
				t.Helper()

				return &tokenhelper.TokenClientImplMock{}
			},
			goLibMock: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
		},
		{
			name: "handAuthTest-failedIOCopy",
			args: args{
				method: "POST",
				url:    "/sumapi/v1/auth",
				body:   `{"username": "", "password": ""}`,
			},
			expectedStatusCode: 500,
			tokenClientMock: func(t *testing.T) *tokenhelper.TokenClientImplMock {
				t.Helper()

				return &tokenhelper.TokenClientImplMock{}
			},
			goLibMock: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{
					CopyFn: func(dst io.Writer, src io.Reader) (written int64, err error) {
						return 0, fmt.Errorf("test error")
					},
				}
			},
		},
		{
			name: "noSuchRoute",
			args: args{
				method: "POST",
				url:    "/sumapi/v1/nosuchroute",
				body:   `{"username": "test", "password": "test"}`,
			},
			expectedStatusCode: 401,
			tokenClientMock: func(t *testing.T) *tokenhelper.TokenClientImplMock {
				t.Helper()

				return &tokenhelper.TokenClientImplMock{}
			},
			goLibMock: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router := chi.NewRouter()
			server := httptest.NewServer(router)

			t.Cleanup(func() { server.Close() })

			router = chi.NewRouter()

			server = httptest.NewServer(router)

			router.Use(golib.Inject(tt.goLibMock(t)))
			router.Use(tokenhelper.Inject(tt.tokenClientMock(t)))

			InstallRoutes(router)

			reader := strings.NewReader(tt.args.body)

			request, err := http.NewRequestWithContext(ctx, tt.args.method, server.URL+tt.args.url, reader)
			if err != nil {
				t.Fatalf("Could not make the request: %v", err)
			}

			response, err := client.Do(request)
			if err != nil {
				t.Fatalf("Could not make the request: %v", err)
			}

			defer response.Body.Close()

			if response.StatusCode != tt.expectedStatusCode {
				t.Fatalf("Response status code: %v does not match expected status code: %v", response.StatusCode, tt.expectedStatusCode)
			}
		})
	}
}

func Test_handleSum(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	config.LoadConfig()

	client := &http.Client{}

	type args struct {
		method string
		body   string
	}

	tests := []struct {
		name                   string
		args                   args
		jsonProviderClientMock func(t *testing.T) *jsonprovider.JSONProviderClientImplMock
		expectedStatusCode     int
		expectedSum            int
		goLibMock              func(t *testing.T) *golib.GoLibImplMock
	}{
		{
			name: "handSumTest-successfulResponse",
			args: args{
				method: "POST",
				body: `{
					"data1": [1,2,3,4],
					"data2": {"a":6,"b":4},
					"data3": [[[2]]],
					"data4": {"a":{"b":4},"c":-2},
					"data5": {"a":[-1,1,"dark"]},
					"data6": [-1,{"a":1, "b":"light"}],
					"data7": [],
					"data8": {}
				}`,
			},
			expectedStatusCode: 200,
			jsonProviderClientMock: func(t *testing.T) *jsonprovider.JSONProviderClientImplMock {
				t.Helper()

				return &jsonprovider.JSONProviderClientImplMock{}
			},
			goLibMock: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
			expectedSum: 24,
		},
		{
			name: "handSumTest-badjson",
			args: args{
				method: "POST",
				body:   `asdf`,
			},
			expectedStatusCode: 400,
			jsonProviderClientMock: func(t *testing.T) *jsonprovider.JSONProviderClientImplMock {
				t.Helper()

				return &jsonprovider.JSONProviderClientImplMock{}
			},
			goLibMock: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{}
			},
			expectedSum: 0,
		},
		{
			name: "handSumTest-ioCopyError",
			args: args{
				method: "POST",
				body: `{
					"data1": [1,2,3,4],
					"data2": {"a":6,"b":4},
					"data3": [[[2]]],
					"data4": {"a":{"b":4},"c":-2},
					"data5": {"a":[-1,1,"dark"]},
					"data6": [-1,{"a":1, "b":"light"}],
					"data7": [],
					"data8": {}
				}`,
			},
			expectedStatusCode: 500,
			jsonProviderClientMock: func(t *testing.T) *jsonprovider.JSONProviderClientImplMock {
				t.Helper()

				return &jsonprovider.JSONProviderClientImplMock{}
			},
			goLibMock: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{
					CopyFn: func(dst io.Writer, src io.Reader) (written int64, err error) {
						return 0, fmt.Errorf("test error")
					},
				}
			},
			expectedSum: 0,
		},
		{
			name: "handSumTest-unMarshalErr",
			args: args{
				method: "POST",
				body: `{
					"data1": [1,2,3,4],
					"data2": {"a":6,"b":4},
					"data3": [[[2]]],
					"data4": {"a":{"b":4},"c":-2},
					"data5": {"a":[-1,1,"dark"]},
					"data6": [-1,{"a":1, "b":"light"}],
					"data7": [],
					"data8": {}
				}`,
			},
			expectedStatusCode: 400,
			jsonProviderClientMock: func(t *testing.T) *jsonprovider.JSONProviderClientImplMock {
				t.Helper()

				return &jsonprovider.JSONProviderClientImplMock{}
			},
			goLibMock: func(t *testing.T) *golib.GoLibImplMock {
				t.Helper()

				return &golib.GoLibImplMock{
					UnmarshalFn: func(data []byte, v interface{}) error {
						return fmt.Errorf("test error")
					},
				}
			},
			expectedSum: 0,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router := chi.NewRouter()
			server := httptest.NewServer(router)

			t.Cleanup(func() { server.Close() })

			router = chi.NewRouter()

			server = httptest.NewServer(router)
			router.Use(golib.Inject(tt.goLibMock(t)))
			router.Use(jsonprovider.Inject(tt.jsonProviderClientMock(t)))

			router.Route("/sumapi/v1", func(router chi.Router) {
				router.Post("/sum", handleSum)
			})

			reader := strings.NewReader(tt.args.body)

			request, err := http.NewRequestWithContext(ctx, tt.args.method, server.URL+"/sumapi/v1/sum", reader)
			if err != nil {
				t.Fatalf("Could not make the request: %v", err)
			}

			response, err := client.Do(request)
			if err != nil {
				t.Fatalf("Could not make the request: %v", err)
			}

			defer response.Body.Close()

			if response.StatusCode != tt.expectedStatusCode {
				t.Fatalf("Response status code: %v does not match expected status code: %v", response.StatusCode, tt.expectedStatusCode)
			}

			responseBytes, err := io.ReadAll(response.Body)
			if err != nil {
				panic(err)
			}

			var sumResponse SumResponse

			json.Unmarshal(responseBytes, &sumResponse)

			if sumResponse.Sum != tt.expectedSum {
				t.Fatalf("response sum: %v does not match expected sum: %v", sumResponse, tt.expectedSum)
			}
		})
	}
}
