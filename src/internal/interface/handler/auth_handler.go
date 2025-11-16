package handler

import (
	"fmt"
	"log"
	"net/http"
	"social-platform-backend/config"
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
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")

	conf := config.GetConfig()
	if token == "" {
		log.Printf("[Err] Missing token in AuthHandler.VerifyEmail")
		redirectURL := fmt.Sprintf("%s/auth/verify-result?success=false&message=Missing+verification+token", conf.Client.Url)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// Verify token
	if err := h.authService.VerifyEmail(token); err != nil {
		log.Printf("[Err] Error verifying email in AuthHandler.VerifyEmail: %v", err)

		if strings.Contains(err.Error(), "expired") {
			redirectURL := fmt.Sprintf("%s/auth/verify-result?success=false&message=Token+has+expired.+Please+request+a+new+verification+email", conf.Client.Url)
			c.Redirect(http.StatusFound, redirectURL)
			return
		}

		redirectURL := fmt.Sprintf("%s/auth/verify-result?success=false&message=Invalid+token.+Please+try+again", conf.Client.Url)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	redirectURL := fmt.Sprintf("%s/auth/verify-result?success=true&message=Email+verified+successfully", conf.Client.Url)
	c.Redirect(http.StatusFound, redirectURL)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in AuthHandler.Login: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Login user
	loginResponse, err := h.authService.Login(&req)
	if err != nil {
		log.Printf("[Err] Error in service layer AuthHandler.Login: %v", err)
		if strings.Contains(err.Error(), "not verified") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "Email not verified. Please verify your email first",
			})
			return
		}

		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Invalid email or password",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    loginResponse,
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req request.ForgotPasswordRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in AuthHandler.ForgotPassword: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Process forgot password
	if err := h.authService.ForgotPassword(&req); err != nil {
		log.Printf("[Err] Error in service layer AuthHandler.ForgotPassword: %v", err)

		if strings.Contains(err.Error(), "not verified") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "Email not verified. Please verify your email first",
			})
			return
		}

		c.JSON(http.StatusOK, response.APIResponse{
			Success: true,
			Message: "If your email is registered, you will receive a password reset link",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Password reset link has been sent to your email",
	})
}

func (h *AuthHandler) VerifyResetToken(c *gin.Context) {
	token := c.Query("token")
	conf := config.GetConfig()

	if token == "" {
		log.Printf("[Err] Missing token in AuthHandler.VerifyResetToken")
		redirectURL := fmt.Sprintf("%s/auth/reset-password?success=false&message=Missing+reset+token", conf.Client.Url)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// Verify token
	validToken, err := h.authService.VerifyResetToken(token)
	if err != nil {
		log.Printf("[Err] Error verifying reset token in AuthHandler.VerifyResetToken: %v", err)

		if strings.Contains(err.Error(), "expired") {
			redirectURL := fmt.Sprintf("%s/auth/reset-password?success=false&message=Token+has+expired", conf.Client.Url)
			c.Redirect(http.StatusFound, redirectURL)
			return
		}

		redirectURL := fmt.Sprintf("%s/auth/reset-password?success=false&message=Invalid+token", conf.Client.Url)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	redirectURL := fmt.Sprintf("%s/auth/reset-password?token=%s", conf.Client.Url, validToken)
	c.Redirect(http.StatusFound, redirectURL)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req request.ResetPasswordRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in AuthHandler.ResetPassword: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Reset password
	if err := h.authService.ResetPassword(&req); err != nil {
		log.Printf("[Err] Error in service layer AuthHandler.ResetPassword: %v", err)

		if strings.Contains(err.Error(), "expired") {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "Token has expired. Please request a new password reset",
			})
			return
		}

		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid or expired token",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Password has been reset successfully. You can now login with your new password",
	})
}

func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	var req request.ResendVerificationRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in AuthHandler.ResendVerificationEmail: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Resend verification email
	if err := h.authService.ResendVerificationEmail(&req); err != nil {
		log.Printf("[Err] Error in service layer AuthHandler.ResendVerificationEmail: %v", err)
		if strings.Contains(err.Error(), "already verified") {
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Message: "Email is already verified",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Email not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to resend verification email",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Verification email has been resent. Please check your email",
	})
}

func (h *AuthHandler) ResendResetPasswordEmail(c *gin.Context) {
	var req request.ResendVerificationRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in AuthHandler.ResendResetPasswordEmail: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Resend reset password email
	if err := h.authService.ResendResetPasswordEmail(&req); err != nil {
		log.Printf("[Err] Error in service layer AuthHandler.ResendResetPasswordEmail: %v", err)
		if strings.Contains(err.Error(), "not verified") {
			c.JSON(http.StatusForbidden, response.APIResponse{
				Success: false,
				Message: "Email not verified. Please verify your email first",
			})
			return
		}

		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, response.APIResponse{
				Success: false,
				Message: "Email not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to resend reset password email",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Reset password email has been resent. Please check your email",
	})
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req request.GoogleLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Err] Error binding JSON in AuthHandler.GoogleLogin: %v", err)
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	loginResponse, err := h.authService.GoogleLogin(&req)
	if err != nil {
		log.Printf("[Err] Error in service layer AuthHandler.GoogleLogin: %v", err)

		if strings.Contains(err.Error(), "invalid Google ID token") {
			c.JSON(http.StatusUnauthorized, response.APIResponse{
				Success: false,
				Message: "Invalid Google ID token. Please try again",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to login with Google",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Google login successful",
		Data:    loginResponse,
	})
}
