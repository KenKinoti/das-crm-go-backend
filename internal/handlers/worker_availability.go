package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kenkinoti/gofiber-ago-crm-backend/internal/models"
	"gorm.io/gorm"
)

// WorkerAvailabilityHandler handles worker availability management endpoints
type WorkerAvailabilityHandler struct {
	handler *Handler
}

// NewWorkerAvailabilityHandler creates a new worker availability handler
func NewWorkerAvailabilityHandler(h *Handler) *WorkerAvailabilityHandler {
	return &WorkerAvailabilityHandler{handler: h}
}

// RegisterRoutes registers worker availability routes
func (wah *WorkerAvailabilityHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Worker availability routes
	availabilityRoutes := router.Group("/availability")
	{
		// Weekly availability patterns
		availabilityRoutes.GET("/", wah.GetAvailability)
		availabilityRoutes.POST("/", wah.CreateAvailability)
		availabilityRoutes.PUT("/:id", wah.UpdateAvailability)
		availabilityRoutes.DELETE("/:id", wah.DeleteAvailability)
		
		// Bulk operations for weekly patterns
		availabilityRoutes.POST("/bulk", wah.BulkUpdateAvailability)
		availabilityRoutes.GET("/template", wah.GetAvailabilityTemplate)
		
		// Availability exceptions (time off, special availability)
		availabilityRoutes.GET("/exceptions", wah.GetAvailabilityExceptions)
		availabilityRoutes.POST("/exceptions", wah.CreateAvailabilityException)
		availabilityRoutes.PUT("/exceptions/:id", wah.UpdateAvailabilityException)
		availabilityRoutes.DELETE("/exceptions/:id", wah.DeleteAvailabilityException)
		availabilityRoutes.POST("/exceptions/:id/approve", wah.ApproveAvailabilityException)
	}
	
	// Worker preferences routes
	preferencesRoutes := router.Group("/preferences")
	{
		preferencesRoutes.GET("/", wah.GetWorkerPreferences)
		preferencesRoutes.POST("/", wah.CreateWorkerPreferences)
		preferencesRoutes.PUT("/:id", wah.UpdateWorkerPreferences)
		preferencesRoutes.DELETE("/:id", wah.DeleteWorkerPreferences)
	}
	
	// Worker skills routes
	skillsRoutes := router.Group("/skills")
	{
		skillsRoutes.GET("/", wah.GetWorkerSkills)
		skillsRoutes.POST("/", wah.CreateWorkerSkill)
		skillsRoutes.PUT("/:id", wah.UpdateWorkerSkill)
		skillsRoutes.DELETE("/:id", wah.DeleteWorkerSkill)
		skillsRoutes.GET("/categories", wah.GetSkillCategories)
		skillsRoutes.GET("/expiring", wah.GetExpiringSkills)
	}
	
	// Location preferences routes
	locationRoutes := router.Group("/locations")
	{
		locationRoutes.GET("/", wah.GetLocationPreferences)
		locationRoutes.POST("/", wah.CreateLocationPreference)
		locationRoutes.PUT("/:id", wah.UpdateLocationPreference)
		locationRoutes.DELETE("/:id", wah.DeleteLocationPreference)
	}
	
	// Capacity and availability reports
	capacityRoutes := router.Group("/capacity")
	{
		capacityRoutes.GET("/summary", wah.GetCapacitySummary)
		capacityRoutes.GET("/weekly", wah.GetWeeklyCapacity)
		capacityRoutes.GET("/conflicts", wah.GetAvailabilityConflicts)
	}
}

