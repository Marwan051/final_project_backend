package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// DecodeJSONBody decodes and sanitizes JSON request bodies with proper error handling
func DecodeJSONBody(r *http.Request, dst any) error {
	// Limit request body size to prevent DoS attacks (10MB max)
	r.Body = http.MaxBytesReader(nil, r.Body, 10*1024*1024)

	// Decode JSON with strict settings
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return errors.New("malformed JSON at position " + string(rune(syntaxError.Offset)))

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("malformed JSON")

		case errors.As(err, &unmarshalTypeError):
			return errors.New("invalid value for field '" + unmarshalTypeError.Field + "'")

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return errors.New("unknown field " + fieldName)

		case errors.Is(err, io.EOF):
			return errors.New("request body must not be empty")

		case err.Error() == "http: request body too large":
			return errors.New("request body too large (max 10MB)")

		default:
			return err
		}
	}

	// Check for additional JSON data after the first object
	if dec.More() {
		return errors.New("request body must contain only a single JSON object")
	}

	return nil
}

// WriteJSONResponse writes a JSON response with proper headers
func WriteJSONResponse(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// WriteJSONError writes a JSON error response
func WriteJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
