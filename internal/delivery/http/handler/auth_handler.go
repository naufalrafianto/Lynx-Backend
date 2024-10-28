package handler

import (
	"auth-service/internal/delivery/http/response"
	"auth-service/internal/domain"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase domain.UserUsecase
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type otpRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

func NewAuthHandler(authUsecase domain.UserUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request", err))
		return
	}

	// Validate password strength
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, response.Error("Password must be at least 8 characters long", nil))
		return
	}

	user := &domain.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	if err := h.authUsecase.Register(user); err != nil {
		log.Printf("Registration error: %v", err)
		c.JSON(http.StatusInternalServerError, response.Error("Registration failed", err))
		return
	}

	message := "Registration successful. "
	if c.GetString("APP_ENV") == "development" {
		message += "Check the server logs for the OTP code."
	} else {
		message += "Please check your email for the OTP code."
	}

	c.JSON(http.StatusCreated, response.Success(message, nil))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request", err))
		return
	}

	token, err := h.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error("Login failed", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Login successful", gin.H{
		"token": token,
	}))
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req otpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request", err))
		return
	}

	if err := h.authUsecase.VerifyOTP(req.Email, req.OTP); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("OTP verification failed", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("Email verified successfully", nil))
}

func (h *AuthHandler) ResendOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request", err))
		return
	}

	if err := h.authUsecase.ResendOTP(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Failed to resend OTP", err))
		return
	}

	c.JSON(http.StatusOK, response.Success("OTP resent successfully", nil))
}
