package models

import (
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
	Role           string         `json:"role" gorm:"type:varchar(50);not null;index"` // admin, manager, care_worker, support_coordinator
	RoleID         *string        `json:"role_id,omitempty" gorm:"type:varchar(36);index"` // New role-based system
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
	Participant    *Participant `json:"participant,omitempty" gorm:"foreignKey:ParticipantID"`
	UploadedByUser User         `json:"uploaded_by_user,omitempty" gorm:"foreignKey:UploadedBy"`
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
		&RefreshToken{},
		&UserPermission{},
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

	// Create default admin user
	admin := User{
		ID:             "user_admin",
		Email:          "kennedy@dasyin.com.au",
		FirstName:      "Ken",
		LastName:       "Kinoti",
		Role:           "super_admin",
		OrganizationID: org.ID,
		IsActive:       true,
		// Note: Password should be hashed in real implementation
		PasswordHash: "$2a$10$n0NvhICRgFPZq/EaeWxW6un3Xrym3.23GJpk4wYchZmpxgETxQani", // "Test123!@#"
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
			"password_hash": admin.PasswordHash,
			"role": admin.Role,
			"is_active": true,
			"first_name": admin.FirstName,
			"last_name": admin.LastName,
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

	return nil
}
