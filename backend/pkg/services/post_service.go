package services

import (
	"encoding/json"
	"net/http"
	"project-root/pkg/api"
	"time"
)

var apiVersion = "1.0"

// HTTPError creates a structured error response
func HTTPError(w http.ResponseWriter, status int, userMessage, logMessage string, authenticated bool, user *api.UserResponse, fieldErrors []api.ValidationError) {
    response := api.Response{
        Status:        "error",
        Message:       userMessage,
        Error: &api.ErrorDetails{
            Code:    status,
            Message: logMessage,
            Details: fieldErrors,
        },
        Authenticated: authenticated,
        User:          user,
        Metadata: api.Metadata{
            Timestamp: time.Now(),
            Version:   apiVersion,
        },
    }
    RespondWithJSON(w, status, response)
}

// RespondWithSuccess sends a JSON response with a success message
func RespondWithSuccess(w http.ResponseWriter, status int, message string, authenticated bool, payload , pagination interface{}, user *api.UserResponse) {
	response := api.Response{
		Status:        "success",
		Message:       message,
		Payload:       payload,
		Pagination: pagination,
		Authenticated: authenticated,
		User:          user,
		Metadata: api.Metadata{
			Timestamp: time.Now(),
			Version:   apiVersion,
		},
	}
	RespondWithJSON(w, status, response)
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}