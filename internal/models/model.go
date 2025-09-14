package models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organization represents the care provider organization
type Organization struct {
	ID        string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null"`
	ABN       string         `json:"abn" gorm:"type:varchar(11);unique"`
	Phone     string         `json:"phone" gorm:"type:varchar(20)"`
	Email     string         `json:"email" gorm:"type:varchar(255)"`
	Website   string         `json:"website" gorm:"type:varchar(255)"`
	Address   Address        `json:"address" gorm:"embedded;embeddedPrefix:address_"`
	NDISReg   NDISReg        `json:"ndis_registration" gorm:"embedded;embeddedPrefix:ndis_"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Users        []User        `json:"users,omitempty" gorm:"foreignKey:OrganizationID"`
	Participants []Participant `json:"participants,omitempty" gorm:"foreignKey:OrganizationID"`
}

// User represents system users (staff, admins, etc.)
type User struct {
	ID             string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	Email          string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash   string         `json:"-" gorm:"type:varchar(255);not null"`
	FirstName      string         `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName       string         `json:"last_name" gorm:"type:varchar(100);not null"`
	Phone          string         `json:"phone" gorm:"type:varchar(20)"`
	Role           string         `json:"role" gorm:"type:varchar(50);not null;index"`     // admin, manager, care_worker, support_coordinator
	RoleID         *string        `json:"role_id,omitempty" gorm:"type:varchar(36);index"` // New role-based system
	Timezone       string         `json:"timezone" gorm:"type:varchar(100);default:'Australia/Adelaide'"` // User's preferred timezone
	OrganizationID string         `json:"organization_id" gorm:"type:varchar(36);not null;index"`
	IsActive       bool           `json:"is_active" gorm:"default:true;index"`
	LastLoginAt    *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Organization    Organization     `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	Shifts          []Shift          `json:"shifts,omitempty" gorm:"foreignKey:StaffID"`
	UploadedDocs    []Document       `json:"uploaded_documents,omitempty" gorm:"foreignKey:UploadedBy"`
	UserPermissions []UserPermission `json:"permissions,omitempty" gorm:"foreignKey:UserID"`
	RefreshTokens   []RefreshToken   `json:"-" gorm:"foreignKey:UserID"`
}

// Participant represents care recipients
type Participant struct {
	ID             string             `json:"id" gorm:"type:varchar(36);primaryKey"`
	FirstName      string             `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName       string             `json:"last_name" gorm:"type:varchar(100);not null"`
	DateOfBirth    time.Time          `json:"date_of_birth" gorm:"not null;index"`
	NDISNumber     string             `json:"ndis_number" gorm:"type:varchar(10);uniqueIndex"`
	Email          string             `json:"email" gorm:"type:varchar(255)"`
	Phone          string             `json:"phone" gorm:"type:varchar(20)"`
	Address        Address            `json:"address" gorm:"embedded;embeddedPrefix:address_"`
	MedicalInfo    MedicalInformation `json:"medical_information" gorm:"embedded;embeddedPrefix:medical_"`
	Funding        FundingInformation `json:"funding" gorm:"embedded;embeddedPrefix:funding_"`
	OrganizationID string             `json:"organization_id" gorm:"type:varchar(36);not null;index"`
	IsActive       bool               `json:"is_active" gorm:"default:true;index"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	DeletedAt      gorm.DeletedAt     `json:"-" gorm:"index"`

	// Relationships
	Organization      Organization       `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	EmergencyContacts []EmergencyContact `json:"emergency_contacts,omitempty" gorm:"foreignKey:ParticipantID"`
	Shifts            []Shift            `json:"shifts,omitempty" gorm:"foreignKey:ParticipantID"`
	Documents         []Document         `json:"documents,omitempty" gorm:"foreignKey:ParticipantID"`
	CarePlans         []CarePlan         `json:"care_plans,omitempty" gorm:"foreignKey:ParticipantID"`
}

// EmergencyContact represents participant emergency contacts
type EmergencyContact struct {
	ID            string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	ParticipantID string    `json:"participant_id" gorm:"type:varchar(36);not null;index"`
	Name          string    `json:"name" gorm:"type:varchar(200);not null"`
	Relationship  string    `json:"relationship" gorm:"type:varchar(50);not null"`
	Phone         string    `json:"phone" gorm:"type:varchar(20);not null"`
	Email         string    `json:"email" gorm:"type:varchar(255)"`
	IsPrimary     bool      `json:"is_primary" gorm:"default:false"`
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relationships
	Participant Participant `json:"participant,omitempty" gorm:"foreignKey:ParticipantID"`
}

