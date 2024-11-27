package res

import (
	"encoding/json"
	"net/http"
)

type statusRes struct {
	Status string `json:"status"`
}

type statusDataRes struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type statusMessageRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type statusDataMessageRes struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func BadRequest(w http.ResponseWriter, error string) {
	Error(w, error, http.StatusBadRequest)
}

func Error(w http.ResponseWriter, error string, code int) {
	Write(w, code, "error", error, nil)
}

func Fail(w http.ResponseWriter, message string, code int) {
	Write(w, code, "fail", message, nil)
}

func Success(w http.ResponseWriter, data interface{}) {
	Write(w, http.StatusOK, "success", "", data)
}

func Write(w http.ResponseWriter, code int, status, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	var err error
	if data == nil && message == "" {
		err = json.NewEncoder(w).Encode(statusRes{status})
	} else if data == nil {
		err = json.NewEncoder(w).Encode(statusMessageRes{status, message})
	} else if message == "" {
		err = json.NewEncoder(w).Encode(statusDataRes{status, data})
	} else {
		err = json.NewEncoder(w).Encode(statusDataMessageRes{status, message, data})
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
