package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kenkinoti/gofiber-ago-crm-backend/internal/models"
	"gorm.io/gorm"
)

// IncidentReportRequest represents the request structure for incident reports
type IncidentReportRequest struct {
	ParticipantID       string    `json:"participant_id" binding:"required"`
	IncidentDate        time.Time `json:"incident_date" binding:"required"`
	IncidentTime        string    `json:"incident_time" binding:"required"`
	Location            string    `json:"location" binding:"required"`
	IncidentType        string    `json:"incident_type" binding:"required"`
	Severity            string    `json:"severity" binding:"required"`
	Description         string    `json:"description" binding:"required"`
	ImmediateAction     string    `json:"immediate_action"`
	InjuriesDescription *string   `json:"injuries_description,omitempty"`
	WitnessesPresent    bool      `json:"witnesses_present"`
	WitnessDetails      *string   `json:"witness_details,omitempty"`
	MedicalAttention    bool      `json:"medical_attention"`
	MedicalDetails      *string   `json:"medical_details,omitempty"`
	EmergencyServices   bool      `json:"emergency_services"`
	EmergencyDetails    *string   `json:"emergency_details,omitempty"`
	FamilyNotified      bool      `json:"family_notified"`
	PreventiveMeasures  *string   `json:"preventive_measures,omitempty"`
	FollowUpRequired    bool      `json:"follow_up_required"`
	FollowUpDetails     *string   `json:"follow_up_details,omitempty"`
}

// EnhancedIncidentReportRequest represents the comprehensive incident report request
type EnhancedIncidentReportRequest struct {
	// Basic incident info
	ParticipantID           string    `json:"participant_id" binding:"required"`
	IncidentDate            time.Time `json:"incident_date" binding:"required"`
	IncidentTime            string    `json:"incident_time" binding:"required"`
	Location                string    `json:"location" binding:"required"`
	IncidentType            string    `json:"incident_type" binding:"required"`
	Severity                string    `json:"severity" binding:"required"`
	Description             string    `json:"description" binding:"required"`
	ImmediateAction         string    `json:"immediate_action"`
	
	// Enhanced fields
	EnvironmentConditions   *string   `json:"environment_conditions,omitempty"`
	WeatherConditions       *string   `json:"weather_conditions,omitempty"`
	ContributingFactors     *string   `json:"contributing_factors,omitempty"`
	CorrectiveActions       *string   `json:"corrective_actions,omitempty"`
	
	// People involved
	PeopleInvolved          []PersonInvolvedRequest `json:"people_involved"`
	
	// Witnesses  
	Witnesses               []WitnessRequest        `json:"witnesses"`
	
	// Notifications
	ManagementNotified      bool      `json:"management_notified"`
	PoliceNotified          bool      `json:"police_notified"`
	PoliceReference         *string   `json:"police_reference,omitempty"`
	InsuranceNotified       bool      `json:"insurance_notified"`
	InsuranceReference      *string   `json:"insurance_reference,omitempty"`
	RegulatoryBodyNotified  bool      `json:"regulatory_body_notified"`
	RegulatoryReference     *string   `json:"regulatory_reference,omitempty"`
	
	// Timeline entries
	Timeline                []TimelineEntryRequest  `json:"timeline"`
	
	// Risk factors
	RiskFactors             []RiskFactorRequest     `json:"risk_factors"`
	
	// Follow-up actions
	FollowUpActions         []FollowUpActionRequest `json:"follow_up_actions"`
}

// PersonInvolvedRequest for people involved in the incident
type PersonInvolvedRequest struct {
	PersonType                 string  `json:"person_type" binding:"required"`
	ParticipantID              *string `json:"participant_id,omitempty"`
	StaffUserID                *string `json:"staff_user_id,omitempty"`
	FirstName                  *string `json:"first_name,omitempty"`
	LastName                   *string `json:"last_name,omitempty"`
	Age                        *int    `json:"age,omitempty"`
	Gender                     *string `json:"gender,omitempty"`
	ContactPhone               *string `json:"contact_phone,omitempty"`
	ContactEmail               *string `json:"contact_email,omitempty"`
	RelationshipToParticipant  *string `json:"relationship_to_participant,omitempty"`
	EmployerOrganization       *string `json:"employer_organization,omitempty"`
	RoleInIncident             string  `json:"role_in_incident" binding:"required"`
	WasInjured                 bool    `json:"was_injured"`
	InjuryDescription          *string `json:"injury_description,omitempty"`
	InjurySeverity             *string `json:"injury_severity,omitempty"`
	MedicalAttentionRequired   bool    `json:"medical_attention_required"`
	MedicalAttentionProvided   bool    `json:"medical_attention_provided"`
	MedicalProvider            *string `json:"medical_provider,omitempty"`
	TransportedToHospital      bool    `json:"transported_to_hospital"`
	HospitalName               *string `json:"hospital_name,omitempty"`
	AmbulanceCalled            bool    `json:"ambulance_called"`
	ImmediateCareProvided      *string `json:"immediate_care_provided,omitempty"`
	OngoingCareRequired        bool    `json:"ongoing_care_required"`
	OngoingCareDetails         *string `json:"ongoing_care_details,omitempty"`
	Injuries                   []InjuryRequest `json:"injuries"`
}

