package services

import (
	"EasySplit/internal/database"
	"EasySplit/internal/models"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func CreateFriendRequest(senderEmail, receiverEmail string) error {
	if senderEmail == receiverEmail {
		return fmt.Errorf("cannot send friend request to yourself")
	}

	isFriend, err := IsFriend(senderEmail, receiverEmail)
	if err != nil {
		return err
	}
	if isFriend {
		return fmt.Errorf("users are already friends")
	}

	exists, err := checkFriendRequestExists(senderEmail, receiverEmail)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("friend request already sent")
	}

	query := `INSERT INTO friend_requests (id, sender_email, receiver_email, status) 
			  VALUES ($1, $2, $3, $4)`
	_, err = database.DB.Exec(query, uuid.New(), senderEmail, receiverEmail, "pending")
	return err
}

func checkFriendRequestExists(senderEmail, receiverEmail string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM friend_requests WHERE sender_email = $1 AND receiver_email = $2 AND status = 'pending'`
	err := database.DB.QueryRow(query, senderEmail, receiverEmail).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func IsFriend(userEmail, friendEmail string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM friends WHERE (user_email = $1 AND friend_email = $2) OR (user_email = $2 AND friend_email = $1)`
	err := database.DB.QueryRow(query, userEmail, friendEmail).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func AcceptFriendRequest(senderEmail, receiverEmail string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE friend_requests SET status = 'accepted' WHERE sender_email = $1 AND receiver_email = $2`
	res, err := tx.Exec(query, senderEmail, receiverEmail)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("friend request not found")
	}

	query = `INSERT INTO friends (id, user_email, friend_email, created_at) 
	         VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)`
	_, err = tx.Exec(query,
		uuid.New(), senderEmail, receiverEmail, time.Now(),
		uuid.New(), receiverEmail, senderEmail, time.Now(),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetFriendByEmail(userEmail, friendEmail string) (models.Friend, error) {
	var friend models.Friend
	query := `
		SELECT f.id, f.user_email, f.friend_email, f.created_at, 
			   u.first_name, u.last_name, u.username
		FROM friends f
		JOIN users u ON u.email = f.friend_email
		WHERE f.user_email = $1 AND f.friend_email = $2
	`
	err := database.DB.QueryRow(query, userEmail, friendEmail).Scan(
		&friend.ID, &friend.UserEmail, &friend.FriendEmail, &friend.CreatedAt,
		&friend.FirstName, &friend.LastName, &friend.Username,
	)
	if err != nil {
		return models.Friend{}, fmt.Errorf("friend not found: %w", err)
	}
	return friend, nil
}


func GetFriendsByEmail(email string) ([]models.Friend, error) {
	var friends []models.Friend
	query := `
		SELECT f.id,  f.friend_email, f.created_at, 
			   u.first_name, u.last_name, u.username
		FROM friends f
		JOIN users u ON u.email = f.friend_email
		WHERE f.user_email = $1
	`
	rows, err := database.DB.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var friend models.Friend
		err := rows.Scan(
			&friend.ID, &friend.FriendEmail, &friend.CreatedAt,
			&friend.FirstName, &friend.LastName, &friend.Username,
		)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	return friends, rows.Err()
}

func RemoveFriend(userEmail, friendEmail string) error {
	query := `DELETE FROM friends WHERE (user_email = $1 AND friend_email = $2) OR (user_email = $2 AND friend_email = $1)`
	_, err := database.DB.Exec(query, userEmail, friendEmail)
	return err
}

func GetPendingFriendRequests(email string) ([]models.FriendRequest, error) {
	var requests []models.FriendRequest
	query := `
		SELECT fr.id, fr.sender_email, fr.receiver_email, fr.created_at, fr.status,
			   u.first_name, u.last_name, u.username
		FROM friend_requests fr
		JOIN users u ON u.email = fr.sender_email
		WHERE fr.receiver_email = $1 AND fr.status = 'pending'
	`
	rows, err := database.DB.Query(query, email)
	if err != nil {
		return nil, fmt.Errorf("error querying pending friend requests: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var request models.FriendRequest
		err := rows.Scan(
			&request.ID, &request.SenderEmail, &request.ReceiverEmail, &request.CreatedAt, &request.Status,
			&request.FirstName, &request.LastName, &request.Username,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning pending friend request: %w", err)
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over pending friend requests: %w", err)
	}

	return requests, nil
}
