package repositories

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"project-root/pkg/api"
	"time"
)

// ChatRepository interacts with chat-related data sources.
type ChatRepository struct {
	DB    *sql.DB
	MsgDB *sql.DB
}

// NewChatRepository creates a new ChatRepository instance with multiple DB connections if needed.
func NewChatRepository() *ChatRepository {
	return &ChatRepository{DB: dbHandler.MainDB, MsgDB: dbHandler.MsgDB}
}

// GenerateChatHash generates a unique hash for the chat based on user IDs and timestamp.
func GenerateChatHash(user1ID, user2ID int) string {
	// Ensure the lower user ID comes first for consistency
	var id1, id2 int
	if user1ID < user2ID {
		id1, id2 = user1ID, user2ID
	} else {
		id1, id2 = user2ID, user1ID
	}

	// Generate the hash using SHA-256
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d-%d-%d", id1, id2, time.Now().UnixNano())))
	return hex.EncodeToString(hash.Sum(nil))
}

// CreateChat creates a new chat with a corresponding table for storing messages.
func (repo *ChatRepository) CreateChat(user1ID, user2ID int) (string, error) {
	// Check if the chat already exists
	var existingHash string
	query := `
		SELECT hash FROM conversations 
		WHERE (user1_id = ? AND user2_id = ?) 
		   OR (user1_id = ? AND user2_id = ?);
	`
	err := repo.DB.QueryRow(query, user1ID, user2ID, user2ID, user1ID).Scan(&existingHash)
	if err == nil {
		// Chat already exists
		return existingHash, nil
	} else if err != sql.ErrNoRows {
		return "", fmt.Errorf("error checking existing chat: %v", err)
	}

	// Generate a new chat hash
	chatHash := GenerateChatHash(user1ID, user2ID)

	// Create a new conversation in the main database
	insertQuery := `
		INSERT INTO conversations (user1_id, user2_id, hash)
		VALUES (?, ?, ?);
	`
	_, err = repo.DB.Exec(insertQuery, user1ID, user2ID, chatHash)
	if err != nil {
		return "", fmt.Errorf("error creating conversation: %v", err)
	}

	// Create a corresponding messages table in the MsgDB
	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS "%s" (
			"sender_id" INTEGER NOT NULL,
			"message_content" TEXT NOT NULL,
			"sent_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`, chatHash)

	_, err = repo.MsgDB.Exec(createTableQuery)
	if err != nil {
		return "", fmt.Errorf("error creating chat messages table: %v", err)
	}

	return chatHash, nil
}

// GetChatsForUser retrieves all chats for a specific user.
func (repo *ChatRepository) GetChatsForUser(userID int) ([]api.ChatInfo, error) {
	query := `
		SELECT c.hash, u1.id AS user1_id, u2.id AS user2_id
		FROM conversations c
		JOIN users u1 ON c.user1_id = u1.id
		JOIN users u2 ON c.user2_id = u2.id
		WHERE u1.id = ? OR u2.id = ?;
	`

	rows, err := repo.DB.Query(query, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying chats: %v", err)
	}
	defer rows.Close()

	var chats []api.ChatInfo
	for rows.Next() {
		var chat api.ChatInfo
		if err := rows.Scan(&chat.ChatHash, &chat.User1ID, &chat.User2ID); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return chats, nil
}

// ErrChatNotFound is returned when a chat is not found.
var ErrChatNotFound = fmt.Errorf("chat not found")

// CheckChatAccess checks if the chat exists and if the user has access to it.
func (repo *ChatRepository) CheckChatAccess(userID int, chatHash string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM conversations c
		WHERE c.hash = ? AND (c.user1_id = ? OR c.user2_id = ?);
	`

	var count int
	err := repo.DB.QueryRow(query, chatHash, userID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking chat access: %v", err)
	}

	if count == 0 {
		return false, ErrChatNotFound
	}

	return true, nil
}

// GetChatDetails retrieves chat details if the user has access to it.
func (repo *ChatRepository) GetChatDetails(userID int, chatHash string) (*api.Chat, error) {
	// Check if the chat exists and if the user is a participant
	hasAccess, err := repo.CheckChatAccess(userID, chatHash)
	if err != nil {
		if err == ErrChatNotFound {
			return nil, ErrChatNotFound
		}
		return nil, fmt.Errorf("error checking chat access: %v", err)
	}

	if !hasAccess {
		return nil, ErrChatNotFound
	}

	query := `
		SELECT c.hash, c.user1_id, c.user2_id
		FROM conversations c
		WHERE c.hash = ? AND c.user1_id = ? OR c.user2_id = ?;
	`

	var chat api.Chat
	err = repo.DB.QueryRow(query, chatHash, userID, userID).Scan(&chat.ChatInfo.ChatHash, &chat.ChatInfo.User1ID, &chat.ChatInfo.User2ID)
	if err == sql.ErrNoRows {
		return nil, ErrChatNotFound
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving chat details: %v", err)
	}

	// Retrieve messages for the chat
	messages, err := repo.GetMessagesForChat(chatHash)
	if err != nil {
		return nil, fmt.Errorf("error retrieving messages: %v", err)
	}
	chat.Message = messages

	return &chat, nil
}

// GetMessagesForChat retrieves all messages for a specific chat.
func (repo *ChatRepository) GetMessagesForChat(chatHash string) ([]api.MessageMessage, error) {
	sanitizedChatHash := fmt.Sprintf(`"%s"`, chatHash)

	query := fmt.Sprintf(`
	SELECT m.sender_id, m.message_content, m.sent_at
	FROM %s m
	ORDER BY m.sent_at;
	`, sanitizedChatHash)

	rows, err := repo.MsgDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying messages: %v", err)
	}
	defer rows.Close()

	var messages []api.MessageMessage

	for rows.Next() {
		var msg api.MessageMessage
		var user api.UserResponse

		if err := rows.Scan(&user.ID, &msg.Message, &msg.SendAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		msg.Sender = &user
		msg.RoomHash = chatHash 

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return messages, nil
}
func (repo *ChatRepository) SaveMessage(chatHash string, senderID  int, messageContent string) error {
	query := fmt.Sprintf(`
		INSERT INTO "%s" (sender_id, message_content, sent_at)
		VALUES (?, ?, CURRENT_TIMESTAMP);
	`, chatHash)

	_, err := repo.MsgDB.Exec(query, senderID, messageContent)
	if err != nil {
		return fmt.Errorf("error saving message: %v", err)
	}

	return nil
}