// Shift represents scheduled work shifts
type Shift struct {
	ID              string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	ParticipantID   string         `json:"participant_id" gorm:"type:varchar(36);not null;index"`
	StaffID         string         `json:"staff_id" gorm:"type:varchar(36);not null;index"`
	StartTime       time.Time      `json:"start_time" gorm:"not null;index"`
	EndTime         time.Time      `json:"end_time" gorm:"not null;index"`
	ActualStartTime *time.Time     `json:"actual_start_time,omitempty"`
	ActualEndTime   *time.Time     `json:"actual_end_time,omitempty"`
	ServiceType     string         `json:"service_type" gorm:"type:varchar(100);not null;index"`
	Location        string         `json:"location" gorm:"type:varchar(100);not null"`
	Status          string         `json:"status" gorm:"type:varchar(50);default:'scheduled';index"` // scheduled, in_progress, completed, cancelled, no_show
	HourlyRate      float64        `json:"hourly_rate" gorm:"type:decimal(10,2);not null"`
	TotalCost       float64        `json:"total_cost" gorm:"type:decimal(10,2)"`
	Notes           string         `json:"notes" gorm:"type:text"`
	CompletionNotes string         `json:"completion_notes" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Participant Participant `json:"participant,omitempty" gorm:"foreignKey:ParticipantID"`
	Staff       User        `json:"staff,omitempty" gorm:"foreignKey:StaffID"`
}

// Document represents uploaded files and documents
type Document struct {
	ID               string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	ParticipantID    *string        `json:"participant_id,omitempty" gorm:"type:varchar(36);index"`
	IncidentReportID *string        `json:"incident_report_id,omitempty" gorm:"type:varchar(36);index"`
	UploadedBy       string         `json:"uploaded_by" gorm:"type:varchar(36);not null;index"`
	Filename         string         `json:"filename" gorm:"type:varchar(255);not null"`
	OriginalFilename string         `json:"original_filename" gorm:"type:varchar(255);not null"`
	Title            string         `json:"title" gorm:"type:varchar(255);not null"`
	Description      string         `json:"description" gorm:"type:text"`
	Category         string         `json:"category" gorm:"type:varchar(100);not null;index"` // care_plan, medical_record, incident_report, assessment, etc.
	FileType         string         `json:"file_type" gorm:"type:varchar(100);not null"`
	FileSize         int64          `json:"file_size" gorm:"not null"`
	FilePath         string         `json:"file_path" gorm:"type:varchar(500);not null"`
	URL              string         `json:"url" gorm:"type:varchar(500)"`
	IsActive         bool           `json:"is_active" gorm:"default:true;index"`
	ExpiryDate       *time.Time     `json:"expiry_date,omitempty" gorm:"index"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Participant    *Participant    `json:"participant,omitempty" gorm:"foreignKey:ParticipantID"`
	IncidentReport *IncidentReport `json:"incident_report,omitempty" gorm:"foreignKey:IncidentReportID"`
	UploadedByUser User            `json:"uploaded_by_user,omitempty" gorm:"foreignKey:UploadedBy"`
}

// CarePlan represents participant care plans
type CarePlan struct {
	ID            string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	ParticipantID string         `json:"participant_id" gorm:"type:varchar(36);not null;index"`
	Title         string         `json:"title" gorm:"type:varchar(255);not null"`
	Description   string         `json:"description" gorm:"type:text"`
	Goals         string         `json:"goals" gorm:"type:text"` // JSON string of goals
	StartDate     time.Time      `json:"start_date" gorm:"not null"`
	EndDate       *time.Time     `json:"end_date,omitempty"`
	Status        string         `json:"status" gorm:"type:varchar(50);default:'active';index"` // active, completed, cancelled
	CreatedBy     string         `json:"created_by" gorm:"type:varchar(36);not null"`
	ApprovedBy    *string        `json:"approved_by,omitempty" gorm:"type:varchar(36)"`
	ApprovedAt    *time.Time     `json:"approved_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Participant Participant `json:"participant,omitempty" gorm:"foreignKey:ParticipantID"`
	Creator     User        `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Approver    *User       `json:"approver,omitempty" gorm:"foreignKey:ApprovedBy"`
}

