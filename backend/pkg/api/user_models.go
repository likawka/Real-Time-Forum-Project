package api

import (
	"time"
)

// User represents a user
type User struct {
	ID               int       `json:"id"`
	Nickname         string    `json:"nickname"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	CreatedAt        time.Time `json:"created_at"`
	AmountOfPosts    int       `json:"amount_of_posts"`
	AmountOfComments int       `json:"amount_of_comments"`
}

// RegistrationRequest represents the request payload for user registration
type RegistrationRequest struct {
	Nickname  string    `json:"nickname" example:"test"`
	Email     string    `json:"email" example:"test@test.com"`
	Password  string    `json:"password" example:"!QAZ2wsx"`
	Age       string    `json:"age" example:"02.02.2002"`
	Gender    string    `json:"gender" example:"male"`
	FirstName string    `json:"first_name" example:"John"`
	LastName  string    `json:"last_name" example:"Doe"`
	CreatedAt time.Time `json:"created_at" swaggerignore:"true"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" example:"test@test.com"`
	Password string `json:"password" example:"!QAZ2wsx"`
}

// // UpdateUserRequest represents the request payload for updating user information
// type UpdateUserRequest struct {
// 	Nickname  string `json:"nickname"`
// 	Age       string `json:"age"`
// 	Gender    string `json:"gender"`
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// 	Email     string `json:"email"`
// 	Password  string `json:"password"`
// }

// UserResponse represents a user
type UserResponse struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
}

type GetUserResponse struct {
	User     User      `json:"user"`
	Posts    []Post    `json:"posts"`
	Comments []Comment `json:"comments"`
}

type GetUsersResponse struct {
	Users []Users `json:"users"`
}

type Users struct {
	ID           int       `json:"id"`
	Nickname     string    `json:"nickname"`
	LastActivity time.Time `json:"last_activity"`
}
