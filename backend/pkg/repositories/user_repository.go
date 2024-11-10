package repositories

import (
	"database/sql"
	"project-root/pkg/api"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{DB: dbHandler.MainDB}
}

func (u *UserRepository) GetAllUsers(page, pageSize int, sortBy string, excludeID int) ([]api.Users, int, int, error) {
	// Define the base query
	query := `
		SELECT 
			users.id, 
			users.nickname 
		FROM users
	`

	// Add condition to exclude a specific user ID if provided
	if excludeID != 0 {
		query += " WHERE users.id != ?"
	}

	// Add sorting if provided
	if sortBy != "" {
		query += " ORDER BY " + sortBy
	} else {
		query += " ORDER BY LOWER(users.nickname) ASC"
	}

	// Calculate pagination offsets
	offset := (page - 1) * pageSize
	query += " LIMIT ? OFFSET ?"

	var rows *sql.Rows
	var err error

	// Execute the query with or without excludeID condition
	if excludeID != 0 {
		rows, err = u.DB.Query(query, excludeID, pageSize, offset)
	} else {
		rows, err = u.DB.Query(query, pageSize, offset)
	}

	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	// Parse the result set
	users := []api.Users{}
	for rows.Next() {
		var user api.Users
		if err := rows.Scan(&user.ID, &user.Nickname); err != nil {
			return nil, 0, 0, err
		}
		users = append(users, user)
	}

	// Count total number of users with or without excludeID condition
	countQuery := "SELECT COUNT(*) FROM users"
	if excludeID != 0 {
		countQuery += " WHERE users.id != ?"
	}

	var totalUsers int
	if excludeID != 0 {
		err = u.DB.QueryRow(countQuery, excludeID).Scan(&totalUsers)
	} else {
		err = u.DB.QueryRow(countQuery).Scan(&totalUsers)
	}
	if err != nil {
		return nil, 0, 0, err
	}

	// Calculate the total number of pages
	totalPages := totalUsers / pageSize
	if totalUsers%pageSize != 0 {
		totalPages++
	}

	return users, totalUsers, totalPages, nil
}


// CreateUser creates a new user in the database and returns a UserResponse
func CreateUser(u *api.RegistrationRequest) (*api.UserResponse, error) {
	stmt, err := dbHandler.MainDB.Prepare("INSERT INTO users (nickname, age, gender, first_name, last_name, email, password_hash, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	

	_, err = stmt.Exec(u.Nickname, u.Age, u.Gender, u.FirstName, u.LastName, u.Email, u.Password, u.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Fetch the newly created user from the database
	var newUser api.UserResponse
	err = dbHandler.MainDB.QueryRow("SELECT id, nickname FROM users WHERE email = ?", u.Email).Scan(&newUser.ID, &newUser.Nickname)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// GetUserByEmail retrieves a user by email from the database
// ChekUserByEmail retrieves a user by email from the database
func ChekUserByEmail(email string) (*api.LoginRequest, *api.UserResponse, error) {
	var user api.LoginRequest
	var userD api.UserResponse
	err := dbHandler.MainDB.QueryRow("SELECT email, password_hash, id, nickname FROM users WHERE email = ?", email).Scan(
		&user.Email, &user.Password, &userD.ID, &userD.Nickname)
	if err == sql.ErrNoRows {
		return nil, nil, err
	} else if err != nil {
		return nil, nil, err
	}
	return &user, &userD, nil
}

// GetUserByNickname retrieves a user by nickname from the database
func GetUserByNickname(nickname string) (*api.User, error) {
	var user api.User
	err := dbHandler.MainDB.QueryRow("SELECT id, nickname, first_name, last_name, created_at, amount_of_posts, amount_of_comments FROM users WHERE nickname = ?", nickname).Scan(
		&user.ID, &user.Nickname, &user.FirstName, &user.LastName, &user.CreatedAt, &user.AmountOfPosts, &user.AmountOfComments)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserBySessionID(sessionID string) (*api.User, error) {
	var user api.User
	err := dbHandler.MainDB.QueryRow(`
		SELECT u.id, u.nickname 
		FROM users u
		INNER JOIN active_sessions s ON u.id = s.user_id
		WHERE s.session_id = ?
	`, sessionID).Scan(&user.ID, &user.Nickname)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found with the session ID
		}
		return nil, err // Other error occurred
	}

	return &user, nil
}

// // UpdateUser updates user information in the database
// func UpdateUser(u *api.UpdateUserRequest) error {
// 	stmt, err := dbHandler.MainDB.Prepare("UPDATE users SET nickname=?, age=?, gender=?, first_name=?, last_name=?, email=? WHERE id=?")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(u.Nickname, u.Age, u.Gender, u.FirstName, u.LastName, u.Email, u.ID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
// // IncrementPostCount increments the post count for a user
// func IncrementPostCount(userID int) error {
// 	stmt, err := dbHandler.MainDB.Prepare("UPDATE users SET amount_of_posts = amount_of_posts + 1 WHERE id = ?")
// 	if err != nil {
// 		return err
// 	}
// 	_, err = stmt.Exec(userID)
// 	return err
// }

// IncrementCount increments a specific count (e.g., amount_of_comments) for a user.
func IncrementCount(userID int, what string) error {
	// Prepare the SQL statement with a placeholder for the count field
	what = "amount_of_" + what
	stmt, err := dbHandler.MainDB.Prepare("UPDATE users SET " + what + " = " + what + " + 1 WHERE id = ?")
	if err != nil {
		return err
	}

	// Execute the prepared statement with the userID parameter
	_, err = stmt.Exec(userID)
	return err
}

func GetUserByID(userID int) (api.UserResponse, error) {
	var user api.UserResponse
	err := dbHandler.MainDB.QueryRow("SELECT id, nickname FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Nickname)
	if err != nil {
		return api.UserResponse{}, err
	}
	return user, nil
}