// InjuryRequest for detailed injury tracking
type InjuryRequest struct {
	InjuryType              string  `json:"injury_type" binding:"required"`
	InjuryCategory          string  `json:"injury_category" binding:"required"`
	BodyPart                string  `json:"body_part" binding:"required"`
	BodySide                *string `json:"body_side,omitempty"`
	SpecificLocation        *string `json:"specific_location,omitempty"`
	SeverityLevel           string  `json:"severity_level" binding:"required"`
	SizeDimensions          *string `json:"size_dimensions,omitempty"`
	DepthDescription        *string `json:"depth_description,omitempty"`
	CauseOfInjury           string  `json:"cause_of_injury" binding:"required"`
	MechanismDescription    *string `json:"mechanism_description,omitempty"`
	FirstAidGiven           bool    `json:"first_aid_given"`
	FirstAidDetails         *string `json:"first_aid_details,omitempty"`
	MedicalTreatmentRequired bool   `json:"medical_treatment_required"`
	TreatmentProvided       *string `json:"treatment_provided,omitempty"`
	OngoingTreatmentRequired bool   `json:"ongoing_treatment_required"`
	FollowUpAppointmentRequired bool `json:"follow_up_appointment_required"`
	FollowUpAppointmentScheduled bool `json:"follow_up_appointment_scheduled"`
	FollowUpDetails         *string `json:"follow_up_details,omitempty"`
	PhotosTaken             bool    `json:"photos_taken"`
	WitnessStatementsTaken  bool    `json:"witness_statements_taken"`
}

// WitnessRequest for witness information
type WitnessRequest struct {
	WitnessType        string  `json:"witness_type" binding:"required"`
	StaffUserID        *string `json:"staff_user_id,omitempty"`
	ParticipantID      *string `json:"participant_id,omitempty"`
	FirstName          *string `json:"first_name,omitempty"`
	LastName           *string `json:"last_name,omitempty"`
	ContactPhone       *string `json:"contact_phone,omitempty"`
	ContactEmail       *string `json:"contact_email,omitempty"`
	Relationship       *string `json:"relationship,omitempty"`
	WitnessStatement   *string `json:"witness_statement,omitempty"`
	StatementMethod    *string `json:"statement_method,omitempty"`
	WitnessCredibility *string `json:"witness_credibility,omitempty"`
	Notes              *string `json:"notes,omitempty"`
}

// TimelineEntryRequest for incident timeline
type TimelineEntryRequest struct {
	SequenceNumber     int    `json:"sequence_number" binding:"required"`
	TimestampRecorded  time.Time `json:"timestamp_recorded" binding:"required"`
	EventDescription   string `json:"event_description" binding:"required"`
	EventType          string `json:"event_type" binding:"required"`
	Source             *string `json:"source,omitempty"`
}

// RiskFactorRequest for risk assessment
type RiskFactorRequest struct {
	FactorType         string  `json:"factor_type" binding:"required"`
	FactorCategory     string  `json:"factor_category" binding:"required"`
	Description        string  `json:"description" binding:"required"`
	SeverityImpact     *string `json:"severity_impact,omitempty"`
	Likelihood         *string `json:"likelihood,omitempty"`
	RiskRating         *string `json:"risk_rating,omitempty"`
	ExistingControls   *string `json:"existing_controls,omitempty"`
	RecommendedActions *string `json:"recommended_actions,omitempty"`
	ActionPriority     *string `json:"action_priority,omitempty"`
	ActionAssignedTo   *string `json:"action_assigned_to,omitempty"`
	ActionDueDate      *time.Time `json:"action_due_date,omitempty"`
}

