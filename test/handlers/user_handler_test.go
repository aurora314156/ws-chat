package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"ws-chat/controller"
	"ws-chat/db"
	"ws-chat/handler"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/supabase-community/supabase-go"
	"golang.org/x/crypto/bcrypt"
)

var supaClient *supabase.Client
var userHandler *handler.UserHandler
var router *gin.Engine

var testEmail string
var testPassword string
var testFullName string

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Warning: Could not load .env file. Assuming environment variables are already set.")
	}
	supaClient, err = db.NewSupabaseClient()
	if err != nil || supaClient == nil {
		log.Fatalf("Failed to initialize Supabase client: %v", err)
	}

	userHandler = handler.NewUserHandler(controller.NewUserController(supaClient))
	gin.SetMode(gin.TestMode)
	router = gin.New()
	router.POST("/signup", userHandler.Signup)

	testEmail = "ricktest111@gmail.com"
	testPassword = "testPassword123"
	testFullName = "Rick Test"
}

func TestUserHander_SignUp(t *testing.T) {
	// 1. prepare request payload
	signupReq := handler.SignupRequest{
		Email:    testEmail,
		Password: testPassword,
		FullName: testFullName,
	}
	body, _ := json.Marshal(signupReq)

	// 2. create HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 3. perform the request
	router.ServeHTTP(w, req)

	// 4. assert response
	assert.Equal(t, http.StatusCreated, w.Code, "Expected HTTP 201 Created on success")

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "User registered successfully", response["message"])

	userIDStr, ok := response["user_id"].(string)
	assert.True(t, ok, "Response should contain 'user_id' as a string")

	// 5. parse user_id to UUID
	createdUUID, err := uuid.Parse(userIDStr)
	assert.NoError(t, err, "The returned user_id should be a valid UUID string")

	// 6. verify user in DB
	userController := controller.NewUserController(supaClient)
	createdUser, err := userController.GetUser(createdUUID)

	assert.NoError(t, err, "GetUser should not return an error")
	assert.NotNil(t, createdUser, "User should be found in the database with the fixed ID")

	// 7. verify fields
	assert.Equal(t, createdUUID, createdUser.ID, "User ID in DB should match the fixed ID")
	assert.Equal(t, testEmail, createdUser.Email, "Email in DB should match the request email")
	assert.Equal(t, testFullName, createdUser.FullName, "FullName in DB should match the request FullName")

	// 8. verify password is hashed
	err = bcrypt.CompareHashAndPassword([]byte(createdUser.HashedPassword), []byte(testPassword))
	assert.NoError(t, err, "Stored HashedPassword should match the testPassword")

}
