package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/kenkinoti/gofiber-ago-crm-backend/internal/config"
	"github.com/kenkinoti/gofiber-ago-crm-backend/internal/middleware"
)

type Handler struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewHandler(db *gorm.DB, cfg *config.Config) *Handler {
	return &Handler{
		DB:     db,
		Config: cfg,
	}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", h.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.RefreshToken)
		}

		// Protected routes (require authentication)
		protected := v1.Group("/")
		protected.Use(middleware.AuthRequired(h.Config))
		{
			// Auth routes that require authentication
			auth := protected.Group("/auth")
			{
				auth.POST("/logout", h.Logout)
			}

			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", h.GetCurrentUser)
				users.GET("", h.GetUsers)
				users.POST("", middleware.RequireRole("admin"), h.CreateUser)
				users.PUT("/:id", h.UpdateUser)
				users.DELETE("/:id", middleware.RequireRole("admin"), h.DeleteUser)
			}

			// Participant routes
			participants := protected.Group("/participants")
			{
				participants.GET("", h.GetParticipants)
				participants.GET("/:id", h.GetParticipant)
				participants.POST("", h.CreateParticipant)
				participants.PUT("/:id", h.UpdateParticipant)
				participants.DELETE("/:id", h.DeleteParticipant)
			}

			// Shift routes
			shifts := protected.Group("/shifts")
			{
				shifts.GET("", h.GetShifts)
				shifts.GET("/:id", h.GetShift)
				shifts.POST("", h.CreateShift)
				shifts.PUT("/:id", h.UpdateShift)
				shifts.PATCH("/:id/status", h.UpdateShiftStatus)
				shifts.DELETE("/:id", h.DeleteShift)
			}

			// Document routes
			documents := protected.Group("/documents")
			{
				documents.GET("", h.GetDocuments)
				documents.GET("/:id", h.GetDocument)
				documents.POST("", h.UploadDocument)
				documents.PUT("/:id", h.UpdateDocument)
				documents.DELETE("/:id", h.DeleteDocument)
				documents.GET("/:id/download", h.DownloadDocument)
			}

			// Organization routes
			organization := protected.Group("/organization")
			{
				organization.GET("", h.GetOrganization)
				organization.PUT("", middleware.RequireRole("admin"), h.UpdateOrganization)
			}
		}
	}
}
