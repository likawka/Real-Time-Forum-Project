package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"project-root/pkg/api"
	"project-root/pkg/repositories"
	"project-root/pkg/services"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HandleRegister registers a new user.
// @Summary Register a new user
// @Description Registers a new user with the provided registration details.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body api.RegistrationRequest true "Registration details"
// @Success 200 {object} api.UserResponse "User registered successfully"
// @Failure 400 {object} api.Response "Invalid request payload"
// @Failure 409 {object} api.Response "User already registered"
// @Router /auth/register [post]
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var registrForm api.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&registrForm); err != nil {
		services.HTTPError(w, http.StatusBadRequest, "Invalid request payload", "Invalid request payload", authenticated, nil, nil)
		return
	}

	// Validate fields before proceeding
	validationErrors := services.ValidateOperation("registration", registrForm)
	if len(validationErrors) > 0 {
		// Handle the HTTP error response
		services.HTTPError(w, http.StatusBadRequest, "Validation error", "Validation error", authenticated, nil, validationErrors)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrForm.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		services.HTTPError(w, http.StatusInternalServerError, "Internal server error", "Error hashing password", authenticated, nil, nil)
		return
	}
	registrForm.Password = string(hashedPassword)
	registrForm.CreatedAt = time.Now()

	// Create the user and get the UserResponse
	userResponse, err := repositories.CreateUser(&registrForm)
	if err != nil {
		log.Println("Error creating user:", err)
		services.HTTPError(w, http.StatusConflict, "User already registered", "Error creating user", authenticated, nil, nil)
		return
	}

	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	session := api.Session{
		UserID:    userResponse.ID,
		SessionID: sessionID,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	if err := repositories.CreateSession(&session); err != nil {
		log.Println("Error creating session for new user:", err)
		services.HTTPError(w, http.StatusInternalServerError, "Internal server error", "Error creating session", authenticated, nil, nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Expires:  expiresAt,
		HttpOnly: true,
		Path:     "/",
	})

	authenticated = true

	// Respond with success
	services.RespondWithSuccess(w, http.StatusOK, "User registered successfully", authenticated, nil, nil, userResponse)
}

// HandleLogin logs in a user.
// @Summary Log in a user
// @Description Logs in a user with the provided email and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body api.LoginRequest true "Login credentials"
// @Success 200 {object} api.UserResponse "User logged in successfully"
// @Failure 400 {object} api.Response "Invalid request payload"
// @Failure 401 {object} api.Response "Invalid email or password"
// @Router /auth/login [post]
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var loginForm api.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginForm); err != nil {
		services.HTTPError(w, http.StatusBadRequest, "Invalid request payload", "Invalid request payload", authenticated, nil, nil)
		return
	}

	storedUser, userResponse, err := repositories.ChekUserByEmail(loginForm.Email)
	if err != nil {
		log.Println("User not found:", err)
		services.HTTPError(w, http.StatusUnauthorized, "Invalid email or password", "User not found", authenticated, nil, nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(loginForm.Password)); err != nil {
		log.Println("Invalid password for user:", loginForm.Email)
		services.HTTPError(w, http.StatusUnauthorized, "Invalid email or password", "Invalid password", authenticated, nil, nil)
		return
	}

	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	session := api.Session{
		UserID:    userResponse.ID,
		SessionID: sessionID,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	if err := repositories.CreateSession(&session); err != nil {
		log.Println("Error creating session:", err)
		services.HTTPError(w, http.StatusInternalServerError, "Internal server error", "Error creating session", authenticated, nil, nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Expires:  expiresAt,
		HttpOnly: true,
		Path:     "/",
	})

	authenticated = true
	services.RespondWithSuccess(w, http.StatusOK, "User logged in successfully", authenticated, nil, nil, userResponse)
}

// HandleLogout logs out a user
// @Summary Log out a user
// @Description Logs out the currently authenticated user.
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "User logged out successfully"
// @Failure 401 {object} api.Response "Missing session ID"
// @Router /auth/logout [delete]
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	authenticated := true
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			services.HTTPError(w, http.StatusUnauthorized, "Missing session ID", "Missing session ID", authenticated, nil, nil)
			return
		}
		services.HTTPError(w, http.StatusBadRequest, "Error reading cookie", "Error reading cookie", authenticated, nil, nil)
		return
	}
	sessionID := cookie.Value

	if err := repositories.DeleteSession(sessionID); err != nil {
		log.Println("Error deleting session:", err)
		services.HTTPError(w, http.StatusInternalServerError, "Internal server error", "Error deleting session", authenticated, nil, nil)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	authenticated = false
	services.RespondWithSuccess(w, http.StatusOK, "User logged out successfully", authenticated, nil, nil, nil)
}
