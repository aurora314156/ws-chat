package controllers_test

import (
	"errors"
	"fmt"
	"log"
	"testing"
	"ws-chat/controller"
	"ws-chat/db"
	"ws-chat/models"
	"ws-chat/tool"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/supabase-community/supabase-go"
)

var supaClient *supabase.Client
var testEmail string
var fullName string
var UUID uuid.UUID

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Warning: Could not load .env file. Assuming environment variables are already set.")
	}
	supaClient, err = db.NewSupabaseClient()
	if err != nil || supaClient == nil {
		log.Fatalf("Failed to initialize Supabase client: %v", err)
	}
	UUID = tool.GenUUID()
	testEmail = "ricktest111@gmail.com"
	fullName = "Rick Test"
}

func TestNewUserController(t *testing.T) {
	userController := controller.NewUserController(supaClient)
	assert.NotNil(t, userController)
}

// --- Test CreateUser ---

func TestUserController_CreateUser_Success(t *testing.T) {
	userController := controller.NewUserController(supaClient)

	testUser := &models.User{
		ID:             UUID,
		Email:          testEmail,
		HashedPassword: "hashed_create_pwd",
		FullName:       fullName,
	}

	err := userController.CreateUser(testUser)
	assert.NoError(t, err, "CreateUser should succeed without error")

	createdUser, err := userController.GetUser(UUID)
	assert.NoError(t, err, "GetUser should succeed without error after creation")
	assert.NotNil(t, createdUser, "Created user should not be nil")
	assert.Equal(t, testUser.Email, createdUser.Email, "Email should match")
	assert.Equal(t, testUser.FullName, createdUser.FullName, "FullName should match")
	assert.Equal(t, testUser.ID, createdUser.ID, "ID should match")
	assert.Equal(t, testUser.HashedPassword, createdUser.HashedPassword, "HashedPassword should match")

	err = userController.DeleteUser(UUID)
	if err != nil && !errors.Is(err, controller.ErrUserNotFound) {
		t.Fatalf("Failed to clean up before test. DeleteUser returned unexpected error: %v", err)
	}

}

func TestUserController_UpdateUser_Success(t *testing.T) {
	userController := controller.NewUserController(supaClient)

	testUser := &models.User{
		ID:             UUID,
		Email:          testEmail,
		HashedPassword: "hashed_create_pwd",
		FullName:       fullName,
	}

	err := userController.CreateUser(testUser)
	assert.NoError(t, err, "CreateUser should succeed without error")
	// update page will get the user info first
	createdUser, err := userController.GetUser(UUID)
	// mock user do update
	userUpdate := createdUser
	userUpdate.Email = "rick new email@gmail.com"
	userUpdate.FullName = "rick new name"
	userUpdate.HashedPassword = "new hashed pwd"

	err = userController.UpdateUser(UUID, userUpdate)
	assert.NoError(t, err, "Updated should succeed without error")

	updatedUser, err := userController.GetUser(UUID)
	assert.NoError(t, err, "GetUser should succeed without error after update")
	assert.NotNil(t, updatedUser, "Updated user should not be nil")
	assert.Equal(t, userUpdate.Email, updatedUser.Email, "Email should match")
	assert.Equal(t, userUpdate.FullName, updatedUser.FullName, "FullName should match")
	assert.Equal(t, userUpdate.HashedPassword, updatedUser.HashedPassword, "HashedPassword should match")

	err = userController.DeleteUser(UUID)
	if err != nil && !errors.Is(err, controller.ErrUserNotFound) {
		t.Fatalf("Failed to clean up before test. DeleteUser returned unexpected error: %v", err)
	}

}
