package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/kenkinoti/gofiber-ago-crm-backend/internal/models"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         UserResponse `json:"user"`
	ExpiresIn    int         `json:"expires_in"`
}

type UserResponse struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Phone          string    `json:"phone"`
	Role           string    `json:"role"`
	OrganizationID string    `json:"organization_id"`
	IsActive       bool      `json:"is_active"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request parameters",
				"details": err.Error(),
			},
		})
		return
	}

	// Find user by email
	var user models.User
	if err := h.DB.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password",
			},
		})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"org_id":  user.OrganizationID,
		"exp":     time.Now().Add(h.Config.JWTExpiry).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.Config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "TOKEN_GENERATION_ERROR",
				"message": "Failed to generate token",
			},
		})
		return
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     time.Now().Add(h.Config.RefreshTokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(h.Config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "TOKEN_GENERATION_ERROR",
				"message": "Failed to generate refresh token",
			},
		})
		return
	}

	// Save refresh token to database
	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(h.Config.RefreshTokenExpiry),
		IsRevoked: false,
	}
	h.DB.Create(&refreshTokenModel)

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	h.DB.Save(&user)

	// Convert to response format
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
		"data": LoginResponse{
			Token:        tokenString,
			RefreshToken: refreshTokenString,
			User:         userResponse,
			ExpiresIn:    int(h.Config.JWTExpiry.Seconds()),
		},
		"message": "Login successful",
	})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request parameters",
			},
		})
		return
	}

	// Parse and validate refresh token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "Invalid refresh token",
			},
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "Invalid token claims",
			},
		})
		return
	}

	userID := claims["user_id"].(string)

	// Check if refresh token exists and is not revoked
	var refreshTokenModel models.RefreshToken
	if err := h.DB.Where("user_id = ? AND token = ? AND is_revoked = ? AND expires_at > ?", 
		userID, req.RefreshToken, false, time.Now()).First(&refreshTokenModel).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_TOKEN",
				"message": "Refresh token not found or expired",
			},
		})
		return
	}

	// Get user details
	var user models.User
	if err := h.DB.Where("id = ? AND is_active = ?", userID, true).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "USER_NOT_FOUND",
				"message": "User not found",
			},
		})
		return
	}

	// Generate new access token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"org_id":  user.OrganizationID,
		"exp":     time.Now().Add(h.Config.JWTExpiry).Unix(),
		"iat":     time.Now().Unix(),
	})

	newTokenString, err := newToken.SignedString([]byte(h.Config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "TOKEN_GENERATION_ERROR",
				"message": "Failed to generate new token",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token":      newTokenString,
			"expires_in": int(h.Config.JWTExpiry.Seconds()),
		},
		"message": "Token refreshed successfully",
	})
}

func (h *Handler) Logout(c *gin.Context) {
	// Get user from context (set by auth middleware)
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

	// Revoke all refresh tokens for this user
	h.DB.Model(&models.RefreshToken{}).Where("user_id = ?", userID).Update("is_revoked", true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}
