package repositories

import (
	"database/sql"
	"project-root/pkg/api"
	"log"
	"time"
)

type CommentRepository struct {
	DB *sql.DB
}

// NewCommentRepository creates a new CommentRepository instance.
func NewCommentRepository() *CommentRepository {
	return &CommentRepository{DB: dbHandler.MainDB}
}

// Create inserts a new comment into the database.
func (r *CommentRepository) Create(c *api.CommentCreateRequest) error {
	stmt, err := dbHandler.MainDB.Prepare("INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	c.CreatedAt = time.Now()
	result, err := stmt.Exec(c.PostID, c.UserID, c.Content, c.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = int(id)
	UpdatePostAmount(c.PostID)
	IncrementCount(c.UserID, "comments")
	return nil
}

func (r *CommentRepository) GetComments(page, pageSize int, sortBy, filterType string, filterValue interface{}, userIdAuth int) (*[]api.Comment, int, int, error) {
	baseQuery := `
        SELECT 
            comments.id, 
            comments.post_id, 
            comments.user_id, 
            users.nickname, 
            comments.content, 
            comments.created_at, 
            comments.rate 
        FROM comments
        JOIN users ON comments.user_id = users.id
    `
	var countQuery string
	var args []interface{}

	switch filterType {
	case "userID":
		baseQuery += " WHERE comments.user_id = ?"
		countQuery = "SELECT COUNT(*) FROM comments WHERE user_id = ?"
		args = append(args, filterValue)
	case "nickname":
		baseQuery += " WHERE users.nickname = ?"
		countQuery = "SELECT COUNT(*) FROM comments JOIN users ON comments.user_id = users.id WHERE users.nickname = ?"
		args = append(args, filterValue)
	default:
		baseQuery += " WHERE comments.post_id = ?"
		countQuery = "SELECT COUNT(*) FROM comments WHERE post_id = ?"
		args = append(args, filterValue)
	}

	if sortBy != "" {
		baseQuery += " ORDER BY " + sortBy
	} else {
		baseQuery += " ORDER BY comments.created_at DESC"
	}

	baseQuery += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, (page-1)*pageSize)


	var totalItems int
	err := dbHandler.MainDB.QueryRow(countQuery, filterValue).Scan(&totalItems)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := totalItems / pageSize
	if totalItems%pageSize != 0 {
		totalPages++
	}


	rows, err := dbHandler.MainDB.Query(baseQuery, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var comments []api.Comment
	for rows.Next() {
		var comment api.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Nickname, &comment.Content, &comment.CreatedAt, &comment.Rate.Rate); err != nil {
			return nil, 0, 0, err
		}

		if userIdAuth != 0 {
			status, err := getRateStatus("comment", comment.ID, userIdAuth)
			if err != nil && err != sql.ErrNoRows {
				return nil, 0, 0, err
			}
			comment.Rate.Status = status
		}

		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	return &comments, totalItems, totalPages, nil
}