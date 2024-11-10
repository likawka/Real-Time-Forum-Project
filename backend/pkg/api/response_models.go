package api

import "time"

// Response is the standard response format
type Response struct {
	Status        string        `json:"status"`
	Message       string        `json:"message"`
	Error         *ErrorDetails `json:"error,omitempty" swaggerignore:"true"`
	Payload       interface{}   `json:"payload,omitempty" swaggerignore:"true"`
	Pagination    interface{}   `json:"pagination,omitempty" swaggerignore:"true"`
	Authenticated bool          `json:"authenticated"`
	User          *UserResponse `json:"user,omitempty"`
	Metadata      Metadata      `json:"metadata"`
}

// ErrorDetails provides detailed information about an error
type ErrorDetails struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Details []ValidationError `json:"details,omitempty"`
}

// FieldError provides detailed information about a specific field error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ChatPagination provides pagination details specific to chat messages
type ChatPagination struct {
	PerPage     int    `json:"per_page"` // Кількість повідомлень на сторінці
	FromID     int64  `json:"from_id,omitempty"` // Початковий ID для фільтрації
	ToID       int64  `json:"to_id,omitempty"` // Кінцевий ID для фільтраці
}

// GeneralPagination provides pagination details for other types of data (e.g., posts)
type GeneralPagination struct {
	TotalCount int    `json:"total_count"` // Загальна кількість елементів
	CurrentPage int    `json:"current_page"` // Поточна сторінка
	PerPage     int    `json:"per_page"` // Кількість елементів на сторінці
	TotalPages  int    `json:"total_pages"` // Загальна кількість сторінок
	OrderBy     string `json:"order_by"` // Тип сортування або порядок
}

// Metadata provides additional metadata about the response
type Metadata struct {
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}
