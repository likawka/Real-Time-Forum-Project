package websockets

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"project-root/pkg/api"
	"project-root/pkg/repositories"
	"project-root/pkg/services"

	"github.com/gorilla/websocket"
)

const (
	pingInterval        = 30 * time.Second
	activeUsersInterval = 5 * time.Second
)

type WebSocketManager struct {
	clients map[*Client]bool
	rooms   map[string]map[int]*Client
	mu      sync.RWMutex
}

func NewWebSocketManager() *WebSocketManager {
	manager := &WebSocketManager{
		clients: make(map[*Client]bool),
		rooms:   make(map[string]map[int]*Client),
	}

	go manager.periodicActiveUsersBroadcast()

	return manager
}

func (manager *WebSocketManager) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()

	user, status := services.AuthenticateUser(r)
	if !status {
		return
	}
	client := &Client{
		conn:    conn,
		manager: manager,
		user:    user,
		stopper: make(chan struct{}),
	}

	manager.addClient(client)
	defer manager.removeClient(client)

	go client.writeMessages()

	client.readMessages()
}

func (manager *WebSocketManager) periodicActiveUsersBroadcast() {
	ticker := time.NewTicker(activeUsersInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			manager.broadcastActiveUsers()
		}
	}
}

func (manager *WebSocketManager) broadcastActiveUsers() {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	activeUsers := make([]*api.UserResponse, 0)
	for client := range manager.clients {
		activeUsers = append(activeUsers, client.user)
	}

	msg := api.MessageResponse{
		Type: "active_users",
		Payload: map[string]interface{}{
			"users": activeUsers,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling active users:", err)
		return
	}

	for client := range manager.clients {
		client.sendMessage(data)
	}
}

func (manager *WebSocketManager) broadcastMessage(msg *api.MessageMessage, roomHash string, sender *Client) {

	msg.SendAt = time.Now()
	msg.Sender = sender.user
	msg.RoomHash = roomHash // Ensure the room hash is set in the message

	err := manager.saveMessageToDB(sender.user.ID, msg)
	if err != nil {
		sender.sendError("Failed to save message")
		// return
	}

	// Get clients in the room
	clients := manager.getRoomClients(roomHash)
	if len(clients) == 0 {
		log.Printf("No clients in room %s to broadcast message", roomHash)
	}

	for _, client := range clients {
		response := api.MessageResponse{
			Type:    "message",
			Payload: msg,
		}
		data, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshalling message:", err)
			continue
		}
		if err := client.sendMessage(data); err != nil {
			log.Println("Error sending message to client:", err)
		}
	}

	// if sender != nil {
	//     response := api.MessageResponse{
	//         Type:    "message_sent",
	//         Payload: msg,
	//     }
	//     data, err := json.Marshal(response)
	//     if err != nil {
	//         log.Println("Error marshalling confirmation message:", err)
	//         return
	//     }
	//     if err := sender.sendMessage(data); err != nil {
	//         log.Println("Error sending confirmation message to sender:", err)
	//     }
	// }
}

func (manager *WebSocketManager) broadcastTypingStatus(msg *api.TypingMessage, sender *Client) {
	clients := manager.getRoomClients(msg.RoomHash)
	msg.Sender = sender.user
	response := api.MessageResponse{
		Type:    "typing",
		Payload: msg,
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshalling typing status:", err)
		return
	}

	for _, client := range clients {
		client.sendMessage(data)
	}
}

func (manager *WebSocketManager) addClient(client *Client) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.clients[client] = true
}

func (manager *WebSocketManager) removeClient(client *Client) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if _, ok := manager.clients[client]; ok {
		client.conn.Close()
		delete(manager.clients, client)
	}

	for roomHash, clients := range manager.rooms {
		if _, ok := clients[client.user.ID]; ok {
			delete(clients, client.user.ID)
		}
		if len(clients) == 0 {
			delete(manager.rooms, roomHash)
		}
	}
}

func (manager *WebSocketManager) addClientToRoom(roomHash string, client *Client) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if !manager.hasAccessToRoom(client.user, roomHash) {
		client.sendError("Access denied to room " + roomHash)
		return
	}

	if _, ok := manager.rooms[roomHash]; !ok {
		manager.rooms[roomHash] = make(map[int]*Client)
	}
	manager.rooms[roomHash][client.user.ID] = client
}

func (manager *WebSocketManager) getRoomClients(roomHash string) []*Client {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	var clients []*Client
	if roomClients, ok := manager.rooms[roomHash]; ok {
		for _, client := range roomClients {
			clients = append(clients, client)
		}
	}
	return clients
}

// Mock function to check if the user has access to the room
func (manager *WebSocketManager) hasAccessToRoom(user *api.UserResponse, roomHash string) bool {
	chatRepo := repositories.NewChatRepository()
	hasAccess, err := chatRepo.CheckChatAccess(user.ID, roomHash)
	if err != nil {
		log.Println("Error checking chat access:", err)
		return false
	}
	return hasAccess
}

func (manager *WebSocketManager) saveMessageToDB(userID int, message *api.MessageMessage) error {
	messageRepo := repositories.NewChatRepository()
	if err := messageRepo.SaveMessage(message.RoomHash, userID, message.Message); err != nil {
		log.Println("Error saving message to DB:", err)
		return err
	}
	return nil
}
