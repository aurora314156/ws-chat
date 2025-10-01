package handler

import (
	"net/http"
	"ws-chat/controller"
	"ws-chat/logger"
	"ws-chat/models"
	"ws-chat/tool"

	"github.com/gin-gonic/gin"
)

// define request payload for signup
type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`    // email format validation
	Password string `json:"password" binding:"required,min=8"` // minimum length validation
	FullName string `json:"full_name"`                         //optional
}

type UserHandler struct {
	UserController *controller.UserController
}

func NewUserHandler(uc *controller.UserController) *UserHandler {
	return &UserHandler{
		UserController: uc,
	}
}

func (h *UserHandler) Signup(c *gin.Context) {
	var req SignupRequest

	// 1. binding and verification payload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request or validation failed", "details": err.Error()})
		return
	}

	// 2. password hashing
	hashedPassword, err := tool.HashedPassword([]byte(req.Password))
	if err != nil {
		logger.Error("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error during password hashing"})
		return
	}

	// 3. convert payload to User model
	newUser := &models.User{
		Email:          req.Email,
		HashedPassword: string(hashedPassword),
		FullName:       req.FullName,
		ID:             tool.GenUUID(),
	}

	// 4. call controller to create user
	if err := h.UserController.CreateUser(newUser); err != nil {
		// if it's a duplicate email error
		if err == controller.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists."})
			return
		}

		// other errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user due to internal error."})
		return
	}

	// 5. respond success
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
