package api

import (
	"time"
)

type Post struct {
	ID               int        `json:"id"`
	UserID           int        `json:"user_id"`
	Nickname         string     `json:"nickname"`
	Title            string     `json:"title"`
	Content          string     `json:"content"`
	CreatedAt        time.Time  `json:"created_at"`
	AmountOfComments int        `json:"amount_of_comments"`
	Rate             Rate        `json:"rate"`
	Categories       []Category `json:"categories"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Nickname  string    `json:"nickname"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Rate      Rate       `json:"rate"`
}

type PostCreateRequest struct {
	UserID     int       `json:"user_id" swaggerignore:"true"`
	ID         int       `json:"id" swaggerignore:"true"`
	Title      string    `json:"title" example:"Test title"`
	Content    string    `json:"content" example:"Test content"`
	Categories string    `json:"categories" example:"#category1 #category2"`
	CreatedAt  time.Time `json:"created_at" swaggerignore:"true"`
}
type PostCreateResponse struct {
	ID int `json:"id"`
}

type CommentCreateRequest struct {
	ID        int       `json:"id" swaggerignore:"true"`
	PostID    int       `json:"post_id" swaggerignore:"true"`
	UserID    int       `json:"user_id" swaggerignore:"true"`
	Content   string    `json:"content" example:"Test comment"`
	CreatedAt time.Time `json:"created_at" swaggerignore:"true"`
}

type CommentResponse struct {
	Comment Comment `json:"comment"`
}

type PostsResponse struct {
	Posts []Post `json:"posts"`
}

type PostAndCommentsResponse struct {
	Post     Post       `json:"post"`
	Comments *[]Comment `json:"comments"`
}