// IncidentReport represents incident reports submitted by care workers
type IncidentReport struct {
	ID                  string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	ParticipantID       string         `json:"participant_id" gorm:"type:varchar(36);not null;index"`
	ReportedBy          string         `json:"reported_by" gorm:"type:varchar(36);not null;index"` // Care worker who reported
	ReviewedBy          *string        `json:"reviewed_by,omitempty" gorm:"type:varchar(36);index"` // Support coordinator who reviewed
	OrganizationID      string         `json:"organization_id" gorm:"type:varchar(36);not null;index"`
	
	// Incident details
	IncidentDate        time.Time      `json:"incident_date" gorm:"not null;index"`
	IncidentTime        string         `json:"incident_time" gorm:"type:varchar(10);not null"` // HH:MM format
	Location            string         `json:"location" gorm:"type:varchar(255);not null"`
	IncidentType        string         `json:"incident_type" gorm:"type:varchar(100);not null;index"` // injury, medication_error, behavioral, property_damage, etc.
	Severity            string         `json:"severity" gorm:"type:varchar(50);not null;index"` // low, medium, high, critical
	
	// Description and details
	Description         string         `json:"description" gorm:"type:text;not null"` // What happened
	ImmediateAction     string         `json:"immediate_action" gorm:"type:text"` // What was done immediately
	InjuriesDescription *string        `json:"injuries_description,omitempty" gorm:"type:text"` // If injuries occurred
	WitnessesPresent    bool           `json:"witnesses_present" gorm:"default:false"`
	WitnessDetails      *string        `json:"witness_details,omitempty" gorm:"type:text"`
	
	// Medical and emergency response
	MedicalAttention    bool           `json:"medical_attention" gorm:"default:false"`
	MedicalDetails      *string        `json:"medical_details,omitempty" gorm:"type:text"`
	EmergencyServices   bool           `json:"emergency_services" gorm:"default:false"`
	EmergencyDetails    *string        `json:"emergency_details,omitempty" gorm:"type:text"`
	
	// Notifications
	FamilyNotified      bool           `json:"family_notified" gorm:"default:false"`
	FamilyNotifiedAt    *time.Time     `json:"family_notified_at,omitempty"`
	NDISNotified        bool           `json:"ndis_notified" gorm:"default:false"`
	NDISNotifiedAt      *time.Time     `json:"ndis_notified_at,omitempty"`
	NDISReference       *string        `json:"ndis_reference,omitempty" gorm:"type:varchar(100)"` // NDIS incident reference number
	
	// Follow-up and prevention
	PreventiveMeasures  *string        `json:"preventive_measures,omitempty" gorm:"type:text"`
	FollowUpRequired    bool           `json:"follow_up_required" gorm:"default:false"`
	FollowUpDetails     *string        `json:"follow_up_details,omitempty" gorm:"type:text"`
	
	// Status and workflow
	Status              string         `json:"status" gorm:"type:varchar(50);default:'submitted';index"` // submitted, under_review, completed, requires_action
	Priority            string         `json:"priority" gorm:"type:varchar(50);default:'medium';index"` // low, medium, high, urgent
	ReviewNotes         *string        `json:"review_notes,omitempty" gorm:"type:text"` // Support coordinator notes
	ReviewedAt          *time.Time     `json:"reviewed_at,omitempty"`
	CompletedAt         *time.Time     `json:"completed_at,omitempty"`
	
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Participant   Participant `json:"participant,omitempty" gorm:"foreignKey:ParticipantID"`
	Reporter      User        `json:"reporter,omitempty" gorm:"foreignKey:ReportedBy"`
	Reviewer      *User       `json:"reviewer,omitempty" gorm:"foreignKey:ReviewedBy"`
	Organization  Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	Documents     []Document  `json:"documents,omitempty" gorm:"foreignKey:IncidentReportID"`
}

// BeforeCreate hook for IncidentReport
func (i *IncidentReport) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return
}

// RefreshToken stores JWT refresh tokens
type RefreshToken struct {
	ID        string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	UserID    string    `json:"user_id" gorm:"type:varchar(36);not null;index"`
	Token     string    `json:"token" gorm:"type:varchar(255);not null;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	IsRevoked bool      `json:"is_revoked" gorm:"default:false;index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UserPermission represents user permissions
type UserPermission struct {
	ID         string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	UserID     string    `json:"user_id" gorm:"type:varchar(36);not null;index"`
	Permission string    `json:"permission" gorm:"type:varchar(100);not null"` // read_participants, create_shifts, etc.
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Embedded structs for common data patterns

// Address represents physical addresses
type Address struct {
	Street   string `json:"street" gorm:"type:varchar(255)"`
	Suburb   string `json:"suburb" gorm:"type:varchar(100)"`
	State    string `json:"state" gorm:"type:varchar(50)"`
	Postcode string `json:"postcode" gorm:"type:varchar(10)"`
	Country  string `json:"country" gorm:"type:varchar(100);default:'Australia'"`
}

// NDISReg represents NDIS registration information
type NDISReg struct {
	RegistrationNumber string     `json:"registration_number" gorm:"type:varchar(50)"`
	RegistrationStatus string     `json:"registration_status" gorm:"type:varchar(50);default:'active'"`
	ExpiryDate         *time.Time `json:"expiry_date,omitempty"`
}

// MedicalInformation represents participant medical details
type MedicalInformation struct {
	Conditions  string `json:"conditions" gorm:"type:text"`  // JSON array of conditions
	Medications string `json:"medications" gorm:"type:text"` // JSON array of medications
	Allergies   string `json:"allergies" gorm:"type:text"`   // JSON array of allergies
	DoctorName  string `json:"doctor_name" gorm:"type:varchar(255)"`
	DoctorPhone string `json:"doctor_phone" gorm:"type:varchar(20)"`
	Notes       string `json:"notes" gorm:"type:text"`
}

// FundingInformation represents NDIS funding details
type FundingInformation struct {
	TotalBudget     float64    `json:"total_budget" gorm:"type:decimal(12,2);default:0"`
	UsedBudget      float64    `json:"used_budget" gorm:"type:decimal(12,2);default:0"`
	RemainingBudget float64    `json:"remaining_budget" gorm:"type:decimal(12,2);default:0"`
	BudgetYear      string     `json:"budget_year" gorm:"type:varchar(20)"` // e.g., "2025-2026"
	PlanStartDate   *time.Time `json:"plan_start_date,omitempty"`
	PlanEndDate     *time.Time `json:"plan_end_date,omitempty"`
}

// BeforeCreate hooks for generating UUIDs
func (o *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}

func (p *Participant) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}

func (e *EmergencyContact) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return
}

func (s *Shift) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	// Calculate total cost based on duration and hourly rate
	duration := s.EndTime.Sub(s.StartTime).Hours()
	s.TotalCost = duration * s.HourlyRate
	return
}

func (d *Document) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return
}