// FollowUpActionRequest for follow-up actions
type FollowUpActionRequest struct {
	ActionType         string    `json:"action_type" binding:"required"`
	ActionCategory     string    `json:"action_category" binding:"required"`
	Title              string    `json:"title" binding:"required"`
	Description        string    `json:"description" binding:"required"`
	Priority           string    `json:"priority" binding:"required"`
	AssignedTo         string    `json:"assigned_to" binding:"required"`
	DueDate            time.Time `json:"due_date" binding:"required"`
	EstimatedHours     *int      `json:"estimated_hours,omitempty"`
	RequiresApproval   bool      `json:"requires_approval"`
}

// IncidentReportUpdateRequest represents the request structure for updating incident reports (support coordinator)
type IncidentReportUpdateRequest struct {
	Status            string     `json:"status"`
	Priority          string     `json:"priority"`
	ReviewNotes       *string    `json:"review_notes,omitempty"`
	NDISNotified      bool       `json:"ndis_notified"`
	NDISReference     *string    `json:"ndis_reference,omitempty"`
	PreventiveMeasures *string   `json:"preventive_measures,omitempty"`
	FollowUpRequired  bool       `json:"follow_up_required"`
	FollowUpDetails   *string    `json:"follow_up_details,omitempty"`
}

// CreateIncidentReport creates a new incident report (submitted by care worker)
func (h *Handler) CreateIncidentReport(c *gin.Context) {
	var req IncidentReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// Get user info from context
	userID := h.GetUserIDFromContext(c)
	if userID == "" {
		h.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get user to get organization ID
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	// Verify participant belongs to the same organization
	var participant models.Participant
	if err := h.DB.First(&participant, "id = ? AND organization_id = ?", req.ParticipantID, user.OrganizationID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusNotFound, "Participant not found", err)
		return
	}

	// Create incident report
	incidentReport := models.IncidentReport{
		ParticipantID:       req.ParticipantID,
		ReportedBy:          userID,
		OrganizationID:      user.OrganizationID,
		IncidentDate:        req.IncidentDate,
		IncidentTime:        req.IncidentTime,
		Location:            req.Location,
		IncidentType:        req.IncidentType,
		Severity:            req.Severity,
		Description:         req.Description,
		ImmediateAction:     req.ImmediateAction,
		InjuriesDescription: req.InjuriesDescription,
		WitnessesPresent:    req.WitnessesPresent,
		WitnessDetails:      req.WitnessDetails,
		MedicalAttention:    req.MedicalAttention,
		MedicalDetails:      req.MedicalDetails,
		EmergencyServices:   req.EmergencyServices,
		EmergencyDetails:    req.EmergencyDetails,
		FamilyNotified:      req.FamilyNotified,
		PreventiveMeasures:  req.PreventiveMeasures,
		FollowUpRequired:    req.FollowUpRequired,
		FollowUpDetails:     req.FollowUpDetails,
		Status:              "submitted",
		Priority:            "medium",
	}

	// Set family notification timestamp if notified
	if req.FamilyNotified {
		now := time.Now()
		incidentReport.FamilyNotifiedAt = &now
	}

	if err := h.DB.Create(&incidentReport).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create incident report", err)
		return
	}

	// Load relationships
	h.DB.Preload("Participant").Preload("Reporter").Find(&incidentReport)

	h.SendSuccessResponse(c, incidentReport)
}

