package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kenkinoti/gofiber-ago-crm-backend/internal/models"
	"gorm.io/gorm"
)

// CreateCareNoteRequest represents the request payload for creating care notes
type CreateCareNoteRequest struct {
	ParticipantID    string    `json:"participant_id" binding:"required"`
	ShiftID          *string   `json:"shift_id,omitempty"`
	Title            string    `json:"title" binding:"required"`
	Content          string    `json:"content" binding:"required"`
	NoteType         string    `json:"note_type" binding:"required"`
	Priority         string    `json:"priority"`
	NoteDate         time.Time `json:"note_date" binding:"required"`
	NoteTime         string    `json:"note_time"`
	IsPrivate        bool      `json:"is_private"`
	IsConfidential   bool      `json:"is_confidential"`
	RequiresFollowUp bool      `json:"requires_follow_up"`
	FollowUpBy       *string   `json:"follow_up_by,omitempty"`
	FollowUpDate     *time.Time `json:"follow_up_date,omitempty"`
	Tags             []string  `json:"tags,omitempty"`
	Category         *string   `json:"category,omitempty"`
}

// UpdateCareNoteRequest represents the request payload for updating care notes
type UpdateCareNoteRequest struct {
	Title            *string    `json:"title,omitempty"`
	Content          *string    `json:"content,omitempty"`
	NoteType         *string    `json:"note_type,omitempty"`
	Priority         *string    `json:"priority,omitempty"`
	NoteDate         *time.Time `json:"note_date,omitempty"`
	NoteTime         *string    `json:"note_time,omitempty"`
	IsPrivate        *bool      `json:"is_private,omitempty"`
	IsConfidential   *bool      `json:"is_confidential,omitempty"`
	RequiresFollowUp *bool      `json:"requires_follow_up,omitempty"`
	FollowUpBy       *string    `json:"follow_up_by,omitempty"`
	FollowUpDate     *time.Time `json:"follow_up_date,omitempty"`
	FollowUpStatus   *string    `json:"follow_up_status,omitempty"`
	FollowUpNotes    *string    `json:"follow_up_notes,omitempty"`
	Tags             []string   `json:"tags,omitempty"`
	Category         *string    `json:"category,omitempty"`
}

// GetCareNotes retrieves care notes with filtering options
func (h *Handler) GetCareNotes(c *gin.Context) {
	userID := h.GetUserIDFromContext(c)
	userRole := h.GetUserRoleFromContext(c)
	
	if userID == "" {
		h.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Parse query parameters
	participantID := c.Query("participant_id")
	staffID := c.Query("staff_id")
	noteType := c.Query("note_type")
	priority := c.Query("priority")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	includePrivate := c.Query("include_private") == "true"
	includeConfidential := c.Query("include_confidential") == "true"
	
	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	// Build query
	query := h.DB.Model(&models.CareNote{}).
		Preload("Participant").
		Preload("Staff").
		Preload("Shift").
		Preload("FollowUpUser")

	// Organization filtering (ensure users only see notes from their organization)
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get user information", err)
		return
	}
	query = query.Where("organization_id = ?", user.OrganizationID)

	// Apply filters
	if participantID != "" {
		query = query.Where("participant_id = ?", participantID)
	}
	
	if staffID != "" {
		query = query.Where("staff_id = ?", staffID)
	}
	
	if noteType != "" {
		query = query.Where("note_type = ?", noteType)
	}
	
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	
	if fromDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", fromDate); err == nil {
			query = query.Where("note_date >= ?", parsedDate)
		}
	}
	
	if toDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", toDate); err == nil {
			query = query.Where("note_date <= ?", parsedDate)
		}
	}

	// Privacy and confidentiality filtering based on user role
	if userRole == "care_worker" {
		// Care workers can only see their own notes and public notes
		query = query.Where("(staff_id = ? AND is_private = false) OR (staff_id != ? AND is_private = false AND is_confidential = false)", userID, userID)
	} else if !includePrivate {
		query = query.Where("is_private = false")
	}
	
	if userRole != "admin" && userRole != "super_admin" && !includeConfidential {
		query = query.Where("is_confidential = false")
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to count care notes", err)
		return
	}

	// Apply pagination and ordering
	var careNotes []models.CareNote
	if err := query.Order("note_date DESC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&careNotes).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve care notes", err)
		return
	}

	// Calculate pagination info
	totalPages := (int(total) + limit - 1) / limit

	response := gin.H{
		"care_notes": careNotes,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_count":  total,
			"limit":        limit,
		},
	}

	h.SendSuccessResponse(c, response)
}

