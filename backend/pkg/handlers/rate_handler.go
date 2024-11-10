package handlers

import (
	"encoding/json"
	"net/http"
	"project-root/pkg/api"
	"project-root/pkg/repositories"
	"project-root/pkg/services"
	"log"
)

// @Summary Rate a post or comment.
// @Tags rate
// @Accept json
// @Produce json
// @Security BearerAuth

// @Param body body api.RateRequest true "Rate details"
// @Success 200 {object} api.Response{payload=api.RateResponse} "Rating updated successfully"
// @Failure 401 {object} api.Response{error=api.ErrorDetails} "Unauthorized"
// @Failure 500 {object} api.Response{error=api.ErrorDetails} "Internal Server Error"
// @Router /rate [PUT]
func HandleRate(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var user *api.UserResponse

	user, authenticated = services.AuthenticateUser(r)
	if !authenticated {
		services.HTTPError(w, http.StatusUnauthorized, "Unauthorized", "User is not authenticated", false, nil, nil)
		return
	}
	
	var rateForm api.RateRequest
	if err := json.NewDecoder(r.Body).Decode(&rateForm); err != nil {
		services.HTTPError(w, http.StatusBadRequest, "Bad Request", "Invalid request payload", authenticated, user, nil)
		return
	}
	log.Printf("Rate form: %v", rateForm)

	rate, status, err := repositories.UpdateRate(user.ID, rateForm)
	if err != nil {
		services.HTTPError(w, http.StatusInternalServerError, "Internal Server Error", "Error updating rating", authenticated, user, nil)
		return
	}

	payload := api.RateResponse{
		Rate: api.Rate{
			Rate:   rate,
			Status: status},
	}

	services.RespondWithJSON(w, http.StatusOK, api.Response{
		Status:        "success",
		Message:       "Rating updated successfully",
		Authenticated: authenticated,
		Payload:       payload,
	})
}
