package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kenkinoti/gofiber-ago-crm-backend/internal/models"
)

func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated",
			},
		})
		return
	}

	var user models.User
	if err := h.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "USER_NOT_FOUND",
				"message": "User not found",
			},
		})
		return
	}

	// Convert to response format (without password hash)
	userResponse := UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Phone:          user.Phone,
		Role:           user.Role,
		OrganizationID: user.OrganizationID,
		IsActive:       user.IsActive,
		LastLoginAt:    user.LastLoginAt,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userResponse,
	})
}

func (h *Handler) GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get users endpoint - TODO: implement",
	})
}

func (h *Handler) CreateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Create user endpoint - TODO: implement",
	})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Update user endpoint - TODO: implement",
	})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete user endpoint - TODO: implement",
	})
}