// GetCareNote retrieves a specific care note by ID
func (h *Handler) GetCareNote(c *gin.Context) {
	userID := h.GetUserIDFromContext(c)
	userRole := h.GetUserRoleFromContext(c)
	careNoteID := c.Param("id")

	if userID == "" {
		h.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Get user's organization
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get user information", err)
		return
	}

	var careNote models.CareNote
	query := h.DB.Preload("Participant").
		Preload("Staff").
		Preload("Shift").
		Preload("FollowUpUser").
		Where("id = ? AND organization_id = ?", careNoteID, user.OrganizationID)

	if err := query.First(&careNote).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.SendErrorResponse(c, http.StatusNotFound, "Care note not found", nil)
			return
		}
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve care note", err)
		return
	}

	// Check access permissions for private/confidential notes
	if careNote.IsPrivate && careNote.StaffID != userID && userRole == "care_worker" {
		h.SendErrorResponse(c, http.StatusForbidden, "Access denied to private care note", nil)
		return
	}

	if careNote.IsConfidential && userRole != "admin" && userRole != "super_admin" {
		h.SendErrorResponse(c, http.StatusForbidden, "Access denied to confidential care note", nil)
		return
	}

	h.SendSuccessResponse(c, careNote)
}

// CreateCareNote creates a new care note
func (h *Handler) CreateCareNote(c *gin.Context) {
	userID := h.GetUserIDFromContext(c)
	if userID == "" {
		h.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Get user's organization
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get user information", err)
		return
	}

	var req CreateCareNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Validate participant exists and belongs to the same organization
	var participant models.Participant
	if err := h.DB.Where("id = ? AND organization_id = ?", req.ParticipantID, user.OrganizationID).First(&participant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.SendErrorResponse(c, http.StatusBadRequest, "Participant not found", nil)
			return
		}
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to validate participant", err)
		return
	}

	// If shift is provided, validate it exists and belongs to the participant
	if req.ShiftID != nil {
		var shift models.Shift
		if err := h.DB.Where("id = ? AND participant_id = ?", *req.ShiftID, req.ParticipantID).First(&shift).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				h.SendErrorResponse(c, http.StatusBadRequest, "Shift not found for this participant", nil)
				return
			}
			h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to validate shift", err)
			return
		}
	}

	// Convert tags to JSON string
	var tagsJSON *string
	if len(req.Tags) > 0 {
		if tagsBytes, err := json.Marshal(req.Tags); err == nil {
			tagsStr := string(tagsBytes)
			tagsJSON = &tagsStr
		}
	}

	// Set defaults
	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	// Create care note
	careNote := models.CareNote{
		ParticipantID:    req.ParticipantID,
		StaffID:          userID,
		ShiftID:          req.ShiftID,
		OrganizationID:   user.OrganizationID,
		Title:            req.Title,
		Content:          req.Content,
		NoteType:         req.NoteType,
		Priority:         priority,
		NoteDate:         req.NoteDate,
		NoteTime:         req.NoteTime,
		IsPrivate:        req.IsPrivate,
		IsConfidential:   req.IsConfidential,
		RequiresFollowUp: req.RequiresFollowUp,
		FollowUpBy:       req.FollowUpBy,
		FollowUpDate:     req.FollowUpDate,
		Tags:             tagsJSON,
		Category:         req.Category,
	}

	if err := h.DB.Create(&careNote).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create care note", err)
		return
	}

	// Load relationships before returning
	if err := h.DB.Preload("Participant").
		Preload("Staff").
		Preload("Shift").
		Preload("FollowUpUser").
		First(&careNote, careNote.ID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to load care note relationships", err)
		return
	}

	h.SendSuccessResponse(c, careNote)
}

