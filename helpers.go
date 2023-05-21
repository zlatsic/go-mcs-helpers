package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const MaxBytes = 1_000_000

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	data    any    `json:"data,omitempty"`
}

func ReadJSON(writer http.ResponseWriter, request *http.Request, data any) error {
	request.Body = http.MaxBytesReader(writer, request.Body, int64(MaxBytes))

	dec := json.NewDecoder(request.Body)

	if err := dec.Decode(data); err != nil {
		return err
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func WriteJSON(writer http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			writer.Header()[key] = value
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_, err = writer.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJSON(writer http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return WriteJSON(writer, statusCode, payload)
}
