package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type errorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, errorBody{Code: status, Message: msg})
}
