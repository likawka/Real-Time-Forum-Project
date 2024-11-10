package handlers

import (
	"log"
	"net/http"

	"encoding/json"
	"github.com/gorilla/mux"
	"project-root/pkg/api"
	"project-root/pkg/repositories"
	"project-root/pkg/services"
)

// HandleGetChats retrieves all chats for a specific user.
// @Summary Get all chats for a user
// @Description Retrieves all chats that the authenticated user is a participant of.
// @Tags chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} api.ChatsList "Successfully retrieved chats"
// @Router /chats [get]
func HandleGetChats(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var user *api.UserResponse

	// Authenticate the user
	user, authenticated = services.AuthenticateUser(r)
	if !authenticated {
		services.HTTPError(w, http.StatusUnauthorized, "Unauthorized", "User not authenticated", false, nil, nil)
		return
	}

	// Create a new Chat repository
	chatRepo := repositories.NewChatRepository()

	// Get the user's chats
	chats, err := chatRepo.GetChatsForUser(user.ID)
	if err != nil {
		log.Println("Error getting chats:", err)
		services.HTTPError(w, http.StatusInternalServerError, "Internal Server Error", "Error retrieving chats", authenticated, user, nil)
		return
	}

	// Prepare response payload
	payload := api.ChatsList{
		Chats: chats,
	}

	// Send the success response
	services.RespondWithSuccess(w, http.StatusOK, "Chats retrieved successfully", authenticated, payload, nil, user)
}

// HandleCreateChat creates a new chat between two users.
// @Summary Create a new chat
// @Description Creates a new chat between two specified users and returns the chat hash.
// @Tags chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param chat body api.CreateChatMessage true "Chat creation details"
// @Success 201 {object} api.ChatCreateResponse "Successfully created chat"
// @Router /chats [post]
func HandleCreateChat(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var user *api.UserResponse

	// Authenticate the user
	user, authenticated = services.AuthenticateUser(r)

	// Parse the request body
	var chat api.CreateChatMessage
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		services.HTTPError(w, http.StatusBadRequest, "Bad Request", "Invalid request payload", authenticated, user, nil)
		return
	}

	// Create a new Chat repository
	chatRepo := repositories.NewChatRepository()

	// Create the chat
	chatHash, err := chatRepo.CreateChat(chat.User1ID, chat.User2ID)
	if err != nil {
		log.Println("Error creating chat:", err)
		services.HTTPError(w, http.StatusInternalServerError, "Internal Server Error", "Error creating chat", authenticated, user, nil)
		return
	}

	// Prepare response payload
	payload := api.ChatCreateResponse{
		ChatHash: chatHash,
	}

	// Send the success response
	services.RespondWithSuccess(w, http.StatusCreated, "Chat created successfully", authenticated, payload, nil, user)
}

// HandleGetChat retrieves a specific chat by its hash.
// @Summary Get chat details by hash
// @Description Retrieves the details of a specific chat and its messages if the user has access to it.
// @Tags chats
// @Produce json
// @Security BearerAuth
// @Param chatHash path string true "Chat hash"
// @Success 200 {object} api.Chat "Successfully retrieved chat details"
// @Router /chats/{chatHash} [get]
// HandleGetChat handles the request to get a specific chat by hash.
func HandleGetChat(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	var user *api.UserResponse

	// Authenticate the user
	user, authenticated = services.AuthenticateUser(r)
	if !authenticated {
		services.HTTPError(w, http.StatusUnauthorized, "Unauthorized", "User not authenticated", false, nil, nil)
		return
	}

	// Get the chat hash from the URL parameters
	vars := mux.Vars(r)
	chatHash := vars["chatHash"]

	// Create a new Chat repository
	chatRepo := repositories.NewChatRepository()

	// Check if the user has access to the chat and get chat details
	chat, err := chatRepo.GetChatDetails(user.ID, chatHash)
	if err != nil {
		if err == repositories.ErrChatNotFound {
			services.HTTPError(w, http.StatusNotFound, "Not Found", "Chat not found", authenticated, user, nil)
		} else {
			services.HTTPError(w, http.StatusInternalServerError, "Internal Server Error", err.Error(), authenticated, user, nil)
		}
		return
	}

	userID_1, _ := repositories.GetUserByID(chat.ChatInfo.User1ID)
	userID_2, _ := repositories.GetUserByID(chat.ChatInfo.User2ID)

	for i := 0; i < len(chat.Message); i++ {
		if chat.Message[i].Sender.ID == chat.ChatInfo.User1ID {
			chat.Message[i].Sender.Nickname = userID_1.Nickname
		} else {
			chat.Message[i].Sender.Nickname = userID_2.Nickname
		}
	}

	// Send the success response
	services.RespondWithSuccess(w, http.StatusOK, "Chat retrieved successfully", authenticated, chat, nil, user)
}