func (c *CarePlan) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return
}

func (p *UserPermission) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}

// BeforeUpdate hooks for maintaining data consistency
func (s *Shift) BeforeUpdate(tx *gorm.DB) (err error) {
	// Recalculate total cost if times have changed
	if s.EndTime.After(s.StartTime) {
		duration := s.EndTime.Sub(s.StartTime).Hours()
		s.TotalCost = duration * s.HourlyRate
	}
	return
}

func (p *Participant) BeforeUpdate(tx *gorm.DB) (err error) {
	// Recalculate remaining budget
	p.Funding.RemainingBudget = p.Funding.TotalBudget - p.Funding.UsedBudget
	return
}

// Enhanced Incident Models for comprehensive incident reporting

// IncidentPersonInvolved represents people involved in incidents (patients, staff, visitors, etc.)
type IncidentPersonInvolved struct {
	ID                     string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	IncidentReportID       string         `json:"incident_report_id" gorm:"type:varchar(36);not null;index"`
	
	// Person identification
	PersonType             string         `json:"person_type" gorm:"type:varchar(50);not null;index"` // participant, staff, visitor, family_member, contractor, other
	ParticipantID          *string        `json:"participant_id,omitempty" gorm:"type:varchar(36);index"`
	StaffUserID            *string        `json:"staff_user_id,omitempty" gorm:"type:varchar(36);index"`
	
	// For non-participants/non-staff
	FirstName              *string        `json:"first_name,omitempty" gorm:"type:varchar(100)"`
	LastName               *string        `json:"last_name,omitempty" gorm:"type:varchar(100)"`
	Age                    *int           `json:"age,omitempty"`
	Gender                 *string        `json:"gender,omitempty" gorm:"type:varchar(20)"`
	ContactPhone           *string        `json:"contact_phone,omitempty" gorm:"type:varchar(20)"`
	ContactEmail           *string        `json:"contact_email,omitempty" gorm:"type:varchar(255)"`
	RelationshipToParticipant *string     `json:"relationship_to_participant,omitempty" gorm:"type:varchar(100)"`
	EmployerOrganization   *string        `json:"employer_organization,omitempty" gorm:"type:varchar(255)"`
	
	// Involvement details
	RoleInIncident         string         `json:"role_in_incident" gorm:"type:varchar(100);not null"` // injured_party, witness, responding_staff, bystander, instigator
	WasInjured             bool           `json:"was_injured" gorm:"default:false;index"`
	InjuryDescription      *string        `json:"injury_description,omitempty" gorm:"type:text"`
	InjurySeverity         *string        `json:"injury_severity,omitempty" gorm:"type:varchar(50)"` // minor, moderate, serious, critical
	
	// Medical response for this person
	MedicalAttentionRequired bool           `json:"medical_attention_required" gorm:"default:false"`
	MedicalAttentionProvided bool           `json:"medical_attention_provided" gorm:"default:false"`
	MedicalProvider          *string        `json:"medical_provider,omitempty" gorm:"type:varchar(255)"`
	TransportedToHospital    bool           `json:"transported_to_hospital" gorm:"default:false"`
	HospitalName             *string        `json:"hospital_name,omitempty" gorm:"type:varchar(255)"`
	AmbulanceCalled          bool           `json:"ambulance_called" gorm:"default:false"`
	
	// Actions taken
	ImmediateCareProvided    *string        `json:"immediate_care_provided,omitempty" gorm:"type:text"`
	OngoingCareRequired      bool           `json:"ongoing_care_required" gorm:"default:false"`
	OngoingCareDetails       *string        `json:"ongoing_care_details,omitempty" gorm:"type:text"`
	
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`

	// Relationships
	IncidentReport   IncidentReport `json:"incident_report,omitempty" gorm:"foreignKey:IncidentReportID"`
	Participant      *Participant   `json:"participant,omitempty" gorm:"foreignKey:ParticipantID"`
	StaffUser        *User          `json:"staff_user,omitempty" gorm:"foreignKey:StaffUserID"`
	Injuries         []IncidentInjury `json:"injuries,omitempty" gorm:"foreignKey:PersonInvolvedID"`
}

// IncidentInjury represents detailed injury tracking
type IncidentInjury struct {
	ID                     string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	IncidentReportID       string         `json:"incident_report_id" gorm:"type:varchar(36);not null;index"`
	PersonInvolvedID       string         `json:"person_involved_id" gorm:"type:varchar(36);not null;index"`
	
	// Injury classification
	InjuryType             string         `json:"injury_type" gorm:"type:varchar(100);not null;index"` // laceration, bruise, fracture, burn, sprain, strain, head_injury, psychological, other
	InjuryCategory         string         `json:"injury_category" gorm:"type:varchar(50);not null"` // physical, psychological, emotional, behavioral
	
	// Body location
	BodyPart               string         `json:"body_part" gorm:"type:varchar(100);not null"` // head, face, neck, chest, back, arms, hands, legs, feet, internal, multiple
	BodySide               *string        `json:"body_side,omitempty" gorm:"type:varchar(20)"` // left, right, bilateral, center
	SpecificLocation       *string        `json:"specific_location,omitempty" gorm:"type:text"`
	
	// Injury details
	SeverityLevel          string         `json:"severity_level" gorm:"type:varchar(50);not null;index"` // superficial, minor, moderate, serious, life_threatening
	SizeDimensions         *string        `json:"size_dimensions,omitempty" gorm:"type:varchar(100)"`
	DepthDescription       *string        `json:"depth_description,omitempty" gorm:"type:varchar(100)"` // surface, shallow, deep, penetrating
	
	// Cause and mechanism
	CauseOfInjury          string         `json:"cause_of_injury" gorm:"type:varchar(100);not null"` // fall, struck_by_object, collision, chemical, electrical, behavioral, self_harm, other
	MechanismDescription   *string        `json:"mechanism_description,omitempty" gorm:"type:text"`
	
	// Treatment
	FirstAidGiven          bool           `json:"first_aid_given" gorm:"default:false"`
	FirstAidDetails        *string        `json:"first_aid_details,omitempty" gorm:"type:text"`
	MedicalTreatmentRequired bool         `json:"medical_treatment_required" gorm:"default:false"`
	TreatmentProvided      *string        `json:"treatment_provided,omitempty" gorm:"type:text"`
	OngoingTreatmentRequired bool         `json:"ongoing_treatment_required" gorm:"default:false"`
	
	// Follow-up
	FollowUpAppointmentRequired bool       `json:"follow_up_appointment_required" gorm:"default:false"`
	FollowUpAppointmentScheduled bool      `json:"follow_up_appointment_scheduled" gorm:"default:false"`
	FollowUpDetails        *string        `json:"follow_up_details,omitempty" gorm:"type:text"`
	
	// Documentation
	PhotosTaken            bool           `json:"photos_taken" gorm:"default:false"`
	WitnessStatementsTaken bool           `json:"witness_statements_taken" gorm:"default:false"`
	
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`

	// Relationships
	IncidentReport   IncidentReport         `json:"incident_report,omitempty" gorm:"foreignKey:IncidentReportID"`
	PersonInvolved   IncidentPersonInvolved `json:"person_involved,omitempty" gorm:"foreignKey:PersonInvolvedID"`
}

