package utilities

import (
	"encoding/json"
	"net/http"
)

func ErrorJsonRespond(w http.ResponseWriter, code int, err error) {
	RespondJson(w, code, map[string]string{"error": err.Error()})
}

func RespondJson(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
