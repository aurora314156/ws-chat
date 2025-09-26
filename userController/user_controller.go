package usercontroller

import (
	"encoding/json"
	"ws-chat/logger"
	"ws-chat/models"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

func CreateUser(client *supabase.Client, newUser models.User) error {
	var results []models.User

	resp, _, err := client.From("users").Insert(newUser, false, "", "representation", "").Execute()
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

func GetUser(client *supabase.Client, userID uuid.UUID) (*models.User, error) {
	var results []models.User

	resp, _, err := client.From("users").Select("*", "exact", false).Filter("id", "eq", userID.String()).Execute()
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

func UpdateUser(client *supabase.Client, userID uuid.UUID, updates map[string]interface{}) error {
	var results []map[string]interface{}

	resp, count, err := client.From("users").Update(updates, "representation", "exact").Filter("id", "eq", userID.String()).Execute()
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

func DeleteUser(client *supabase.Client, userID uuid.UUID) error {
	_, count, err := client.From("users").Delete("minimal", "exact").Filter("id", "eq", userID.String()).Execute()
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
