package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kymnguyen/mvta/services/vehicle/internal/application/dto"
)

func respondSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func decodeJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func scanInt(s string, v *int) (int, error) {
	n, err := fmt.Sscanf(s, "%d", v)
	if err != nil || n != 1 {
		return 0, fmt.Errorf("invalid integer: %s", s)
	}
	return *v, nil
}
