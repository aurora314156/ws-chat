package controller

import (
	"encoding/json"
	"strings"
	"ws-chat/logger"
	"ws-chat/models"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

type UserController struct {
	Client *supabase.Client
}

func NewUserController(client *supabase.Client) *UserController {
	return &UserController{Client: client}
}

func (c *UserController) CreateUser(newUser models.User) error {
	var results []models.User

	resp, _, err := c.Client.From("users").Insert(newUser, false, "", "representation", "").Execute()
	if err != nil {
		logger.Error("error inserting user: %w", err)
		return err
	}

	err = json.Unmarshal(resp, &results)
	if err != nil {
		logger.Error("error unmarshaling response: %v", err)
		return err
	}

	logger.Info("Successfully created user: %+v", results)

	return nil
}

func (c *UserController) GetUser(userID uuid.UUID) (*models.User, error) {
	var results []models.User

	resp, _, err := c.Client.From("users").Select("*", "exact", false).Filter("id", "eq", userID.String()).Execute()
	if err != nil {
		logger.Error("error querying user: %v", err)
		return nil, err
	}

	err = json.Unmarshal(resp, &results)
	if err != nil {
		logger.Error("error unmarshaling response: %w", err)
		return nil, err
	}

	if len(results) == 0 {
		logger.Debug("user with ID %s not found", userID)
		return nil, err
	}
	logger.Info("Successfully Get user: %+v", results)

	return &results[0], nil
}

func (c *UserController) UpdateUser(userID uuid.UUID, updates map[string]any) error {
	var results []map[string]any

	resp, count, err := c.Client.From("users").Update(updates, "representation", "exact").Filter("id", "eq", userID.String()).Execute()
	if err != nil {
		logger.Error("error updating user: %w", err)
		return err
	}

	if count == 0 {
		logger.Error("user with ID %s not found or no changes made", userID)
		return err
	}

	err = json.Unmarshal(resp, &results)
	if err != nil {
		logger.Error("error unmarshaling response: %w", err)
		return err
	}

	logger.Info("Successfully updated user: %+v", results)

	return nil
}

func (c *UserController) DeleteUser(userID uuid.UUID) error {
	_, count, err := c.Client.From("users").Delete("minimal", "exact").Filter("id", "eq", userID.String()).Execute()
	if err != nil {
		logger.Error("error deleting user: %w", err)
		return err
	}

	if count == 0 {
		logger.Error("user with ID %s not found or already deleted", userID)
		return err
	}

	logger.Info("Successfully deleted %d user(s) with ID: %s", count, userID)

	return nil
}

func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errString := err.Error()

	// Check for common substrings in duplicate key errors for PostgreSQL and other databases
	isDuplicate := strings.Contains(errString, "duplicate key value") || strings.Contains(errString, "23505")

	return isDuplicate
}
