package handler

import (
	"auth-service/internal/delivery/http/response"
	"auth-service/internal/domain"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(userUsecase domain.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, response.Error("Unauthorized", domain.ErrUnauthorized))
		return
	}

	user, err := h.userUsecase.GetUserByID(userID)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, response.Error("User not found", err))
		case domain.ErrCacheUnavailable:
			// Maybe try without cache
			c.JSON(http.StatusInternalServerError, response.Error("Service temporarily unavailable", err))
		default:
			c.JSON(http.StatusInternalServerError, response.Error("Failed to get user data", err))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("User data retrieved successfully", user))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Get user ID from JWT token
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, response.Error("Unauthorized", domain.ErrUnauthorized))
		return
	}

	// Check if this is a permanent deletion
	permanent := c.DefaultQuery("permanent", "false") == "true"

	// Add additional security for permanent deletion
	if permanent {
		// Only allow permanent deletion with specific confirmation
		confirmation := c.GetHeader("X-Confirm-Delete")
		if confirmation != "permanent-delete" {
			c.JSON(http.StatusBadRequest, response.Error(
				"Permanent deletion requires confirmation",
				errors.New("add 'X-Confirm-Delete: permanent-delete' header for permanent deletion"),
			))
			return
		}
	}

	// Delete the user
	err := h.userUsecase.DeleteUser(userID, permanent)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, response.Error("User not found", err))
		case domain.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, response.Error("Unauthorized", err))
		default:
			c.JSON(http.StatusInternalServerError, response.Error("Failed to delete user", err))
		}
		return
	}

	message := "User account deactivated successfully"
	if permanent {
		message = "User account permanently deleted"
	}

	c.JSON(http.StatusOK, response.Success(message, nil))
}
