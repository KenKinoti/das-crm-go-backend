package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// JSONB is a custom type for handling PostgreSQL JSONB data
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface for JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return errors.New("cannot scan into JSONB")
	}
}

// WorkerAvailability represents regular weekly availability patterns
type WorkerAvailability struct {
	ID                   string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	UserID               string         `json:"user_id" gorm:"type:varchar(36);not null;index"`
	DayOfWeek            int            `json:"day_of_week" gorm:"type:integer;not null;check:day_of_week >= 0 AND day_of_week <= 6"` // 0 = Sunday, 6 = Saturday
	StartTime            time.Time      `json:"start_time" gorm:"type:time;not null"`
	EndTime              time.Time      `json:"end_time" gorm:"type:time;not null"`
	IsAvailable          bool           `json:"is_available" gorm:"type:boolean;default:true"`
	MaxHoursPerDay       float64        `json:"max_hours_per_day" gorm:"type:decimal(4,2);default:8.00"`
	PreferredHoursPerDay float64        `json:"preferred_hours_per_day" gorm:"type:decimal(4,2);default:6.00"`
	HourlyRate           *float64       `json:"hourly_rate" gorm:"type:decimal(10,2)"`
	BreakDurationMinutes int            `json:"break_duration_minutes" gorm:"type:integer;default:30"`
	TravelTimeMinutes    int            `json:"travel_time_minutes" gorm:"type:integer;default:30"`
	Notes                string         `json:"notes" gorm:"type:text"`
	IsActive             bool           `json:"is_active" gorm:"type:boolean;default:true;index"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for WorkerAvailability
func (WorkerAvailability) TableName() string {
	return "worker_availability"
}

// WorkerAvailabilityException represents date-specific availability overrides
type WorkerAvailabilityException struct {
	ID            string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	UserID        string         `json:"user_id" gorm:"type:varchar(36);not null;index"`
	ExceptionDate time.Time      `json:"exception_date" gorm:"type:date;not null;index"`
	ExceptionType string         `json:"exception_type" gorm:"type:varchar(50);not null;check:exception_type IN ('unavailable', 'limited', 'extended', 'special_rate')"`
	StartTime     *time.Time     `json:"start_time" gorm:"type:time"`
	EndTime       *time.Time     `json:"end_time" gorm:"type:time"`
	MaxHours      *float64       `json:"max_hours" gorm:"type:decimal(4,2)"`
	HourlyRate    *float64       `json:"hourly_rate" gorm:"type:decimal(10,2)"`
	Reason        string         `json:"reason" gorm:"type:varchar(255)"`
	Notes         string         `json:"notes" gorm:"type:text"`
	IsApproved    bool           `json:"is_approved" gorm:"type:boolean;default:false;index"`
	ApprovedBy    *string        `json:"approved_by" gorm:"type:varchar(36)"`
	ApprovedAt    *time.Time     `json:"approved_at" gorm:"type:timestamptz"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User     User  `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Approver *User `json:"approver,omitempty" gorm:"foreignKey:ApprovedBy;constraint:OnDelete:SET NULL"`
}

// TableName specifies the table name for WorkerAvailabilityException
func (WorkerAvailabilityException) TableName() string {
	return "worker_availability_exceptions"
}

// WorkerPreferences represents worker capacity and general preferences
type WorkerPreferences struct {
	ID                       string          `json:"id" gorm:"type:varchar(36);primaryKey"`
	UserID                   string          `json:"user_id" gorm:"type:varchar(36);not null;unique;index"`
	MaxHoursPerWeek          float64         `json:"max_hours_per_week" gorm:"type:decimal(5,2);default:38.00"`
	PreferredHoursPerWeek    float64         `json:"preferred_hours_per_week" gorm:"type:decimal(5,2);default:30.00"`
	MaxConsecutiveDays       int             `json:"max_consecutive_days" gorm:"type:integer;default:5"`
	MinHoursBetweenShifts    int             `json:"min_hours_between_shifts" gorm:"type:integer;default:10"`
	MaxTravelDistanceKm      int             `json:"max_travel_distance_km" gorm:"type:integer;default:50"`
	PreferredShiftTypes      pq.StringArray  `json:"preferred_shift_types" gorm:"type:text[]"`
	WillingWeekendWork       bool            `json:"willing_weekend_work" gorm:"type:boolean;default:true"`
	WillingEveningWork       bool            `json:"willing_evening_work" gorm:"type:boolean;default:true"`
	WillingEarlyMorningWork  bool            `json:"willing_early_morning_work" gorm:"type:boolean;default:true"`
	RequiresOwnTransport     bool            `json:"requires_own_transport" gorm:"type:boolean;default:false"`
	HasOwnVehicle            bool            `json:"has_own_vehicle" gorm:"type:boolean;default:true"`
	SpecialSkills            pq.StringArray  `json:"special_skills" gorm:"type:text[]"`
	ClientPreferences        string          `json:"client_preferences" gorm:"type:text"`
	NotificationPreferences  JSONB           `json:"notification_preferences" gorm:"type:jsonb;default:'{\"email\": true, \"sms\": true, \"app\": true}'"`
	IsActive                 bool            `json:"is_active" gorm:"type:boolean;default:true;index"`
	CreatedAt                time.Time       `json:"created_at"`
	UpdatedAt                time.Time       `json:"updated_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for WorkerPreferences
func (WorkerPreferences) TableName() string {
	return "worker_preferences"
}

// WorkerSkill represents worker skills and qualifications
type WorkerSkill struct {
	ID                  string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	UserID              string         `json:"user_id" gorm:"type:varchar(36);not null;index"`
	SkillCategory       string         `json:"skill_category" gorm:"type:varchar(100);not null;index"`
	SkillName           string         `json:"skill_name" gorm:"type:varchar(200);not null"`
	ProficiencyLevel    string         `json:"proficiency_level" gorm:"type:varchar(50);check:proficiency_level IN ('beginner', 'intermediate', 'advanced', 'expert')"`
	CertificationNumber string         `json:"certification_number" gorm:"type:varchar(100)"`
	ExpiryDate          *time.Time     `json:"expiry_date" gorm:"type:date;index"`
	VerifiedBy          *string        `json:"verified_by" gorm:"type:varchar(36)"`
	VerifiedAt          *time.Time     `json:"verified_at" gorm:"type:timestamptz"`
	Notes               string         `json:"notes" gorm:"type:text"`
	IsActive            bool           `json:"is_active" gorm:"type:boolean;default:true;index"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User     User  `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Verifier *User `json:"verifier,omitempty" gorm:"foreignKey:VerifiedBy;constraint:OnDelete:SET NULL"`
}

// TableName specifies the table name for WorkerSkill
func (WorkerSkill) TableName() string {
	return "worker_skills"
}

// WorkerLocationPreference represents location preferences for workers
type WorkerLocationPreference struct {
	ID                      string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	UserID                  string         `json:"user_id" gorm:"type:varchar(36);not null;index"`
	LocationType            string         `json:"location_type" gorm:"type:varchar(50);not null;check:location_type IN ('suburb', 'postcode', 'region', 'address');index"`
	LocationValue           string         `json:"location_value" gorm:"type:varchar(255);not null"`
	PreferenceLevel         string         `json:"preference_level" gorm:"type:varchar(20);check:preference_level IN ('preferred', 'acceptable', 'avoid');default:'preferred';index"`
	MaxTravelTimeMinutes    int            `json:"max_travel_time_minutes" gorm:"type:integer;default:45"`
	AdditionalTravelRate    float64        `json:"additional_travel_rate" gorm:"type:decimal(8,2);default:0.00"`
	Notes                   string         `json:"notes" gorm:"type:text"`
	IsActive                bool           `json:"is_active" gorm:"type:boolean;default:true;index"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	DeletedAt               gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for WorkerLocationPreference
func (WorkerLocationPreference) TableName() string {
	return "worker_location_preferences"
}

// BeforeCreate sets the CreatedAt and UpdatedAt fields before creating
func (wa *WorkerAvailability) BeforeCreate(tx *gorm.DB) error {
	if wa.ID == "" {
		wa.ID = "avail_" + generateID()
	}
	return nil
}

// BeforeUpdate sets the UpdatedAt field before updating
func (wa *WorkerAvailability) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

// Similar hooks for other models
func (wae *WorkerAvailabilityException) BeforeCreate(tx *gorm.DB) error {
	if wae.ID == "" {
		wae.ID = "avail_ex_" + generateID()
	}
	return nil
}

func (wae *WorkerAvailabilityException) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

func (wp *WorkerPreferences) BeforeCreate(tx *gorm.DB) error {
	if wp.ID == "" {
		wp.ID = "pref_" + generateID()
	}
	return nil
}

func (wp *WorkerPreferences) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

func (ws *WorkerSkill) BeforeCreate(tx *gorm.DB) error {
	if ws.ID == "" {
		ws.ID = "skill_" + generateID()
	}
	return nil
}

func (ws *WorkerSkill) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

func (wlp *WorkerLocationPreference) BeforeCreate(tx *gorm.DB) error {
	if wlp.ID == "" {
		wlp.ID = "loc_" + generateID()
	}
	return nil
}

func (wlp *WorkerLocationPreference) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

// generateID generates a unique ID using UUID
func generateID() string {
	id := uuid.New().String()
	return strings.ReplaceAll(id, "-", "")
}

// Helper methods for WorkerAvailability
func (wa *WorkerAvailability) GetDayName() string {
	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	if wa.DayOfWeek >= 0 && wa.DayOfWeek < len(days) {
		return days[wa.DayOfWeek]
	}
	return "Unknown"
}

// GetFormattedTimeSlot returns a formatted string of the time slot
func (wa *WorkerAvailability) GetFormattedTimeSlot() string {
	return wa.StartTime.Format("15:04") + " - " + wa.EndTime.Format("15:04")
}

// GetTotalAvailableHours calculates total available hours for the time slot
func (wa *WorkerAvailability) GetTotalAvailableHours() float64 {
	duration := wa.EndTime.Sub(wa.StartTime)
	hours := duration.Hours()
	
	// Subtract break time
	breakHours := float64(wa.BreakDurationMinutes) / 60.0
	return hours - breakHours
}

// Helper methods for WorkerPreferences
func (wp *WorkerPreferences) GetWorkCapacityPercentage() float64 {
	if wp.MaxHoursPerWeek == 0 {
		return 0
	}
	return (wp.PreferredHoursPerWeek / wp.MaxHoursPerWeek) * 100
}

// IsSkillExpiringSoon checks if a skill expires within the next 30 days
func (ws *WorkerSkill) IsSkillExpiringSoon() bool {
	if ws.ExpiryDate == nil {
		return false
	}
	
	thirtyDaysFromNow := time.Now().AddDate(0, 0, 30)
	return ws.ExpiryDate.Before(thirtyDaysFromNow)
}

// IsSkillExpired checks if a skill has expired
func (ws *WorkerSkill) IsSkillExpired() bool {
	if ws.ExpiryDate == nil {
		return false
	}
	
	return ws.ExpiryDate.Before(time.Now())
}

// DTO structures for API responses
type WorkerAvailabilityDTO struct {
	ID                   string  `json:"id"`
	UserID               string  `json:"user_id"`
	DayOfWeek            int     `json:"day_of_week"`
	DayName              string  `json:"day_name"`
	StartTime            string  `json:"start_time"`
	EndTime              string  `json:"end_time"`
	IsAvailable          bool    `json:"is_available"`
	MaxHoursPerDay       float64 `json:"max_hours_per_day"`
	PreferredHoursPerDay float64 `json:"preferred_hours_per_day"`
	HourlyRate           *float64 `json:"hourly_rate"`
	BreakDurationMinutes int     `json:"break_duration_minutes"`
	TravelTimeMinutes    int     `json:"travel_time_minutes"`
	Notes                string  `json:"notes"`
	IsActive             bool    `json:"is_active"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}

type WorkerPreferencesDTO struct {
	ID                      string   `json:"id"`
	UserID                  string   `json:"user_id"`
	MaxHoursPerWeek         float64  `json:"max_hours_per_week"`
	PreferredHoursPerWeek   float64  `json:"preferred_hours_per_week"`
	MaxConsecutiveDays      int      `json:"max_consecutive_days"`
	MinHoursBetweenShifts   int      `json:"min_hours_between_shifts"`
	MaxTravelDistanceKm     int      `json:"max_travel_distance_km"`
	PreferredShiftTypes     []string `json:"preferred_shift_types"`
	WillingWeekendWork      bool     `json:"willing_weekend_work"`
	WillingEveningWork      bool     `json:"willing_evening_work"`
	WillingEarlyMorningWork bool     `json:"willing_early_morning_work"`
	RequiresOwnTransport    bool     `json:"requires_own_transport"`
	HasOwnVehicle           bool     `json:"has_own_vehicle"`
	SpecialSkills           []string `json:"special_skills"`
	ClientPreferences       string   `json:"client_preferences"`
	NotificationPreferences JSONB    `json:"notification_preferences"`
	WorkCapacityPercentage  float64  `json:"work_capacity_percentage"`
	IsActive                bool     `json:"is_active"`
	CreatedAt               string   `json:"created_at"`
	UpdatedAt               string   `json:"updated_at"`
}

type WorkerSkillDTO struct {
	ID                  string `json:"id"`
	UserID              string `json:"user_id"`
	SkillCategory       string `json:"skill_category"`
	SkillName           string `json:"skill_name"`
	ProficiencyLevel    string `json:"proficiency_level"`
	CertificationNumber string `json:"certification_number"`
	ExpiryDate          string `json:"expiry_date"`
	IsExpired           bool   `json:"is_expired"`
	IsExpiringSoon      bool   `json:"is_expiring_soon"`
	VerifiedBy          string `json:"verified_by"`
	VerifiedAt          string `json:"verified_at"`
	Notes               string `json:"notes"`
	IsActive            bool   `json:"is_active"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

// Convert model to DTO
func (wa *WorkerAvailability) ToDTO() WorkerAvailabilityDTO {
	return WorkerAvailabilityDTO{
		ID:                   wa.ID,
		UserID:               wa.UserID,
		DayOfWeek:            wa.DayOfWeek,
		DayName:              wa.GetDayName(),
		StartTime:            wa.StartTime.Format("15:04"),
		EndTime:              wa.EndTime.Format("15:04"),
		IsAvailable:          wa.IsAvailable,
		MaxHoursPerDay:       wa.MaxHoursPerDay,
		PreferredHoursPerDay: wa.PreferredHoursPerDay,
		HourlyRate:           wa.HourlyRate,
		BreakDurationMinutes: wa.BreakDurationMinutes,
		TravelTimeMinutes:    wa.TravelTimeMinutes,
		Notes:                wa.Notes,
		IsActive:             wa.IsActive,
		CreatedAt:            wa.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            wa.UpdatedAt.Format(time.RFC3339),
	}
}

func (wp *WorkerPreferences) ToDTO() WorkerPreferencesDTO {
	return WorkerPreferencesDTO{
		ID:                      wp.ID,
		UserID:                  wp.UserID,
		MaxHoursPerWeek:         wp.MaxHoursPerWeek,
		PreferredHoursPerWeek:   wp.PreferredHoursPerWeek,
		MaxConsecutiveDays:      wp.MaxConsecutiveDays,
		MinHoursBetweenShifts:   wp.MinHoursBetweenShifts,
		MaxTravelDistanceKm:     wp.MaxTravelDistanceKm,
		PreferredShiftTypes:     wp.PreferredShiftTypes,
		WillingWeekendWork:      wp.WillingWeekendWork,
		WillingEveningWork:      wp.WillingEveningWork,
		WillingEarlyMorningWork: wp.WillingEarlyMorningWork,
		RequiresOwnTransport:    wp.RequiresOwnTransport,
		HasOwnVehicle:           wp.HasOwnVehicle,
		SpecialSkills:           wp.SpecialSkills,
		ClientPreferences:       wp.ClientPreferences,
		NotificationPreferences: wp.NotificationPreferences,
		WorkCapacityPercentage:  wp.GetWorkCapacityPercentage(),
		IsActive:                wp.IsActive,
		CreatedAt:               wp.CreatedAt.Format(time.RFC3339),
		UpdatedAt:               wp.UpdatedAt.Format(time.RFC3339),
	}
}

func (ws *WorkerSkill) ToDTO() WorkerSkillDTO {
	var expiryDate string
	if ws.ExpiryDate != nil {
		expiryDate = ws.ExpiryDate.Format("2006-01-02")
	}
	
	var verifiedBy string
	if ws.VerifiedBy != nil {
		verifiedBy = *ws.VerifiedBy
	}
	
	var verifiedAt string
	if ws.VerifiedAt != nil {
		verifiedAt = ws.VerifiedAt.Format(time.RFC3339)
	}
	
	return WorkerSkillDTO{
		ID:                  ws.ID,
		UserID:              ws.UserID,
		SkillCategory:       ws.SkillCategory,
		SkillName:           ws.SkillName,
		ProficiencyLevel:    ws.ProficiencyLevel,
		CertificationNumber: ws.CertificationNumber,
		ExpiryDate:          expiryDate,
		IsExpired:           ws.IsSkillExpired(),
		IsExpiringSoon:      ws.IsSkillExpiringSoon(),
		VerifiedBy:          verifiedBy,
		VerifiedAt:          verifiedAt,
		Notes:               ws.Notes,
		IsActive:            ws.IsActive,
		CreatedAt:           ws.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           ws.UpdatedAt.Format(time.RFC3339),
	}
}