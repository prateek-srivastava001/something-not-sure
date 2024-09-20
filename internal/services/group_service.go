package services

import (
	"EasySplit/internal/database"
	"EasySplit/internal/models"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func CreateGroup(creatorEmail, name string) (uuid.UUID, error) {
	var count int
	query := `SELECT COUNT(*) FROM groups WHERE name = $1 AND creator_email = $2`
	err := database.DB.QueryRow(query, name, creatorEmail).Scan(&count)
	if err != nil {
		return uuid.Nil, err
	}
	if count > 0 {
		return uuid.Nil, fmt.Errorf("group with this name already exists")
	}

	groupID := uuid.New()
	query = `INSERT INTO groups (id, name, creator_email, created_at) VALUES ($1, $2, $3, $4)`
	_, err = database.DB.Exec(query, groupID, name, creatorEmail, time.Now())
	if err != nil {
		return uuid.Nil, err
	}
	return groupID, nil
}

func AddUserToGroup(groupID uuid.UUID, userEmail, creatorEmail string) error {
	// Check if the user is already in the group
	var count int
	query := `SELECT COUNT(*) FROM group_members WHERE group_id = $1 AND member_email = $2`
	err := database.DB.QueryRow(query, groupID, userEmail).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("user already in group")
	}

	// Check if the group exists
	var groupCount int
	query = `SELECT COUNT(*) FROM groups WHERE id = $1`
	err = database.DB.QueryRow(query, groupID).Scan(&groupCount)
	if err != nil {
		return err
	}
	if groupCount == 0 {
		return fmt.Errorf("group does not exist")
	}

	// Check if the user is a valid user
	var userCount int
	query = `SELECT COUNT(*) FROM users WHERE email = $1`
	err = database.DB.QueryRow(query, userEmail).Scan(&userCount)
	if err != nil {
		return err
	}
	if userCount == 0 {
		return fmt.Errorf("user does not exist")
	}

	// Insert user into the group
	query = `INSERT INTO group_members (id, group_id, member_email, joined_at) VALUES ($1, $2, $3, $4)`
	_, err = database.DB.Exec(query, uuid.New(), groupID, userEmail, time.Now())
	return err
}

func GetGroupByID(groupID uuid.UUID) (*models.Group, error) {
	var group models.Group
	query := `SELECT id, name, creator_email, created_at FROM groups WHERE id = $1`
	err := database.DB.QueryRow(query, groupID).Scan(&group.ID, &group.Name, &group.CreatorEmail, &group.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}
	return &group, nil
}

func GetAllGroups(userEmail string) ([]models.Group, error) {
	var groups []models.Group
	query := `
		SELECT DISTINCT g.id, g.name, g.creator_email, g.created_at
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.member_email = $1
	`
	rows, err := database.DB.Query(query, userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.CreatorEmail, &group.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, rows.Err()
}

func DeleteGroup(groupID uuid.UUID, creatorEmail string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var groupCreator string
	query := `SELECT creator_email FROM groups WHERE id = $1`
	err = tx.QueryRow(query, groupID).Scan(&groupCreator)
	if err != nil {
		return err
	}
	if groupCreator != creatorEmail {
		return fmt.Errorf("user is not the creator of the group")
	}

	_, err = tx.Exec(`DELETE FROM group_members WHERE group_id = $1`, groupID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM groups WHERE id = $1`, groupID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