// GetAvailability gets worker's weekly availability pattern
func (wah *WorkerAvailabilityHandler) GetAvailability(c *gin.Context) {
	userID := wah.handler.GetUserIDFromContext(c)
	if userID == "" {
		wah.handler.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	
	// Allow managers to view other users' availability
	targetUserID := c.Query("user_id")
	if targetUserID != "" && !wah.handler.CanUserAccessResource(c, "view_availability", targetUserID) {
		wah.handler.SendErrorResponse(c, http.StatusForbidden, "Access denied", nil)
		return
	}
	
	if targetUserID == "" {
		targetUserID = userID
	}
	
	var availabilities []models.WorkerAvailability
	query := wah.handler.DB.Where("user_id = ? AND is_active = ? AND deleted_at IS NULL", targetUserID, true)
	
	// Optional day filter
	if dayOfWeek := c.Query("day_of_week"); dayOfWeek != "" {
		if day, err := strconv.Atoi(dayOfWeek); err == nil && day >= 0 && day <= 6 {
			query = query.Where("day_of_week = ?", day)
		}
	}
	
	if err := query.Order("day_of_week ASC, start_time ASC").Find(&availabilities).Error; err != nil {
		wah.handler.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch availability", err)
		return
	}
	
	// Convert to DTOs
	var dtos []models.WorkerAvailabilityDTO
	for _, availability := range availabilities {
		dtos = append(dtos, availability.ToDTO())
	}
	
	wah.handler.SendSuccessResponse(c, gin.H{
		"availability": dtos,
		"user_id":      targetUserID,
	})
}

// CreateAvailability creates a new availability slot
func (wah *WorkerAvailabilityHandler) CreateAvailability(c *gin.Context) {
	userID := wah.handler.GetUserIDFromContext(c)
	if userID == "" {
		wah.handler.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	
	var req struct {
		UserID               string   `json:"user_id"`
		DayOfWeek            int      `json:"day_of_week" binding:"required,min=0,max=6"`
		StartTime            string   `json:"start_time" binding:"required"`
		EndTime              string   `json:"end_time" binding:"required"`
		MaxHoursPerDay       *float64 `json:"max_hours_per_day"`
		PreferredHoursPerDay *float64 `json:"preferred_hours_per_day"`
		HourlyRate           *float64 `json:"hourly_rate"`
		BreakDurationMinutes *int     `json:"break_duration_minutes"`
		TravelTimeMinutes    *int     `json:"travel_time_minutes"`
		Notes                string   `json:"notes"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		wah.handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}
	
	// Use current user if no user_id specified or check permissions
	targetUserID := userID
	if req.UserID != "" && req.UserID != userID {
		if !wah.handler.CanUserAccessResource(c, "manage_availability", req.UserID) {
			wah.handler.SendErrorResponse(c, http.StatusForbidden, "Access denied", nil)
			return
		}
		targetUserID = req.UserID
	}
	
	// Parse time strings
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		wah.handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid start time format", err)
		return
	}
	
	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		wah.handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid end time format", err)
		return
	}
	
	// Validate time range
	if !endTime.After(startTime) {
		wah.handler.SendErrorResponse(c, http.StatusBadRequest, "End time must be after start time", nil)
		return
	}
	
	availability := models.WorkerAvailability{
		UserID:      targetUserID,
		DayOfWeek:   req.DayOfWeek,
		StartTime:   startTime,
		EndTime:     endTime,
		IsAvailable: true,
		IsActive:    true,
	}
	
	// Set optional fields with defaults
	if req.MaxHoursPerDay != nil {
		availability.MaxHoursPerDay = *req.MaxHoursPerDay
	} else {
		availability.MaxHoursPerDay = 8.0
	}
	
	if req.PreferredHoursPerDay != nil {
		availability.PreferredHoursPerDay = *req.PreferredHoursPerDay
	} else {
		availability.PreferredHoursPerDay = 6.0
	}
	
	if req.HourlyRate != nil {
		availability.HourlyRate = req.HourlyRate
	}
	
	if req.BreakDurationMinutes != nil {
		availability.BreakDurationMinutes = *req.BreakDurationMinutes
	} else {
		availability.BreakDurationMinutes = 30
	}
	
	if req.TravelTimeMinutes != nil {
		availability.TravelTimeMinutes = *req.TravelTimeMinutes
	} else {
		availability.TravelTimeMinutes = 30
	}
	
	availability.Notes = req.Notes
	
	if err := wah.handler.DB.Create(&availability).Error; err != nil {
		wah.handler.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create availability", err)
		return
	}
	
	wah.handler.SendSuccessResponse(c, gin.H{
		"availability": availability.ToDTO(),
		"message":      "Availability created successfully",
	})
}

// BulkUpdateAvailability handles bulk update of weekly availability
func (wah *WorkerAvailabilityHandler) BulkUpdateAvailability(c *gin.Context) {
	userID := wah.handler.GetUserIDFromContext(c)
	if userID == "" {
		wah.handler.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	
	var req struct {
		UserID       string                          `json:"user_id"`
		Availability []models.WorkerAvailabilityDTO `json:"availability" binding:"required"`
		ReplaceAll   bool                            `json:"replace_all"` // If true, replace all existing availability
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		wah.handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}
	
	// Use current user if no user_id specified or check permissions
	targetUserID := userID
	if req.UserID != "" && req.UserID != userID {
		if !wah.handler.CanUserAccessResource(c, "manage_availability", req.UserID) {
			wah.handler.SendErrorResponse(c, http.StatusForbidden, "Access denied", nil)
			return
		}
		targetUserID = req.UserID
	}
	
	// Start transaction
	tx := wah.handler.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// If replace_all is true, soft delete existing availability
	if req.ReplaceAll {
		if err := tx.Model(&models.WorkerAvailability{}).
			Where("user_id = ? AND deleted_at IS NULL", targetUserID).
			Update("deleted_at", time.Now()).Error; err != nil {
			tx.Rollback()
			wah.handler.SendErrorResponse(c, http.StatusInternalServerError, "Failed to clear existing availability", err)
			return
		}
	}
	
	var createdAvailability []models.WorkerAvailabilityDTO
	
	for _, avail := range req.Availability {
		// Parse time strings
		startTime, err := time.Parse("15:04", avail.StartTime)
		if err != nil {
			tx.Rollback()
			wah.handler.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid start time format: %s", avail.StartTime), err)
			return
		}
		
		endTime, err := time.Parse("15:04", avail.EndTime)
		if err != nil {
			tx.Rollback()
			wah.handler.SendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid end time format: %s", avail.EndTime), err)
			return
		}
		
		availability := models.WorkerAvailability{
			UserID:               targetUserID,
			DayOfWeek:            avail.DayOfWeek,
			StartTime:            startTime,
			EndTime:              endTime,
			IsAvailable:          avail.IsAvailable,
			MaxHoursPerDay:       avail.MaxHoursPerDay,
			PreferredHoursPerDay: avail.PreferredHoursPerDay,
			HourlyRate:           avail.HourlyRate,
			BreakDurationMinutes: avail.BreakDurationMinutes,
			TravelTimeMinutes:    avail.TravelTimeMinutes,
			Notes:                avail.Notes,
			IsActive:             true,
		}
		
		if err := tx.Create(&availability).Error; err != nil {
			tx.Rollback()
			wah.handler.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create availability", err)
			return
		}
		
		createdAvailability = append(createdAvailability, availability.ToDTO())
	}
	
	if err := tx.Commit().Error; err != nil {
		wah.handler.SendErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}
	
	wah.handler.SendSuccessResponse(c, gin.H{
		"availability": createdAvailability,
		"message":      fmt.Sprintf("Successfully updated %d availability slots", len(createdAvailability)),
	})
}

// GetWorkerPreferences gets worker preferences and capacity settings
func (wah *WorkerAvailabilityHandler) GetWorkerPreferences(c *gin.Context) {
	userID := wah.handler.GetUserIDFromContext(c)
	if userID == "" {
		wah.handler.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	
	// Allow managers to view other users' preferences
	targetUserID := c.Query("user_id")
	if targetUserID != "" && !wah.handler.CanUserAccessResource(c, "view_preferences", targetUserID) {
		wah.handler.SendErrorResponse(c, http.StatusForbidden, "Access denied", nil)
		return
	}
	
	if targetUserID == "" {
		targetUserID = userID
	}
	
	var preferences models.WorkerPreferences
	err := wah.handler.DB.Where("user_id = ? AND is_active = ?", targetUserID, true).First(&preferences).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default preferences if none exist
			defaultPrefs := models.WorkerPreferences{
				UserID:                  targetUserID,
				MaxHoursPerWeek:         38.0,
				PreferredHoursPerWeek:   30.0,
				MaxConsecutiveDays:      5,
				MinHoursBetweenShifts:   10,
				MaxTravelDistanceKm:     50,
				WillingWeekendWork:      true,
				WillingEveningWork:      true,
				WillingEarlyMorningWork: true,
				HasOwnVehicle:           true,
				IsActive:                true,
			}
			
			wah.handler.SendSuccessResponse(c, gin.H{
				"preferences": defaultPrefs.ToDTO(),
				"is_default":  true,
			})
			return
		}
		
		wah.handler.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch preferences", err)
		return
	}
	
	wah.handler.SendSuccessResponse(c, gin.H{
		"preferences": preferences.ToDTO(),
		"is_default":  false,
	})
}

// CreateWorkerPreferences creates or updates worker preferences
func (wah *WorkerAvailabilityHandler) CreateWorkerPreferences(c *gin.Context) {
	userID := wah.handler.GetUserIDFromContext(c)
	if userID == "" {
		wah.handler.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	
	var req models.WorkerPreferences
	if err := c.ShouldBindJSON(&req); err != nil {
		wah.handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}
	
	// Use current user if no user_id specified or check permissions
	targetUserID := userID
	if req.UserID != "" && req.UserID != userID {
		if !wah.handler.CanUserAccessResource(c, "manage_preferences", req.UserID) {
			wah.handler.SendErrorResponse(c, http.StatusForbidden, "Access denied", nil)
			return
		}
		targetUserID = req.UserID
	}
	req.UserID = targetUserID
	
	// Try to update existing preferences first
	var existingPrefs models.WorkerPreferences
	err := wah.handler.DB.Where("user_id = ?", targetUserID).First(&existingPrefs).Error
	
	if err == nil {
		// Update existing preferences
		req.ID = existingPrefs.ID
		req.CreatedAt = existingPrefs.CreatedAt
		
		if err := wah.handler.DB.Save(&req).Error; err != nil {
			wah.handler.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update preferences", err)
			return
		}
	} else if err == gorm.ErrRecordNotFound {
		// Create new preferences
		req.IsActive = true
		if err := wah.handler.DB.Create(&req).Error; err != nil {
			wah.handler.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create preferences", err)
			return
		}
	} else {
		wah.handler.SendErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}
	
	wah.handler.SendSuccessResponse(c, gin.H{
		"preferences": req.ToDTO(),
		"message":     "Preferences saved successfully",
	})
}

// GetCapacitySummary provides a summary of worker capacity and availability
func (wah *WorkerAvailabilityHandler) GetCapacitySummary(c *gin.Context) {
	userID := wah.handler.GetUserIDFromContext(c)
	if userID == "" {
		wah.handler.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	
	// Allow managers to view other users' capacity
	targetUserID := c.Query("user_id")
	if targetUserID != "" && !wah.handler.CanUserAccessResource(c, "view_capacity", targetUserID) {
		wah.handler.SendErrorResponse(c, http.StatusForbidden, "Access denied", nil)
		return
	}
	
	if targetUserID == "" {
		targetUserID = userID
	}
	
	// Get worker preferences
	var preferences models.WorkerPreferences
	prefErr := wah.handler.DB.Where("user_id = ? AND is_active = ?", targetUserID, true).First(&preferences).Error
	
	// Get availability
	var availabilities []models.WorkerAvailability
	availErr := wah.handler.DB.Where("user_id = ? AND is_active = ? AND deleted_at IS NULL", targetUserID, true).
		Order("day_of_week ASC").Find(&availabilities).Error
	
	// Get skills
	var skills []models.WorkerSkill
	wah.handler.DB.Where("user_id = ? AND is_active = ? AND deleted_at IS NULL", targetUserID, true).
		Order("skill_category ASC, skill_name ASC").Find(&skills)
	
	// Calculate weekly capacity
	totalWeeklyHours := 0.0
	totalPreferredHours := 0.0
	daysAvailable := 0
	
	for _, avail := range availabilities {
		if avail.IsAvailable {
			totalWeeklyHours += avail.GetTotalAvailableHours()
			totalPreferredHours += avail.PreferredHoursPerDay
			daysAvailable++
		}
	}
	
	// Count expiring skills
	expiringSkills := 0
	expiredSkills := 0
	for _, skill := range skills {
		if skill.IsSkillExpiringSoon() {
			expiringSkills++
		}
		if skill.IsSkillExpired() {
			expiredSkills++
		}
	}
	
	summary := gin.H{
		"user_id": targetUserID,
		"capacity": gin.H{
			"weekly_available_hours":   totalWeeklyHours,
			"weekly_preferred_hours":   totalPreferredHours,
			"days_available_per_week":  daysAvailable,
			"capacity_utilization":     0.0, // Would need shift data to calculate
		},
		"preferences_configured": prefErr == nil,
		"availability_configured": availErr == nil && len(availabilities) > 0,
		"skills": gin.H{
			"total_skills":    len(skills),
			"expiring_skills": expiringSkills,
			"expired_skills":  expiredSkills,
		},
		"status": "active",
	}
	
	if prefErr == nil {
		summary["preferences"] = preferences.ToDTO()
		if totalWeeklyHours > 0 {
			capacity := summary["capacity"].(gin.H)
			capacity["capacity_utilization"] = (preferences.PreferredHoursPerWeek / totalWeeklyHours) * 100
		}
	}
	
	wah.handler.SendSuccessResponse(c, summary)
}

// Additional helper methods would go here (UpdateAvailability, DeleteAvailability, 
// CreateAvailabilityException, GetWorkerSkills, CreateWorkerSkill, etc.)
// For brevity, I'll include a few key ones:

// UpdateAvailability updates an existing availability slot
func (wah *WorkerAvailabilityHandler) UpdateAvailability(c *gin.Context) {
	userID := wah.handler.GetUserIDFromContext(c)
	if userID == "" {
		wah.handler.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	
	availabilityID := c.Param("id")
	if availabilityID == "" {
		wah.handler.SendErrorResponse(c, http.StatusBadRequest, "Availability ID is required", nil)
		return
	}
	
	var availability models.WorkerAvailability
	if err := wah.handler.DB.Where("id = ? AND deleted_at IS NULL", availabilityID).First(&availability).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			wah.handler.SendErrorResponse(c, http.StatusNotFound, "Availability not found", nil)
			return
		}
		wah.handler.SendErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}
	
	// Check permissions
	if availability.UserID != userID && !wah.handler.CanUserAccessResource(c, "manage_availability", availability.UserID) {
		wah.handler.SendErrorResponse(c, http.StatusForbidden, "Access denied", nil)
		return
	}
	
	var req struct {
		DayOfWeek            *int     `json:"day_of_week"`
		StartTime            *string  `json:"start_time"`
		EndTime              *string  `json:"end_time"`
		IsAvailable          *bool    `json:"is_available"`
		MaxHoursPerDay       *float64 `json:"max_hours_per_day"`
		PreferredHoursPerDay *float64 `json:"preferred_hours_per_day"`
		HourlyRate           *float64 `json:"hourly_rate"`
		BreakDurationMinutes *int     `json:"break_duration_minutes"`
		TravelTimeMinutes    *int     `json:"travel_time_minutes"`
		Notes                *string  `json:"notes"`
		IsActive             *bool    `json:"is_active"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		wah.handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}
	
	// Update fields if provided
	if req.DayOfWeek != nil {
		if *req.DayOfWeek < 0 || *req.DayOfWeek > 6 {
			wah.handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid day of week", nil)
			return
		}
		availability.DayOfWeek = *req.DayOfWeek
	}
	
	if req.StartTime != nil {
		startTime, err := time.Parse("15:04", *req.StartTime)
		if err != nil {
			wah.handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid start time format", err)
			return
		}
		availability.StartTime = startTime
	}
	
	if req.EndTime != nil {
		endTime, err := time.Parse("15:04", *req.EndTime)
		if err != nil {
			wah.handler.SendErrorResponse(c, http.StatusBadRequest, "Invalid end time format", err)
			return
		}
		availability.EndTime = endTime
	}
	
	// Validate time range if both times are being updated
	if !availability.EndTime.After(availability.StartTime) {
		wah.handler.SendErrorResponse(c, http.StatusBadRequest, "End time must be after start time", nil)
		return
	}
	
	if req.IsAvailable != nil {
		availability.IsAvailable = *req.IsAvailable
	}
	if req.MaxHoursPerDay != nil {
		availability.MaxHoursPerDay = *req.MaxHoursPerDay
	}
	if req.PreferredHoursPerDay != nil {
		availability.PreferredHoursPerDay = *req.PreferredHoursPerDay
	}
	if req.HourlyRate != nil {
		availability.HourlyRate = req.HourlyRate
	}
	if req.BreakDurationMinutes != nil {
		availability.BreakDurationMinutes = *req.BreakDurationMinutes
	}
	if req.TravelTimeMinutes != nil {
		availability.TravelTimeMinutes = *req.TravelTimeMinutes
	}
	if req.Notes != nil {
		availability.Notes = *req.Notes
	}
	if req.IsActive != nil {
		availability.IsActive = *req.IsActive
	}
	
	if err := wah.handler.DB.Save(&availability).Error; err != nil {
		wah.handler.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update availability", err)
		return
	}
	
	wah.handler.SendSuccessResponse(c, gin.H{
		"availability": availability.ToDTO(),
		"message":      "Availability updated successfully",
	})
}

// DeleteAvailability soft deletes an availability slot
func (wah *WorkerAvailabilityHandler) DeleteAvailability(c *gin.Context) {
	userID := wah.handler.GetUserIDFromContext(c)
	if userID == "" {
		wah.handler.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	
	availabilityID := c.Param("id")
	if availabilityID == "" {
		wah.handler.SendErrorResponse(c, http.StatusBadRequest, "Availability ID is required", nil)
		return
	}
	
	var availability models.WorkerAvailability
	if err := wah.handler.DB.Where("id = ? AND deleted_at IS NULL", availabilityID).First(&availability).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			wah.handler.SendErrorResponse(c, http.StatusNotFound, "Availability not found", nil)
			return
		}
		wah.handler.SendErrorResponse(c, http.StatusInternalServerError, "Database error", err)
		return
	}
	
	// Check permissions
	if availability.UserID != userID && !wah.handler.CanUserAccessResource(c, "manage_availability", availability.UserID) {
		wah.handler.SendErrorResponse(c, http.StatusForbidden, "Access denied", nil)
		return
	}
	
	// Soft delete
	if err := wah.handler.DB.Model(&availability).Update("deleted_at", time.Now()).Error; err != nil {
		wah.handler.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete availability", err)
		return
	}
	
	wah.handler.SendSuccessResponse(c, gin.H{
		"message": "Availability deleted successfully",
	})
}

// GetAvailabilityTemplate returns a template for weekly availability setup
func (wah *WorkerAvailabilityHandler) GetAvailabilityTemplate(c *gin.Context) {
	template := []gin.H{
		{"day_of_week": 1, "day_name": "Monday", "start_time": "08:00", "end_time": "17:00", "is_available": true, "max_hours_per_day": 8.0, "preferred_hours_per_day": 6.0},
		{"day_of_week": 2, "day_name": "Tuesday", "start_time": "08:00", "end_time": "17:00", "is_available": true, "max_hours_per_day": 8.0, "preferred_hours_per_day": 6.0},
		{"day_of_week": 3, "day_name": "Wednesday", "start_time": "08:00", "end_time": "17:00", "is_available": true, "max_hours_per_day": 8.0, "preferred_hours_per_day": 6.0},
		{"day_of_week": 4, "day_name": "Thursday", "start_time": "08:00", "end_time": "17:00", "is_available": true, "max_hours_per_day": 8.0, "preferred_hours_per_day": 6.0},
		{"day_of_week": 5, "day_name": "Friday", "start_time": "08:00", "end_time": "17:00", "is_available": true, "max_hours_per_day": 8.0, "preferred_hours_per_day": 6.0},
		{"day_of_week": 6, "day_name": "Saturday", "start_time": "09:00", "end_time": "14:00", "is_available": false, "max_hours_per_day": 5.0, "preferred_hours_per_day": 4.0},
		{"day_of_week": 0, "day_name": "Sunday", "start_time": "10:00", "end_time": "16:00", "is_available": false, "max_hours_per_day": 6.0, "preferred_hours_per_day": 4.0},
	}
	
	wah.handler.SendSuccessResponse(c, gin.H{
		"template": template,
		"message":  "Default weekly availability template",
	})
}

// Placeholder methods for remaining endpoints
func (wah *WorkerAvailabilityHandler) GetAvailabilityExceptions(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) CreateAvailabilityException(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) UpdateAvailabilityException(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) DeleteAvailabilityException(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) ApproveAvailabilityException(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) UpdateWorkerPreferences(c *gin.Context) {
	// This would be similar to CreateWorkerPreferences but with ID param
	wah.CreateWorkerPreferences(c) // For now, reuse create logic
}

func (wah *WorkerAvailabilityHandler) DeleteWorkerPreferences(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) GetWorkerSkills(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) CreateWorkerSkill(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) UpdateWorkerSkill(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) DeleteWorkerSkill(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) GetSkillCategories(c *gin.Context) {
	categories := []string{
		"Safety", "Healthcare", "Physical", "Communication", "Transport", 
		"Administration", "Technology", "Specialized Care", "Emergency Response",
	}
	
	wah.handler.SendSuccessResponse(c, gin.H{
		"categories": categories,
	})
}

func (wah *WorkerAvailabilityHandler) GetExpiringSkills(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) GetLocationPreferences(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) CreateLocationPreference(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) UpdateLocationPreference(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) DeleteLocationPreference(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) GetWeeklyCapacity(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}

func (wah *WorkerAvailabilityHandler) GetAvailabilityConflicts(c *gin.Context) {
	wah.handler.SendErrorResponse(c, http.StatusNotImplemented, "Not implemented yet", nil)
}