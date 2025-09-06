package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/kenkinoti/gofiber-ago-crm-backend/internal/config"
	"github.com/kenkinoti/gofiber-ago-crm-backend/internal/middleware"
	"gorm.io/gorm"
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

// Helper methods for handlers
func (h *Handler) GetUserIDFromContext(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

func (h *Handler) GetUserRoleFromContext(c *gin.Context) string {
	if role, exists := c.Get("user_role"); exists {
		if r, ok := role.(string); ok {
			return r
		}
	}
	return ""
}

func (h *Handler) CanUserAccessResource(c *gin.Context, permission string, resourceUserID string) bool {
	currentUserID := h.GetUserIDFromContext(c)
	currentUserRole := h.GetUserRoleFromContext(c)
	
	// User can access their own resources
	if currentUserID == resourceUserID {
		return true
	}
	
	// Admins and managers can access other users' resources
	switch currentUserRole {
	case "super_admin", "admin", "manager":
		return true
	default:
		return false
	}
}

func (h *Handler) SendErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	response := gin.H{
		"success": false,
		"error": gin.H{
			"message": message,
		},
	}
	
	if err != nil {
		response["error"].(gin.H)["details"] = err.Error()
	}
	
	c.JSON(statusCode, response)
}

func (h *Handler) SendSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
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
			// Handle OPTIONS requests explicitly
			auth.OPTIONS("/login", func(c *gin.Context) {
				c.Status(200)
			})
			auth.OPTIONS("/refresh", func(c *gin.Context) {
				c.Status(200)
			})

			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.RefreshToken)
			auth.GET("/test-accounts", h.GetTestAccounts)
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
				users.GET("/timezones", h.GetSupportedTimezones)
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
				organization.PUT("", middleware.RequireRole("admin", "manager"), h.UpdateOrganization)

				// Organization branding routes
				organization.GET("/branding", h.GetOrganizationBranding)
				organization.PUT("/branding", middleware.RequireRole("admin"), h.UpdateOrganizationBranding)

				// Organization settings routes
				organization.GET("/settings", h.GetOrganizationSettings)
				organization.PUT("/settings", middleware.RequireRole("admin"), h.UpdateOrganizationSettings)

				// Organization subscription routes
				organization.GET("/subscription", middleware.RequireRole("admin"), h.GetOrganizationSubscription)
				organization.PUT("/subscription", middleware.RequireRole("admin"), h.UpdateOrganizationSubscription)
			}

			// Super Admin routes (require super_admin role)
			superAdmin := protected.Group("/super-admin")
			superAdmin.Use(middleware.RequireSuperAdmin())
			{
				// Organization management
				organizations := superAdmin.Group("/organizations")
				{
					organizations.GET("", h.GetAllOrganizations)
					organizations.GET("/:id", h.GetOrganizationById)
					organizations.POST("", h.CreateOrganization)
					organizations.PATCH("/:id/status", h.UpdateOrganizationStatus)
					organizations.DELETE("/:id", h.DeleteOrganization)
				}
			}

			// Admin routes (require admin role)
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin", "super_admin"))
			{
				admin.POST("/seed", h.SeedDatabase)
				admin.POST("/seed-organizations", h.SeedOrganizations)
				admin.POST("/seed-advanced", h.SeedAdvanced)
				admin.DELETE("/clear-test-data", middleware.RequireElevatedAuth(), h.ClearTestData)
				admin.DELETE("/truncate", middleware.RequireElevatedAuth(), middleware.RequirePasswordConfirmation(), h.TruncateDatabase)

				// Database management routes (admin can view stats)
				admin.GET("/stats", h.GetSystemStats)
				admin.GET("/tables", h.GetTableStats)
				admin.POST("/backup", middleware.RequireSuperAdmin(), middleware.RequireElevatedAuth(), h.DatabaseBackup)
				admin.POST("/maintenance", middleware.RequireSuperAdmin(), middleware.RequireElevatedAuth(), h.DatabaseMaintenance)
				admin.POST("/cleanup", middleware.RequireSuperAdmin(), middleware.RequireElevatedAuth(), middleware.RequirePasswordConfirmation(), h.DatabaseCleanup)
				admin.GET("/tables/:table", middleware.RequireSuperAdmin(), h.GetTableData)
			}

			// Emergency Contact routes
			emergencyContacts := protected.Group("/emergency-contacts")
			{
				emergencyContacts.GET("", h.GetEmergencyContacts)
				emergencyContacts.GET("/:id", h.GetEmergencyContact)
				emergencyContacts.POST("", h.CreateEmergencyContact)
				emergencyContacts.PUT("/:id", h.UpdateEmergencyContact)
				emergencyContacts.DELETE("/:id", h.DeleteEmergencyContact)
			}

			// Care Plan routes
			carePlans := protected.Group("/care-plans")
			{
				carePlans.GET("", h.GetCarePlans)
				carePlans.GET("/:id", h.GetCarePlan)
				carePlans.POST("", h.CreateCarePlan)
				carePlans.PUT("/:id", h.UpdateCarePlan)
				carePlans.PATCH("/:id/approve", middleware.RequireRole("admin,manager"), h.ApproveCarePlan)
				carePlans.DELETE("/:id", h.DeleteCarePlan)
			}

			// Billing routes
			billing := protected.Group("/billing")
			{
				billing.GET("", h.GetBilling)
				billing.GET("/:id", h.GetBillingRecord)
				billing.POST("/generate", h.GenerateInvoice)
				billing.POST("/:id/payment", h.MarkAsPaid)
				billing.GET("/:id/download", h.DownloadInvoice)
			}

			// Reports routes
			reports := protected.Group("/reports")
			{
				reports.GET("/dashboard", h.GetDashboardStats)
				reports.GET("/revenue", h.GetRevenueReport)
				reports.GET("/shifts", h.GetShiftsReport)
				reports.GET("/service-hours", h.GetServiceHoursReport)
				reports.GET("/participants", h.GetParticipantReport)
				reports.GET("/staff-performance", h.GetStaffPerformance)
				reports.GET("/:type/export", h.ExportReport)
				reports.GET("/templates", h.GetReportTemplates)
			}

			// Worker availability and preferences routes
			worker := protected.Group("/worker")
			{
				workerHandler := NewWorkerAvailabilityHandler(h)
				workerHandler.RegisterRoutes(worker)
			}
		}
	}
}
