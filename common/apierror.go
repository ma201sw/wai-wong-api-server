package common

import (
	"encoding/json"
	"log"
	"net/http"
)

func Error(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w)
}

func WriteError(respWriter http.ResponseWriter, status int, code, desc string) {
	errVal := &struct {
		HTTPStatus int    `json:"http_status"`
		Code       string `json:"code"`
		Desc       string `json:"desc"`
	}{
		HTTPStatus: status,
		Code:       code,
		Desc:       desc,
	}

	errValBytes, err := json.Marshal(errVal)
	if err != nil {
		log.Printf("marshal error: %v", err)

		return
	}

	respWriter.WriteHeader(status)
	_, respWriterWriteErr := respWriter.Write(errValBytes)

	if respWriterWriteErr != nil {
		log.Printf("response writer write error error: %v", respWriterWriteErr)

		return
	}
}

func WriteInternalError(w http.ResponseWriter) {
	WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
}