// UpdateCareNote updates an existing care note
func (h *Handler) UpdateCareNote(c *gin.Context) {
	userID := h.GetUserIDFromContext(c)
	userRole := h.GetUserRoleFromContext(c)
	careNoteID := c.Param("id")

	if userID == "" {
		h.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Get user's organization
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get user information", err)
		return
	}

	var req UpdateCareNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Find existing care note
	var careNote models.CareNote
	if err := h.DB.Where("id = ? AND organization_id = ?", careNoteID, user.OrganizationID).First(&careNote).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.SendErrorResponse(c, http.StatusNotFound, "Care note not found", nil)
			return
		}
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to find care note", err)
		return
	}

	// Check if user can edit this note
	canEdit := false
	if userRole == "admin" || userRole == "super_admin" || userRole == "manager" {
		canEdit = true
	} else if careNote.StaffID == userID {
		canEdit = true
	}

	if !canEdit {
		h.SendErrorResponse(c, http.StatusForbidden, "You can only edit your own care notes", nil)
		return
	}

	// Prepare updates map
	updates := make(map[string]interface{})

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.NoteType != nil {
		updates["note_type"] = *req.NoteType
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.NoteDate != nil {
		updates["note_date"] = *req.NoteDate
	}
	if req.NoteTime != nil {
		updates["note_time"] = *req.NoteTime
	}
	if req.IsPrivate != nil {
		updates["is_private"] = *req.IsPrivate
	}
	if req.IsConfidential != nil {
		updates["is_confidential"] = *req.IsConfidential
	}
	if req.RequiresFollowUp != nil {
		updates["requires_follow_up"] = *req.RequiresFollowUp
	}
	if req.FollowUpBy != nil {
		updates["follow_up_by"] = req.FollowUpBy
	}
	if req.FollowUpDate != nil {
		updates["follow_up_date"] = req.FollowUpDate
	}
	if req.FollowUpStatus != nil {
		updates["follow_up_status"] = *req.FollowUpStatus
	}
	if req.FollowUpNotes != nil {
		updates["follow_up_notes"] = req.FollowUpNotes
	}
	if req.Category != nil {
		updates["category"] = req.Category
	}

	// Handle tags
	if req.Tags != nil {
		if len(req.Tags) > 0 {
			if tagsBytes, err := json.Marshal(req.Tags); err == nil {
				tagsStr := string(tagsBytes)
				updates["tags"] = &tagsStr
			}
		} else {
			updates["tags"] = nil
		}
	}

	// Apply updates
	if err := h.DB.Model(&careNote).Updates(updates).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update care note", err)
		return
	}

	// Load updated care note with relationships
	if err := h.DB.Preload("Participant").
		Preload("Staff").
		Preload("Shift").
		Preload("FollowUpUser").
		First(&careNote, careNote.ID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to load updated care note", err)
		return
	}

	h.SendSuccessResponse(c, careNote)
}

// DeleteCareNote soft deletes a care note
func (h *Handler) DeleteCareNote(c *gin.Context) {
	userID := h.GetUserIDFromContext(c)
	userRole := h.GetUserRoleFromContext(c)
	careNoteID := c.Param("id")

	if userID == "" {
		h.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Get user's organization
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get user information", err)
		return
	}

	// Find care note
	var careNote models.CareNote
	if err := h.DB.Where("id = ? AND organization_id = ?", careNoteID, user.OrganizationID).First(&careNote).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.SendErrorResponse(c, http.StatusNotFound, "Care note not found", nil)
			return
		}
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to find care note", err)
		return
	}

	// Check permissions - only admins, managers, or the original author can delete
	canDelete := false
	if userRole == "admin" || userRole == "super_admin" || userRole == "manager" {
		canDelete = true
	} else if careNote.StaffID == userID {
		canDelete = true
	}

	if !canDelete {
		h.SendErrorResponse(c, http.StatusForbidden, "You can only delete your own care notes", nil)
		return
	}

	// Soft delete
	if err := h.DB.Delete(&careNote).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete care note", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Care note deleted successfully",
	})
}

// GetCareNoteStats returns statistics about care notes
func (h *Handler) GetCareNoteStats(c *gin.Context) {
	userID := h.GetUserIDFromContext(c)
	if userID == "" {
		h.SendErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Get user's organization
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get user information", err)
		return
	}

	// Base query with organization filter
	baseQuery := h.DB.Model(&models.CareNote{}).Where("organization_id = ?", user.OrganizationID)

	// Total notes count
	var totalNotes int64
	baseQuery.Count(&totalNotes)

	// Notes by type
	var noteTypeStats []struct {
		NoteType string `json:"note_type"`
		Count    int64  `json:"count"`
	}
	h.DB.Model(&models.CareNote{}).
		Where("organization_id = ?", user.OrganizationID).
		Select("note_type, count(*) as count").
		Group("note_type").
		Scan(&noteTypeStats)

	// Notes by priority
	var priorityStats []struct {
		Priority string `json:"priority"`
		Count    int64  `json:"count"`
	}
	h.DB.Model(&models.CareNote{}).
		Where("organization_id = ?", user.OrganizationID).
		Select("priority, count(*) as count").
		Group("priority").
		Scan(&priorityStats)

	// Follow-up required count
	var followUpRequired int64
	baseQuery.Where("requires_follow_up = ? AND follow_up_status != ?", true, "completed").Count(&followUpRequired)

	stats := gin.H{
		"total_notes":        totalNotes,
		"follow_up_required": followUpRequired,
		"by_type":           noteTypeStats,
		"by_priority":       priorityStats,
	}

	h.SendSuccessResponse(c, stats)
}