// IncidentWitness represents detailed witness information
type IncidentWitness struct {
	ID                     string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	IncidentReportID       string         `json:"incident_report_id" gorm:"type:varchar(36);not null;index"`
	
	// Witness identification
	WitnessType            string         `json:"witness_type" gorm:"type:varchar(50);not null;index"` // staff, participant, family_member, visitor, contractor, other
	StaffUserID            *string        `json:"staff_user_id,omitempty" gorm:"type:varchar(36);index"`
	ParticipantID          *string        `json:"participant_id,omitempty" gorm:"type:varchar(36);index"`
	
	// For external witnesses
	FirstName              *string        `json:"first_name,omitempty" gorm:"type:varchar(100)"`
	LastName               *string        `json:"last_name,omitempty" gorm:"type:varchar(100)"`
	ContactPhone           *string        `json:"contact_phone,omitempty" gorm:"type:varchar(20)"`
	ContactEmail           *string        `json:"contact_email,omitempty" gorm:"type:varchar(255)"`
	Relationship           *string        `json:"relationship,omitempty" gorm:"type:varchar(100)"`
	
	// Witness account
	WitnessStatement       *string        `json:"witness_statement,omitempty" gorm:"type:text"`
	StatementTakenBy       *string        `json:"statement_taken_by,omitempty" gorm:"type:varchar(36)"`
	StatementDate          *time.Time     `json:"statement_date,omitempty"`
	StatementMethod        *string        `json:"statement_method,omitempty" gorm:"type:varchar(50)"` // verbal, written, recorded, signed
	
	// Witness reliability
	WitnessCredibility     *string        `json:"witness_credibility,omitempty" gorm:"type:varchar(50)"` // high, medium, low, unknown
	Notes                  *string        `json:"notes,omitempty" gorm:"type:text"`
	
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`

	// Relationships
	IncidentReport   IncidentReport `json:"incident_report,omitempty" gorm:"foreignKey:IncidentReportID"`
	StaffUser        *User          `json:"staff_user,omitempty" gorm:"foreignKey:StaffUserID"`
	Participant      *Participant   `json:"participant,omitempty" gorm:"foreignKey:ParticipantID"`
	StatementTaker   *User          `json:"statement_taker,omitempty" gorm:"foreignKey:StatementTakenBy"`
}

// IncidentDocument represents enhanced document management for incidents
type IncidentDocument struct {
	ID                     string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	IncidentReportID       string         `json:"incident_report_id" gorm:"type:varchar(36);not null;index"`
	
	// Document classification
	DocumentType           string         `json:"document_type" gorm:"type:varchar(100);not null;index"` // photo, witness_statement, medical_report, cctv_footage, floor_plan, policy_document, correspondence, investigation_report, other
	DocumentCategory       string         `json:"document_category" gorm:"type:varchar(50);not null"` // evidence, medical, administrative, legal, follow_up
	
	// File information
	OriginalFilename       string         `json:"original_filename" gorm:"type:varchar(255);not null"`
	StoredFilename         string         `json:"stored_filename" gorm:"type:varchar(255);not null"`
	FilePath               string         `json:"file_path" gorm:"type:varchar(500);not null"`
	FileSizeBytes          *int           `json:"file_size_bytes,omitempty"`
	MimeType               *string        `json:"mime_type,omitempty" gorm:"type:varchar(100)"`
	FileHash               *string        `json:"file_hash,omitempty" gorm:"type:varchar(64)"`
	
	// Document metadata
	Title                  string         `json:"title" gorm:"type:varchar(255);not null"`
	Description            *string        `json:"description,omitempty" gorm:"type:text"`
	DateCreated            *time.Time     `json:"date_created,omitempty"`
	DateUploaded           time.Time      `json:"date_uploaded"`
	UploadedBy             string         `json:"uploaded_by" gorm:"type:varchar(36);not null"`
	
	// Access and security
	IsConfidential         bool           `json:"is_confidential" gorm:"default:false"`
	AccessLevel            string         `json:"access_level" gorm:"type:varchar(50);default:'standard'"` // public, standard, restricted, confidential
	RetentionPeriodYears   int            `json:"retention_period_years" gorm:"default:7"`
	DestructionDate        *time.Time     `json:"destruction_date,omitempty"`
	
	// Version control
	VersionNumber          int            `json:"version_number" gorm:"default:1"`
	ParentDocumentID       *string        `json:"parent_document_id,omitempty" gorm:"type:varchar(36)"`
	IsCurrentVersion       bool           `json:"is_current_version" gorm:"default:true"`
	
	// Legal and compliance
	IsLegalEvidence        bool           `json:"is_legal_evidence" gorm:"default:false"`
	ChainOfCustody         *string        `json:"chain_of_custody,omitempty" gorm:"type:text"`
	
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	IncidentReport   IncidentReport     `json:"incident_report,omitempty" gorm:"foreignKey:IncidentReportID"`
	Uploader         User               `json:"uploader,omitempty" gorm:"foreignKey:UploadedBy"`
	ParentDocument   *IncidentDocument  `json:"parent_document,omitempty" gorm:"foreignKey:ParentDocumentID"`
}

// IncidentNotification represents tracking of all notifications sent
type IncidentNotification struct {
	ID                     string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	IncidentReportID       string         `json:"incident_report_id" gorm:"type:varchar(36);not null;index"`
	
	// Notification details
	NotificationType       string         `json:"notification_type" gorm:"type:varchar(50);not null;index"` // family, ndis, management, police, insurance, regulatory, other
	RecipientName          string         `json:"recipient_name" gorm:"type:varchar(255);not null"`
	RecipientContact       *string        `json:"recipient_contact,omitempty" gorm:"type:text"`
	RecipientOrganization  *string        `json:"recipient_organization,omitempty" gorm:"type:varchar(255)"`
	
	// Notification method and timing
	Method                 string         `json:"method" gorm:"type:varchar(50);not null"` // phone, email, letter, in_person, fax, online_portal
	NotificationSentAt     time.Time      `json:"notification_sent_at"`
	NotificationSentBy     string         `json:"notification_sent_by" gorm:"type:varchar(36);not null"`
	
	// Response tracking
	AcknowledgmentRequired bool           `json:"acknowledgment_required" gorm:"default:false"`
	AcknowledgmentReceived bool           `json:"acknowledgment_received" gorm:"default:false"`
	AcknowledgmentReceivedAt *time.Time   `json:"acknowledgment_received_at,omitempty"`
	ResponseReference      *string        `json:"response_reference,omitempty" gorm:"type:varchar(100)"`
	ResponseDetails        *string        `json:"response_details,omitempty" gorm:"type:text"`
	
	// Follow-up requirements
	FollowUpRequired       bool           `json:"follow_up_required" gorm:"default:false"`
	FollowUpDueDate        *time.Time     `json:"follow_up_due_date,omitempty"`
	FollowUpCompleted      bool           `json:"follow_up_completed" gorm:"default:false"`
	FollowUpNotes          *string        `json:"follow_up_notes,omitempty" gorm:"type:text"`
	
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`

	// Relationships
	IncidentReport IncidentReport `json:"incident_report,omitempty" gorm:"foreignKey:IncidentReportID"`
	SentBy         User           `json:"sent_by,omitempty" gorm:"foreignKey:NotificationSentBy"`
}

