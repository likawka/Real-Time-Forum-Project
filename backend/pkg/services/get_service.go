package services

import (
	"net/http"
	"project-root/pkg/api"
	"project-root/pkg/repositories"
	"github.com/gorilla/mux"
)
func GetRouteParams(r *http.Request) map[string]string {
    // Використовує mux.Vars для отримання параметрів шляху
    return mux.Vars(r)
}

// GetSessionID retrieves the session ID from the request cookie
func getSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// AuthenticateUser retrieves user information based on session ID.
func AuthenticateUser(r *http.Request) (*api.UserResponse, bool) {
	authenticated := false
	var user *api.UserResponse

	sessionID, err := getSessionID(r)
	if err != nil {
		return &api.UserResponse{ID: 0, Nickname: ""}, false
	}

	u, err := repositories.GetUserBySessionID(sessionID)
	if err != nil || u == nil {
		return &api.UserResponse{ID: 0, Nickname: ""}, false
	}

	user = &api.UserResponse{
		ID:       u.ID,
		Nickname: u.Nickname,
	}
	authenticated = true

	return user, authenticated
}