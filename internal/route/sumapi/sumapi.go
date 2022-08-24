package sumapi

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-wai-wong/common"
	"go-wai-wong/internal/constant"
	"go-wai-wong/internal/golib"
	"go-wai-wong/internal/provider/jsonprovider"
	"go-wai-wong/internal/tokenhelper"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresIn uint32 `json:"expires_in"`
}

type SumResponse struct {
	SHA256 string `json:"sha256"`
	Sum    int    `json:"sum"`
}

type AuthRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func writeResponse(respWriter http.ResponseWriter, data interface{}) {
	responseBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal response: %v", err)
		common.WriteInternalError(respWriter)

		return
	}

	if _, err := respWriter.Write(responseBytes); err != nil {
		log.Printf("failed to write response: %v", err)
		common.WriteInternalError(respWriter)

		return
	}
}

func handleAuth(respWriter http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	var goLibSrv golib.Service

	if err := golib.FromContextAs(
		ctx,
		&goLibSrv); err != nil {
		log.Printf("golib service type assert error")
		common.WriteInternalError(respWriter)

		return
	}

	var tokenHelperSrv tokenhelper.Service

	if err := tokenhelper.FromContextAs(
		ctx,
		&tokenHelperSrv); err != nil {
		log.Printf("token helper service type assert error")
		common.WriteInternalError(respWriter)

		return
	}

	requestBodyBuf := &bytes.Buffer{}

	_, err := goLibSrv.Copy(requestBodyBuf, request.Body)
	if err != nil {
		log.Printf("io copy error: %v", err)
		common.WriteInternalError(respWriter)

		return
	}

	var authRequestBody AuthRequestBody

	if unMarshalErr := goLibSrv.Unmarshal(requestBodyBuf.Bytes(), &authRequestBody); unMarshalErr != nil {
		log.Printf("failed to unmarshal: %v", unMarshalErr)
		common.WriteInternalError(respWriter)

		return
	}

	// TO DO: should validate username and password input with regex to stop attacks
	if authRequestBody.Username == "" || authRequestBody.Password == "" {
		log.Printf("username or password is empty")
		common.WriteError(respWriter, http.StatusForbidden, "FORBIDDEN", "")

		return
	}

	token, err := tokenHelperSrv.GenToken(ctx, authRequestBody.Username)
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		common.WriteInternalError(respWriter)

		return
	}

	response := &AuthResponse{
		Token:     token,
		ExpiresIn: uint32(viper.GetDuration(constant.TokenExpiresIn) / time.Second),
	}

	writeResponse(respWriter, response)
}

func handleSum(respWriter http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	var jsonProviderSrv jsonprovider.Service

	if err := jsonprovider.FromContextAs(
		ctx,
		&jsonProviderSrv); err != nil {
		log.Printf("json provider service type assert error")
		common.WriteInternalError(respWriter)

		return
	}

	var goLibSrv golib.Service

	if err := golib.FromContextAs(
		ctx,
		&goLibSrv); err != nil {
		log.Printf("golib service type assert error")
		common.WriteInternalError(respWriter)

		return
	}

	requestBodyBuf := &bytes.Buffer{}

	_, err := goLibSrv.Copy(requestBodyBuf, request.Body)
	if err != nil {
		log.Printf("io copy error: %v", err)
		common.WriteInternalError(respWriter)

		return
	}

	var jsonRequestBody map[string]interface{}

	if unmarshalErr := goLibSrv.Unmarshal(requestBodyBuf.Bytes(), &jsonRequestBody); unmarshalErr != nil {
		log.Printf("failed to unmarshal: %v", unmarshalErr)
		common.WriteError(respWriter, http.StatusBadRequest, "BAD REQUEST", "")

		return
	}

	floatSlice := []float64{}
	jsonProviderSrv.JSONMapToFloatSliceAs(jsonRequestBody, &floatSlice)

	sumResult := 0
	for _, v := range floatSlice {
		sumResult += int(v)
	}

	sha256Hash := sha256.New()

	if _, sha256WriteErr := sha256Hash.Write([]byte(strconv.Itoa(sumResult))); sha256WriteErr != nil {
		log.Printf("failed to write bytes: %v", sha256WriteErr)
		common.WriteInternalError(respWriter)

		return
	}

	hash := fmt.Sprintf("%x", sha256Hash.Sum(nil))

	response := &SumResponse{
		SHA256: hash,
		Sum:    sumResult,
	}

	writeResponse(respWriter, response)
}

func validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		var tokenHelperSrv tokenhelper.Service

		if err := tokenhelper.FromContextAs(
			ctx,
			&tokenHelperSrv); err != nil {
			log.Printf("token helper service type assert error")
			common.WriteInternalError(respWriter)

			return
		}

		// check that we have a bearer token
		auth := request.Header.Get("Authorization")

		if request.URL.Path == "/sumapi/v1/auth" {
			next.ServeHTTP(respWriter, request)

			return
		}

		if !strings.HasPrefix(auth, "Bearer ") {
			respWriter.Header().Add("WWW-Authenticate", "Bearer")
			common.WriteError(respWriter, http.StatusUnauthorized, "UNAUTHORIZED", "")

			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		_, err := tokenHelperSrv.VerifyToken(ctx, token)
		if err != nil {
			log.Printf("failed to verify token: %v", err)
			common.WriteError(respWriter, http.StatusUnauthorized, "INVALID_TOKEN", "auth token invalid")

			return
		}

		next.ServeHTTP(respWriter, request)
	})
}

func InstallRoutes(r chi.Router) {
	r.Route("/sumapi/v1", func(router chi.Router) {
		router.NotFound(func(w http.ResponseWriter, r *http.Request) {
			common.WriteError(w, http.StatusNotFound, "NOT_FOUND", "not found")
		})
		router.Use(validateToken)
		router.Post("/auth", handleAuth)
		router.Post("/sum", handleSum)
	})
}
