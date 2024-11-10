package websockets

import (
	"encoding/json"
	"log"
	"time"

	"project-root/pkg/api"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	manager *WebSocketManager
	user    *api.UserResponse
	stopper chan struct{}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var msg api.MessageResponse
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Error unmarshalling JSON:", err)
			c.sendError("Invalid message format")
			continue
		}

		payload, err := json.Marshal(msg.Payload)

		switch msg.Type {
		case "join_room":
			var joinMsg api.JoinRoomMessage
			if err := json.Unmarshal(payload, &joinMsg); err != nil {
				log.Println("Error unmarshalling join room message:", err)
				c.sendError("Invalid join room message")
				continue
			}
			c.handleJoinRoom(&joinMsg)
		case "message":
			var msgMsg api.MessageMessage
			if err := json.Unmarshal(payload, &msgMsg); err != nil {
				log.Println("Error unmarshalling message message:", err)
				c.sendError("Invalid message message")
				continue
			}
			c.handleMessage(&msgMsg)
		case "typing":
			var typingMsg api.TypingMessage
			if err := json.Unmarshal(payload, &typingMsg); err != nil {
				log.Println("Error unmarshalling typing message:", err)
				c.sendError("Invalid typing message")
				continue
			}
			c.handleTyping(&typingMsg)
		default:
			c.sendError("Unknown message type")
		}
	}
}

func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-c.stopper:
			return
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Error sending ping:", err)
				return
			}
		}
	}
}

func (client *Client) sendMessage(message []byte) error {
	if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println("Error sending message:", err)
		return err
	}
	return nil
}

func (c *Client) sendError(errorMessage string) {
	msg := api.MessageResponse{
		Type:    "error",
		Payload: api.ErrorMessage{Message: errorMessage},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling error message:", err)
		return
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println("Error sending error message:", err)
	}
}

func (c *Client) handleJoinRoom(msg *api.JoinRoomMessage) {
	c.manager.addClientToRoom(msg.RoomHash, c)
}

func (c *Client) handleMessage(msg *api.MessageMessage) {
	c.manager.broadcastMessage(msg, msg.RoomHash, c)
}

func (c *Client) handleTyping(msg *api.TypingMessage) {
	c.manager.broadcastTypingStatus(msg, c)
}
