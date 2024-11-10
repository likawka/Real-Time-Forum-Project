package api

type Rate struct {
	Rate   int    `json:"rate"`
	Status string `json:"status"`
}

type RateResponse struct {
	Rate Rate `json:"rate"`
}

type RateRequest struct {
	PostID    int    `json:"post_id" example:"1"`
	CommentID int    `json:"comment_id" example:"2"`
	Status    string `json:"status" example:"up"`
}
