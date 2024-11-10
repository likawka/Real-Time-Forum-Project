package repositories

import (
	"database/sql"
	"log"
	"project-root/pkg/api"
	"time"
)

// PostRepository provides access to the post storage.
type PostRepository struct {
	DB *sql.DB
}

// NewPostRepository creates a new PostRepository instance.
func NewPostRepository() *PostRepository {
	return &PostRepository{DB: dbHandler.MainDB}
}

// GetAllPosts retrieves all posts with pagination, sorting, and optional filtering.
func (r *PostRepository) GetPosts(page, pageSize int, sortBy, filterType, filterValue string) ([]api.Post, int, int, error) {
	query := `
        SELECT 
            posts.id, 
            posts.user_id, 
            users.nickname, 
            posts.title, 
            posts.content, 
            posts.created_at, 
            posts.amount_of_comments, 
            posts.rate 
        FROM posts
        JOIN users ON posts.user_id = users.id
    `

	// Add filter conditions to the query
	if filterType != "" && filterValue != "" {
		switch filterType {
		case "nickname":
			query += " WHERE users.nickname = ?"
		case "user_id":
			query += " WHERE posts.user_id = ?"
		}
	}

	// Add sorting to the query
	if sortBy != "" {
		query += " ORDER BY " + sortBy
	} else {
		query += " ORDER BY posts.created_at DESC"
	}

	query += " LIMIT ? OFFSET ?"

	// Count total items with the same filtering condition
	countQuery := "SELECT COUNT(*) FROM posts"
	if filterType == "nickname" && filterValue != "" {
		countQuery += " JOIN users ON posts.user_id = users.id WHERE users.nickname = ?"
	} else if filterType == "user_id" && filterValue != "" {
		countQuery += " WHERE posts.user_id = ?"
	}

	var totalItems int
	var err error
	if filterType != "" && filterValue != "" {
		err = dbHandler.MainDB.QueryRow(countQuery, filterValue).Scan(&totalItems)
	} else {
		err = dbHandler.MainDB.QueryRow(countQuery).Scan(&totalItems)
	}
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := totalItems / pageSize
	if totalItems%pageSize != 0 {
		totalPages++
	}

	// Execute the query with the pagination and filter parameters
	var rows *sql.Rows
	if filterType != "" && filterValue != "" {
		rows, err = dbHandler.MainDB.Query(query, filterValue, pageSize, (page-1)*pageSize)
	} else {
		rows, err = dbHandler.MainDB.Query(query, pageSize, (page-1)*pageSize)
	}
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var posts []api.Post
	for rows.Next() {
		var post api.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Nickname, &post.Title, &post.Content, &post.CreatedAt, &post.AmountOfComments, &post.Rate.Rate); err != nil {
			return nil, 0, 0, err
		}

		categories, err := r.GetCategoriesByPostID(post.ID)
		if err != nil {
			return nil, 0, 0, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	return posts, totalItems, totalPages, nil
}

// GetPostByID retrieves a post by its ID.
func (r *PostRepository) GetPostByID(postID int, userIdAuth int) (*api.Post, error) {
	var post api.Post
	err := dbHandler.MainDB.QueryRow(`
        SELECT 
            posts.id, 
            posts.user_id, 
            users.nickname, 
            posts.title, 
            posts.content, 
            posts.created_at, 
            posts.amount_of_comments, 
            posts.rate 
        FROM posts
        JOIN users ON posts.user_id = users.id
        WHERE posts.id = ?`, postID).Scan(
		&post.ID, &post.UserID, &post.Nickname, &post.Title, &post.Content, &post.CreatedAt, &post.AmountOfComments, &post.Rate.Rate)
	if err != nil {
		return nil, err
	}

	categories, err := r.GetCategoriesByPostID(post.ID)
	if err != nil {
		return nil, err
	}
	post.Categories = categories

	// If userIdAuth is provided and not zero, fetch rate status
	if userIdAuth != 0 {
		status, err := getRateStatus("post", post.ID, userIdAuth)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("error retrieving rate status: %v", err)
			return nil, err
		}
		post.Rate.Status = status
	}

	return &post, nil
}

// GetCategoriesByPostID retrieves categories for a given post ID.
func (r *PostRepository) GetCategoriesByPostID(postID int) ([]api.Category, error) {
	query := `
        SELECT 
            categories.id, 
            categories.name 
        FROM categories
        JOIN post_categories ON categories.id = post_categories.category_id
        WHERE post_categories.post_id = ?
    `
	rows, err := dbHandler.MainDB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []api.Category
	for rows.Next() {
		var category api.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// Create inserts a new post into the database.
func (r *PostRepository) Create(post *api.PostCreateRequest) error {
	stmt, err := dbHandler.MainDB.Prepare("INSERT INTO posts (user_id, title, content, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	post.CreatedAt = time.Now()

	result, err := stmt.Exec(post.UserID, post.Title, post.Content, post.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	post.ID = int(id)
	IncrementCount(int(post.UserID), "posts")
	return nil
}

// AddPostCategories inserts categories for a specific post into the database.
func (r *PostRepository) AddPostCategories(postID int, categories []string) error {
	for _, category := range categories {
		var categoryID int
		err := dbHandler.MainDB.QueryRow("SELECT id FROM categories WHERE name = ?", category).Scan(&categoryID)
		if err != nil {
			if err == sql.ErrNoRows {
				stmt, err := dbHandler.MainDB.Prepare("INSERT INTO categories (name) VALUES (?)")
				if err != nil {
					log.Fatal(err)
				}
				result, err := stmt.Exec(category)
				if err != nil {
					return err
				}
				id, err := result.LastInsertId()
				if err != nil {
					return err
				}
				categoryID = int(id)
			} else {
				return err
			}
		}

		stmt, err := dbHandler.MainDB.Prepare("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(postID, categoryID)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdatePostAmount updates the amount of comments for a specific post.
func UpdatePostAmount(postID int) error {
	log.Println("Updating post amount")
	stmt, err := dbHandler.MainDB.Prepare("UPDATE posts SET amount_of_comments = amount_of_comments + 1 WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(postID)
	return err
}
