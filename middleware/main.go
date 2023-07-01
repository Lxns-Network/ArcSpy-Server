package middleware

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ApplyMiddleware(next http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		next = middleware(next)
	}

	return next
}

func RespondWithJSON(w http.ResponseWriter, code int, message string, data interface{}) {
	if code != 200 {
		w.WriteHeader(code)
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
	_, _ = w.Write(response)
}