// CreateEnhancedIncidentReport creates a comprehensive incident report with all related entities
func (h *Handler) CreateEnhancedIncidentReport(c *gin.Context) {
	var req EnhancedIncidentReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// Get user info from context
	userID := h.GetUserIDFromContext(c)
	if userID == "" {
		h.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get user to get organization ID
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	// Verify participant belongs to the same organization
	var participant models.Participant
	if err := h.DB.First(&participant, "id = ? AND organization_id = ?", req.ParticipantID, user.OrganizationID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusNotFound, "Participant not found", err)
		return
	}

	// Begin transaction
	tx := h.DB.Begin()
	if tx.Error != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to begin transaction", tx.Error)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create incident report
	incidentReport := models.IncidentReport{
		ParticipantID:         req.ParticipantID,
		ReportedBy:           userID,
		OrganizationID:       user.OrganizationID,
		IncidentDate:         req.IncidentDate,
		IncidentTime:         req.IncidentTime,
		Location:             req.Location,
		IncidentType:         req.IncidentType,
		Severity:             req.Severity,
		Description:          req.Description,
		ImmediateAction:      req.ImmediateAction,
		Status:               "submitted",
		Priority:             "medium",
	}

	// Basic incident report created - additional fields will be handled separately
	now := time.Now()

	if err := tx.Create(&incidentReport).Error; err != nil {
		tx.Rollback()
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create incident report", err)
		return
	}

	// Create people involved
	for _, personReq := range req.PeopleInvolved {
		person := models.IncidentPersonInvolved{
			IncidentReportID:         incidentReport.ID,
			PersonType:              personReq.PersonType,
			ParticipantID:           personReq.ParticipantID,
			StaffUserID:             personReq.StaffUserID,
			FirstName:               personReq.FirstName,
			LastName:                personReq.LastName,
			Age:                     personReq.Age,
			Gender:                  personReq.Gender,
			ContactPhone:            personReq.ContactPhone,
			ContactEmail:            personReq.ContactEmail,
			RelationshipToParticipant: personReq.RelationshipToParticipant,
			EmployerOrganization:    personReq.EmployerOrganization,
			RoleInIncident:          personReq.RoleInIncident,
			WasInjured:              personReq.WasInjured,
			InjuryDescription:       personReq.InjuryDescription,
			InjurySeverity:          personReq.InjurySeverity,
			MedicalAttentionRequired: personReq.MedicalAttentionRequired,
			MedicalAttentionProvided: personReq.MedicalAttentionProvided,
			MedicalProvider:         personReq.MedicalProvider,
			TransportedToHospital:   personReq.TransportedToHospital,
			HospitalName:            personReq.HospitalName,
			AmbulanceCalled:         personReq.AmbulanceCalled,
			ImmediateCareProvided:   personReq.ImmediateCareProvided,
			OngoingCareRequired:     personReq.OngoingCareRequired,
			OngoingCareDetails:      personReq.OngoingCareDetails,
		}

		if err := tx.Create(&person).Error; err != nil {
			tx.Rollback()
			h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create person involved", err)
			return
		}

		// Create injuries for this person
		for _, injuryReq := range personReq.Injuries {
			injury := models.IncidentInjury{
				IncidentReportID:            incidentReport.ID,
				PersonInvolvedID:            person.ID,
				InjuryType:                 injuryReq.InjuryType,
				InjuryCategory:             injuryReq.InjuryCategory,
				BodyPart:                   injuryReq.BodyPart,
				BodySide:                   injuryReq.BodySide,
				SpecificLocation:           injuryReq.SpecificLocation,
				SeverityLevel:              injuryReq.SeverityLevel,
				SizeDimensions:             injuryReq.SizeDimensions,
				DepthDescription:           injuryReq.DepthDescription,
				CauseOfInjury:              injuryReq.CauseOfInjury,
				MechanismDescription:       injuryReq.MechanismDescription,
				FirstAidGiven:              injuryReq.FirstAidGiven,
				FirstAidDetails:            injuryReq.FirstAidDetails,
				MedicalTreatmentRequired:   injuryReq.MedicalTreatmentRequired,
				TreatmentProvided:          injuryReq.TreatmentProvided,
				OngoingTreatmentRequired:   injuryReq.OngoingTreatmentRequired,
				FollowUpAppointmentRequired: injuryReq.FollowUpAppointmentRequired,
				FollowUpAppointmentScheduled: injuryReq.FollowUpAppointmentScheduled,
				FollowUpDetails:            injuryReq.FollowUpDetails,
				PhotosTaken:                injuryReq.PhotosTaken,
				WitnessStatementsTaken:     injuryReq.WitnessStatementsTaken,
			}

			if err := tx.Create(&injury).Error; err != nil {
				tx.Rollback()
				h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create injury", err)
				return
			}
		}
	}

	// Create witnesses
	for _, witnessReq := range req.Witnesses {
		witness := models.IncidentWitness{
			IncidentReportID:   incidentReport.ID,
			WitnessType:       witnessReq.WitnessType,
			StaffUserID:       witnessReq.StaffUserID,
			ParticipantID:     witnessReq.ParticipantID,
			FirstName:         witnessReq.FirstName,
			LastName:          witnessReq.LastName,
			ContactPhone:      witnessReq.ContactPhone,
			ContactEmail:      witnessReq.ContactEmail,
			Relationship:      witnessReq.Relationship,
			WitnessStatement:  witnessReq.WitnessStatement,
			StatementTakenBy:  &userID,
			StatementDate:     &now,
			StatementMethod:   witnessReq.StatementMethod,
			WitnessCredibility: witnessReq.WitnessCredibility,
			Notes:             witnessReq.Notes,
		}

		if err := tx.Create(&witness).Error; err != nil {
			tx.Rollback()
			h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create witness", err)
			return
		}
	}

	// TODO: Create timeline entries (models need to be implemented)
	// for _, timelineReq := range req.Timeline {
	// 	timeline := models.IncidentTimeline{...}
	// 	...
	// }

	// TODO: Create risk factors (models need to be implemented)
	// for _, riskReq := range req.RiskFactors {
	// 	risk := models.IncidentRiskFactor{...}
	// 	...
	// }

	// TODO: Create follow-up actions (models need to be implemented)
	// for _, actionReq := range req.FollowUpActions {
	// 	action := models.IncidentFollowUpAction{...}
	// 	...
	// }

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}

	// Load relationships for response
	h.DB.Preload("Participant").Preload("Reporter").Find(&incidentReport)

	h.SendSuccessResponse(c, incidentReport)
}

