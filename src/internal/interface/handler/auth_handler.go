package handler

import (
	"log"
	"net/http"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in AuthHandler.Register: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Register user
	if err := h.authService.Register(&req); err != nil {
		log.Printf("[Err] Error in service layer AuthHandler.Register: %v", err)

		// Check for specific error types
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, response.APIResponse{
				Success: false,
				Message: "Email already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to register user",
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Registration successful. Please check your email to verify your account.",
		Data: response.RegisterResponse{
			Message: "A verification email has been sent to your email address",
			Email:   req.Email,
		},
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		log.Printf("[Err] Missing token in AuthHandler.VerifyEmail")
		c.Redirect(http.StatusFound, "http://localhost:3000/verify-result?success=false&message=Missing+verification+token")
		return
	}

	// Verify token
	if err := h.authService.VerifyEmail(token); err != nil {
		log.Printf("[Err] Error verifying email in AuthHandler.VerifyEmail: %v", err)

		if strings.Contains(err.Error(), "expired") {
			c.Redirect(http.StatusFound, "http://localhost:3000/verify-result?success=false&message=Token+has+expired.+Please+request+a+new+verification+email")
			return
		}

		c.Redirect(http.StatusFound, "http://localhost:3000/verify-result?success=false&message=Invalid+token.+Please+try+again")
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:3000/verify-result?success=true&message=Email+verified+successfully")
}
