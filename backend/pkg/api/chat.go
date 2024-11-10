package api

type ChatsList struct {
	Chats []ChatInfo `json:"chats"`
}

type ChatInfo struct {
	ChatHash string `json:"chatHash"`
	User1ID int `json:"user1_id"`
	User2ID int `json:"user2_id"`
}

type Chat struct {
	ChatInfo ChatInfo `json:"ChatInfo"`
	Message []MessageMessage `json:"messages"`
}

type CreateChatMessage struct {
	User1ID int `json:"user1_id"`
	User2ID int `json:"user2_id"`
}

type ChatCreateResponse struct {
	ChatHash string `json:"chatHash"`
}