// CareNote represents care notes taken by staff during shifts
type CareNote struct {
	ID            string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	ParticipantID string         `json:"participant_id" gorm:"type:varchar(36);not null;index"`
	StaffID       string         `json:"staff_id" gorm:"type:varchar(36);not null;index"`
	ShiftID       *string        `json:"shift_id,omitempty" gorm:"type:varchar(36);index"` // Optional link to specific shift
	OrganizationID string        `json:"organization_id" gorm:"type:varchar(36);not null;index"`
	
	// Note content and metadata
	Title         string         `json:"title" gorm:"type:varchar(255);not null"`
	Content       string         `json:"content" gorm:"type:text;not null"`
	NoteType      string         `json:"note_type" gorm:"type:varchar(50);not null;index"` // daily_progress, medication, behaviour, health, communication, achievement, concern
	Priority      string         `json:"priority" gorm:"type:varchar(20);default:'normal';index"` // low, normal, high, urgent
	
	// Timing information
	NoteDate      time.Time      `json:"note_date" gorm:"not null;index"` // Date the observation/event occurred
	NoteTime      string         `json:"note_time" gorm:"type:varchar(10)"` // HH:MM format for specific time if relevant
	
	// Visibility and access
	IsPrivate     bool           `json:"is_private" gorm:"default:false"` // Only visible to admin/management
	IsConfidential bool          `json:"is_confidential" gorm:"default:false"` // Requires special access
	
	// Follow-up tracking
	RequiresFollowUp bool        `json:"requires_follow_up" gorm:"default:false"`
	FollowUpBy     *string       `json:"follow_up_by,omitempty" gorm:"type:varchar(36)"` // User ID who should follow up
	FollowUpDate   *time.Time    `json:"follow_up_date,omitempty"`
	FollowUpStatus string        `json:"follow_up_status" gorm:"type:varchar(50);default:'pending'"` // pending, in_progress, completed, cancelled
	FollowUpNotes  *string       `json:"follow_up_notes,omitempty" gorm:"type:text"`
	
	// Categories and tags
	Tags          *string        `json:"tags,omitempty" gorm:"type:text"` // JSON array of tags for filtering
	Category      *string        `json:"category,omitempty" gorm:"type:varchar(100)"` // Additional categorization
	
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Participant   Participant    `json:"participant,omitempty" gorm:"foreignKey:ParticipantID"`
	Staff         User           `json:"staff,omitempty" gorm:"foreignKey:StaffID"`
	Shift         *Shift         `json:"shift,omitempty" gorm:"foreignKey:ShiftID"`
	Organization  Organization   `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	FollowUpUser  *User          `json:"follow_up_user,omitempty" gorm:"foreignKey:FollowUpBy"`
}

// BeforeCreate hook for CareNote
func (cn *CareNote) BeforeCreate(tx *gorm.DB) (err error) {
	if cn.ID == "" {
		cn.ID = uuid.New().String()
	}
	return
}

// BeforeCreate hooks for enhanced incident models
func (ipi *IncidentPersonInvolved) BeforeCreate(tx *gorm.DB) (err error) {
	if ipi.ID == "" {
		ipi.ID = uuid.New().String()
	}
	return
}

func (ii *IncidentInjury) BeforeCreate(tx *gorm.DB) (err error) {
	if ii.ID == "" {
		ii.ID = uuid.New().String()
	}
	return
}

func (iw *IncidentWitness) BeforeCreate(tx *gorm.DB) (err error) {
	if iw.ID == "" {
		iw.ID = uuid.New().String()
	}
	return
}

func (id *IncidentDocument) BeforeCreate(tx *gorm.DB) (err error) {
	if id.ID == "" {
		id.ID = uuid.New().String()
	}
	return
}

func (in *IncidentNotification) BeforeCreate(tx *gorm.DB) (err error) {
	if in.ID == "" {
		in.ID = uuid.New().String()
	}
	return
}

// Database migration function
func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&Organization{},
		&User{},
		&Participant{},
		&EmergencyContact{},
		&Shift{},
		&Document{},
		&CarePlan{},
		&CareNote{},
		&RefreshToken{},
		&UserPermission{},
		&WorkerAvailability{},
		&WorkerAvailabilityException{},
		&WorkerPreferences{},
		&WorkerSkill{},
		&WorkerLocationPreference{},
		&IncidentReport{},
		&IncidentPersonInvolved{},
		&IncidentInjury{},
		&IncidentWitness{},
		&IncidentDocument{},
		&IncidentNotification{},
	)
}

// Index creation function for better performance
func CreateIndexes(db *gorm.DB) error {
	// Composite indexes for better query performance
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_shifts_participant_date ON shifts(participant_id, start_time)").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_shifts_staff_date ON shifts(staff_id, start_time)").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_documents_participant_category ON documents(participant_id, category)").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_participants_ndis_org ON participants(ndis_number, organization_id)").Error; err != nil {
		return err
	}

	return nil
}

// Sample data seeding function (for development/testing)
func SeedDatabase(db *gorm.DB) error {
	// Create default organization
	org := Organization{
		ID:    "org_default",
		Name:  "DASYIN - ADL Services",
		ABN:   "12345678901",
		Phone: "+61887654321",
		Email: "info@dasyin.com.au",
		Address: Address{
			Street:   "789 Business Ave",
			Suburb:   "Adelaide",
			State:    "SA",
			Postcode: "5000",
			Country:  "Australia",
		},
		NDISReg: NDISReg{
			RegistrationNumber: "REG123456",
			RegistrationStatus: "active",
		},
	}

	if err := db.FirstOrCreate(&org, "id = ?", org.ID).Error; err != nil {
		return err
	}

	// Standard password hash for "password"
	passwordHash := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"

	// Create default admin user
	admin := User{
		ID:             "user_admin",
		Email:          "kennedy@dasyin.com.au",
		FirstName:      "Ken",
		LastName:       "Kinoti",
		Role:           "admin",
		OrganizationID: org.ID,
		IsActive:       true,
		PasswordHash:   passwordHash,
	}

	// Always ensure the protected system admin exists and has correct credentials
	// This is critical for system security and access recovery
	var existingUser User
	err := db.Where("email = ?", admin.Email).First(&existingUser).Error
	if err != nil {
		// User doesn't exist, create it
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
	} else {
		// User exists, ALWAYS update to ensure correct password and role
		// This protects against any manual changes or corruption
		if err := db.Model(&existingUser).Updates(map[string]interface{}{
			"password_hash":   admin.PasswordHash,
			"role":            admin.Role,
			"is_active":       true,
			"first_name":      admin.FirstName,
			"last_name":       admin.LastName,
			"organization_id": admin.OrganizationID,
		}).Error; err != nil {
			return err
		}
	}

	// Add default permissions for admin
	permissions := []string{
		"create_users", "read_users", "update_users", "delete_users",
		"create_participants", "read_participants", "update_participants", "delete_participants",
		"create_shifts", "read_shifts", "update_shifts", "delete_shifts",
		"create_documents", "read_documents", "update_documents", "delete_documents",
		"create_care_plans", "read_care_plans", "update_care_plans", "delete_care_plans",
		"view_reports", "manage_organization",
	}

	for _, perm := range permissions {
		userPerm := UserPermission{
			UserID:     admin.ID,
			Permission: perm,
		}
		db.FirstOrCreate(&userPerm, "user_id = ? AND permission = ?", admin.ID, perm)
	}

	// Create all required test users
	testUsers := []User{
		{
			ID:             "user_super_admin",
			Email:          "admin@dasyin.com.au",
			FirstName:      "Super",
			LastName:       "Admin",
			Role:           "super_admin",
			OrganizationID: org.ID,
			IsActive:       true,
			PasswordHash:   passwordHash,
		},
		{
			ID:             "user_manager",
			Email:          "manager@dasyin.com.au",
			FirstName:      "Team",
			LastName:       "Manager",
			Role:           "manager",
			OrganizationID: org.ID,
			IsActive:       true,
			PasswordHash:   passwordHash,
		},
		{
			ID:             "user_coordinator",
			Email:          "coordinator@dasyin.com.au",
			FirstName:      "Support",
			LastName:       "Coordinator",
			Role:           "support_coordinator",
			OrganizationID: org.ID,
			IsActive:       true,
			PasswordHash:   passwordHash,
		},
		{
			ID:             "user_careworker",
			Email:          "careworker@dasyin.com.au",
			FirstName:      "Care",
			LastName:       "Worker",
			Role:           "care_worker",
			OrganizationID: org.ID,
			IsActive:       true,
			PasswordHash:   passwordHash,
		},
		{
			ID:             "user_org2_admin",
			Email:          "org2admin@dasyin.com.au",
			FirstName:      "Org2",
			LastName:       "Admin",
			Role:           "admin",
			OrganizationID: org.ID,
			IsActive:       true,
			PasswordHash:   passwordHash,
		},
	}

	// Create or update each test user
	for _, user := range testUsers {
		var existingUser User
		err := db.Where("email = ?", user.Email).First(&existingUser).Error
		if err != nil {
			// User doesn't exist, create it
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		} else {
			// User exists, update to ensure correct credentials
			if err := db.Model(&existingUser).Updates(map[string]interface{}{
				"password_hash":   user.PasswordHash,
				"role":            user.Role,
				"is_active":       true,
				"first_name":      user.FirstName,
				"last_name":       user.LastName,
				"organization_id": user.OrganizationID,
			}).Error; err != nil {
				return err
			}
		}

		// Add permissions based on role
		var userPerms []string
		switch user.Role {
		case "super_admin":
			userPerms = permissions // All permissions
		case "admin":
			userPerms = []string{
				"create_users", "read_users", "update_users", "delete_users",
				"create_participants", "read_participants", "update_participants", "delete_participants",
				"create_shifts", "read_shifts", "update_shifts", "delete_shifts",
				"create_documents", "read_documents", "update_documents", "delete_documents",
				"create_care_plans", "read_care_plans", "update_care_plans", "delete_care_plans",
				"view_reports", "manage_organization",
			}
		case "manager":
			userPerms = []string{
				"read_users", "update_users",
				"create_participants", "read_participants", "update_participants",
				"create_shifts", "read_shifts", "update_shifts", "delete_shifts",
				"create_documents", "read_documents", "update_documents",
				"create_care_plans", "read_care_plans", "update_care_plans",
				"view_reports",
			}
		case "support_coordinator":
			userPerms = []string{
				"read_participants", "update_participants",
				"read_shifts", "update_shifts",
				"create_documents", "read_documents", "update_documents",
				"create_care_plans", "read_care_plans", "update_care_plans",
			}
		case "care_worker":
			userPerms = []string{
				"read_participants",
				"read_shifts", "update_shifts",
				"read_documents",
				"read_care_plans",
			}
		}

		for _, perm := range userPerms {
			userPerm := UserPermission{
				UserID:     user.ID,
				Permission: perm,
			}
			db.FirstOrCreate(&userPerm, "user_id = ? AND permission = ?", user.ID, perm)
		}
	}

	// Create default organization settings for the default organization
	if err := SetupOrganizationDefaults(db, org.ID); err != nil {
		log.Printf("Warning: Failed to setup organization defaults: %v", err)
	}

	// Sample data will be created via SQL script
	// Run supabase_schema.sql for participant, care plan, and shift data

	return nil
}