// GetIncidentReports retrieves incident reports (with filtering)
func (h *Handler) GetIncidentReports(c *gin.Context) {
	// Get user info from context
	userID := h.GetUserIDFromContext(c)
	userRole := h.GetUserRoleFromContext(c)
	if userID == "" {
		h.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get user to get organization ID
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	// Build query
	query := h.DB.Where("organization_id = ?", user.OrganizationID)

	// Filter based on user role
	switch userRole {
	case "care_worker":
		// Care workers can only see their own reports
		query = query.Where("reported_by = ?", userID)
	case "support_coordinator", "admin", "manager":
		// Support coordinators and admins can see all reports in their organization
	default:
		h.SendErrorResponse(c, http.StatusForbidden, "Insufficient permissions", nil)
		return
	}

	// Apply filters from query parameters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if severity := c.Query("severity"); severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if incidentType := c.Query("incident_type"); incidentType != "" {
		query = query.Where("incident_type = ?", incidentType)
	}
	if participantID := c.Query("participant_id"); participantID != "" {
		query = query.Where("participant_id = ?", participantID)
	}

	// Pagination
	page := 1
	limit := 20
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	offset := (page - 1) * limit

	var incidentReports []models.IncidentReport
	var total int64

	// Get total count
	query.Model(&models.IncidentReport{}).Count(&total)

	// Get paginated results with relationships
	if err := query.
		Preload("Participant").
		Preload("Reporter").
		Preload("Reviewer").
		Preload("Documents").
		Preload("PeopleInvolved").
		Preload("PeopleInvolved.Participant").
		Preload("PeopleInvolved.StaffUser").
		Preload("Witnesses").
		Preload("Timeline").
		Preload("RiskFactors").
		Preload("FollowUpActions").
		Preload("Notifications").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&incidentReports).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch incident reports", err)
		return
	}

	response := gin.H{
		"incident_reports": incidentReports,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	h.SendSuccessResponse(c, response)
}

// GetIncidentReport retrieves a single incident report
func (h *Handler) GetIncidentReport(c *gin.Context) {
	id := c.Param("id")
	userID := h.GetUserIDFromContext(c)
	userRole := h.GetUserRoleFromContext(c)

	if userID == "" {
		h.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get user to get organization ID
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	var incidentReport models.IncidentReport
	query := h.DB.Where("id = ? AND organization_id = ?", id, user.OrganizationID)

	// Check access based on role
	if userRole == "care_worker" {
		query = query.Where("reported_by = ?", userID)
	}

	if err := query.
		Preload("Participant").
		Preload("Reporter").
		Preload("Reviewer").
		Preload("Documents").
		Preload("PeopleInvolved").
		Preload("PeopleInvolved.Participant").
		Preload("PeopleInvolved.StaffUser").
		Preload("PeopleInvolved.Injuries").
		Preload("Witnesses").
		Preload("Witnesses.StaffUser").
		Preload("Witnesses.Participant").
		Preload("Timeline").
		Preload("Timeline.RecordedByUser").
		Preload("RiskFactors").
		Preload("RiskFactors.ActionAssignedToUser").
		Preload("FollowUpActions").
		Preload("FollowUpActions.AssignedToUser").
		Preload("FollowUpActions.AssignedByUser").
		Preload("Notifications").
		Preload("Notifications.NotificationSentByUser").
		First(&incidentReport).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.SendErrorResponse(c, http.StatusNotFound, "Incident report not found", err)
		} else {
			h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch incident report", err)
		}
		return
	}

	h.SendSuccessResponse(c, incidentReport)
}

