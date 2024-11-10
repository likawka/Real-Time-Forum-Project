package repositories

import (
	"database/sql"
	"project-root/pkg/api"
	"time"
)

type RateableRepository interface {
	UpdateRateM(delta, postID, commentID int) (int, error)
}

func (p *PostRepository) UpdateRateM(delta, postID, commentID int) (int, error) {
	_, err := dbHandler.MainDB.Exec("UPDATE posts SET rate = rate + ? WHERE id = ?", delta, postID)
	if err != nil {
		return 0, err
	}

	var newRate int
	err = dbHandler.MainDB.QueryRow("SELECT rate FROM posts WHERE id = ?", postID).Scan(&newRate)
	if err != nil {
		return 0, err
	}

	return newRate, nil
}

func (c *CommentRepository) UpdateRateM(delta, postID, commentID int) (int, error) {
	_, err := dbHandler.MainDB.Exec("UPDATE comments SET rate = rate + ? WHERE id = ? AND post_id = ?", delta, commentID, postID)
	if err != nil {
		return 0, err
	}

	var newRate int
	err = dbHandler.MainDB.QueryRow("SELECT rate FROM comments WHERE id = ? AND post_id = ?", commentID, postID).Scan(&newRate)
	if err != nil {
		return 0, err
	}

	return newRate, nil
}

func UpdateRate(userID int, req api.RateRequest) (int, string, error) {
	var rateable RateableRepository
	var itemType string
	var searched_id int

	if req.CommentID != 0 {
		rateable = &CommentRepository{}
		itemType = "comment_id"
		searched_id = req.CommentID
	} else {
		rateable = &PostRepository{}
		itemType = "post_id"
		searched_id = req.PostID
	}

	var existingStatus string
	query := "SELECT status FROM rates WHERE user_id = ? AND " + itemType + " = ?"
	err := dbHandler.MainDB.QueryRow(query, userID, searched_id).Scan(&existingStatus)
	if err != nil && err != sql.ErrNoRows {
		return 0, "", err
	}

	delta := 0
	newStatus := req.Status
	var newRate int
	switch {
	case existingStatus == req.Status:
		deleteQuery := "DELETE FROM rates WHERE user_id = ? AND " + itemType + " = ?"
		_, err = dbHandler.MainDB.Exec(deleteQuery, userID, searched_id)
		if err != nil {
			return 0, "", err
		}
		if req.Status == "up" {
			delta = -1
		} else {
			delta = 1
		}
		newStatus = ""
	case existingStatus == "":
		insertQuery := "INSERT INTO rates (user_id, " + itemType + ", status, rated_at) VALUES (?, ?, ?, ?)"
		_, err = dbHandler.MainDB.Exec(insertQuery, userID, searched_id, req.Status, time.Now())
		if err != nil {
			return 0, "", err
		}
		if req.Status == "up" {
			delta = 1
		} else {
			delta = -1
		}
	default:
		updateQuery := "UPDATE rates SET status = ?, rated_at = ? WHERE user_id = ? AND " + itemType + " = ?"
		_, err = dbHandler.MainDB.Exec(updateQuery, req.Status, time.Now(), userID, searched_id)
		if err != nil {
			return 0, "", err
		}
		if req.Status == "up" && existingStatus == "down" {
			delta = 2
		} else if req.Status == "down" && existingStatus == "up" {
			delta = -2
		}
	}
	var updateErr error
	if req.CommentID != 0 {
		newRate, updateErr = rateable.UpdateRateM(delta, req.PostID, req.CommentID)
	} else {
		newRate, updateErr = rateable.UpdateRateM(delta, req.PostID, 0)
	}

	if updateErr != nil {
		return 0, "", updateErr
	}

	return newRate, newStatus, nil
}

func getRateStatus(itemType string, itemID, userID int) (string, error) {
	var status sql.NullString
	var query string

	switch itemType {
	case "post":
		query = "SELECT status FROM rates WHERE user_id = ? AND post_id = ?"
	case "comment":
		query = "SELECT status FROM rates WHERE user_id = ? AND comment_id = ?"
	default:
		return "", nil
	}

	err := dbHandler.MainDB.QueryRow(query, userID, itemID).Scan(&status)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if status.Valid {
		return status.String, nil
	}
	return "", nil
}
