package respond

import (
	"encoding/json"
	"log"
	"net/http"
)

type WithJSONOptions struct {
	Code int
	Body interface{}
}

func WithJSON(w http.ResponseWriter, o WithJSONOptions) {
	data, err := json.Marshal(o.Body)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(o.Code)
	w.Write(data)
}

type WithErrorOptions struct {
	Code    int
	Message string
}

func WithError(w http.ResponseWriter, o WithErrorOptions) {
	if o.Message == "" {
		WithStatus(w, o.Code)
		return
	}

	body := struct {
		Error string `json:"error"`
	}{
		Error: o.Message,
	}

	WithJSON(w, WithJSONOptions{
		Code: o.Code,
		Body: body,
	})
}

func WithStatus(w http.ResponseWriter, code int) {
	body := struct {
		Message string `json:"message"`
	}{
		Message: http.StatusText(code),
	}

	WithJSON(w, WithJSONOptions{
		Code: code,
		Body: body,
	})
}