// UpdateIncidentReport updates an incident report (support coordinator only)
func (h *Handler) UpdateIncidentReport(c *gin.Context) {
	id := c.Param("id")
	userID := h.GetUserIDFromContext(c)
	userRole := h.GetUserRoleFromContext(c)

	// Only support coordinators and admins can update incident reports
	if userRole != "support_coordinator" && userRole != "admin" && userRole != "manager" {
		h.SendErrorResponse(c, http.StatusForbidden, "Only support coordinators can update incident reports", nil)
		return
	}

	var req IncidentReportUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// Get user to get organization ID
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	var incidentReport models.IncidentReport
	if err := h.DB.First(&incidentReport, "id = ? AND organization_id = ?", id, user.OrganizationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.SendErrorResponse(c, http.StatusNotFound, "Incident report not found", err)
		} else {
			h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch incident report", err)
		}
		return
	}

	// Update fields
	updates := map[string]interface{}{}
	
	if req.Status != "" {
		updates["status"] = req.Status
		if req.Status == "completed" {
			now := time.Now()
			updates["completed_at"] = &now
		}
	}
	
	if req.Priority != "" {
		updates["priority"] = req.Priority
	}
	
	if req.ReviewNotes != nil {
		updates["review_notes"] = *req.ReviewNotes
	}
	
	if req.NDISNotified {
		updates["ndis_notified"] = req.NDISNotified
		if incidentReport.NDISNotifiedAt == nil {
			now := time.Now()
			updates["ndis_notified_at"] = &now
		}
	}
	
	if req.NDISReference != nil {
		updates["ndis_reference"] = *req.NDISReference
	}
	
	if req.PreventiveMeasures != nil {
		updates["preventive_measures"] = *req.PreventiveMeasures
	}
	
	updates["follow_up_required"] = req.FollowUpRequired
	if req.FollowUpDetails != nil {
		updates["follow_up_details"] = *req.FollowUpDetails
	}

	// Set reviewer info
	updates["reviewed_by"] = userID
	now := time.Now()
	updates["reviewed_at"] = &now

	if err := h.DB.Model(&incidentReport).Updates(updates).Error; err != nil {
		h.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update incident report", err)
		return
	}

	// Reload with relationships
	h.DB.Preload("Participant").Preload("Reporter").Preload("Reviewer").Preload("Documents").First(&incidentReport)

	h.SendSuccessResponse(c, incidentReport)
}

// GetIncidentReportStats gets statistics for incident reports dashboard
func (h *Handler) GetIncidentReportStats(c *gin.Context) {
	userID := h.GetUserIDFromContext(c)
	if userID == "" {
		h.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get user to get organization ID
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		h.SendErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	// Base query for organization
	baseQuery := h.DB.Model(&models.IncidentReport{}).Where("organization_id = ?", user.OrganizationID)

	// Total reports
	var totalReports int64
	baseQuery.Count(&totalReports)

	// Reports by status
	var statusStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	baseQuery.Select("status, COUNT(*) as count").Group("status").Find(&statusStats)

	// Reports by severity
	var severityStats []struct {
		Severity string `json:"severity"`
		Count    int64  `json:"count"`
	}
	baseQuery.Select("severity, COUNT(*) as count").Group("severity").Find(&severityStats)

	// Reports by incident type
	var typeStats []struct {
		IncidentType string `json:"incident_type"`
		Count        int64  `json:"count"`
	}
	baseQuery.Select("incident_type, COUNT(*) as count").Group("incident_type").Find(&typeStats)

	// Recent reports (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var recentReports int64
	baseQuery.Where("created_at >= ?", thirtyDaysAgo).Count(&recentReports)

	// NDIS notifications required
	var ndisNotificationsRequired int64
	baseQuery.Where("ndis_notified = ? AND (severity = ? OR severity = ?)", false, "high", "critical").Count(&ndisNotificationsRequired)

	response := gin.H{
		"total_reports":                totalReports,
		"recent_reports":               recentReports,
		"ndis_notifications_required":  ndisNotificationsRequired,
		"status_breakdown":             statusStats,
		"severity_breakdown":           severityStats,
		"incident_type_breakdown":      typeStats,
	}

	h.SendSuccessResponse(c, response)
}