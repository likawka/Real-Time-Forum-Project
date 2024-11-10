package repositories

import (
	"log"
	"project-root/pkg/api"
	"time"
)

// CreateSession creates a new session in the database
func CreateSession(s *api.Session) error {
	stmt, err := dbHandler.MainDB.Prepare("INSERT INTO active_sessions (user_id, session_id, created_at, expires_at, last_activity) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(s.UserID, s.SessionID, s.CreatedAt, s.ExpiresAt, s.LastActivity)
	if err != nil {
		return err
	}
	return nil
}

// UpdateLastActivity updates the last activity timestamp for a session
func UpdateLastActivity(sessionID string) error {
	stmt, err := dbHandler.MainDB.Prepare("UPDATE active_sessions SET last_activity = ? WHERE session_id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(time.Now(), sessionID)
	return err
}

// DeleteSession deletes a session from the database
func DeleteSession(sessionID string) error {
	stmt, err := dbHandler.MainDB.Prepare("DELETE FROM active_sessions WHERE session_id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sessionID)
	if err != nil {
		return err
	}
	return nil
}
