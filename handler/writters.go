package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func writeOk(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("ok"))

}

func writeJSON(rw http.ResponseWriter, data any) {
	marshal, err := json.Marshal(data)
	if err != nil {
		writeBadRequest(nil, rw, fmt.Errorf("incorrect data input, %w", err))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(marshal)
}

func writeBadRequest(req *http.Request, rw http.ResponseWriter, data error) {
	req = req.WithContext(context.WithValue(req.Context(), "error", data.Error()))
	http.Error(rw, data.Error(), http.StatusInternalServerError)
}

func writeInternalError(req *http.Request, rw http.ResponseWriter, data error) {
	req = req.WithContext(context.WithValue(req.Context(), "error", data.Error()))
	http.Error(rw, "internal server error", http.StatusInternalServerError)
}

func writeMethodNotAllowed(req *http.Request, rw http.ResponseWriter, msg string) {
	req = req.WithContext(context.WithValue(req.Context(), "error", msg))
	http.Error(rw, msg, http.StatusMethodNotAllowed)
}
