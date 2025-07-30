package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Base model with audit fields
type BaseModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uint      `json:"created_by"`
	UpdatedBy uint      `json:"updated_by"`
}

// User represents a system user
type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	
	// Relationships
	Preferences *UserPreferences `gorm:"foreignKey:UserID" json:"preferences,omitempty"`
	CreatedPools []Pool          `gorm:"foreignKey:CreatedBy" json:"-"`
	UpdatedPools []Pool          `gorm:"foreignKey:UpdatedBy" json:"-"`
}

// UserPreferences stores user display preferences
type UserPreferences struct {
	BaseModel
	UserID     uint   `gorm:"not null;uniqueIndex" json:"user_id"`
	UnitSystem string `gorm:"not null;default:'imperial'" json:"unit_system"` // imperial, metric
}

// Pool represents a pool or hot tub
type Pool struct {
	BaseModel
	Name            string  `gorm:"uniqueIndex;not null" json:"name"`
	VolumeGallons   *float64 `json:"volume_gallons,omitempty"`
	Type            string  `json:"type"` // pool, hot_tub
	SystemDescription *string `json:"system_description,omitempty"`
	
	// Relationships
	Samples []Sample `gorm:"foreignKey:PoolID" json:"samples,omitempty"`
}

// Kit represents test kits and equipment
type Kit struct {
	BaseModel
	Name           string     `gorm:"not null" json:"name"`
	Description    *string    `json:"description,omitempty"`
	PurchasedDate  *time.Time `json:"purchased_date,omitempty"`
	ReplenishedDate *time.Time `json:"replenished_date,omitempty"`
	
	// Relationships
	Samples []Sample `gorm:"foreignKey:KitID" json:"samples,omitempty"`
}

// KitJSON is a helper struct for JSON unmarshaling with string dates
type KitJSON struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	Description     *string `json:"description,omitempty"`
	PurchasedDate   *string `json:"purchased_date,omitempty"`
	ReplenishedDate *string `json:"replenished_date,omitempty"`
	CreatedAt       *string `json:"created_at,omitempty"`
	UpdatedAt       *string `json:"updated_at,omitempty"`
	CreatedBy       uint    `json:"created_by"`
	UpdatedBy       uint    `json:"updated_by"`
}

// UnmarshalJSON custom unmarshaler for Kit to handle date-only strings
func (k *Kit) UnmarshalJSON(data []byte) error {
	var kitJSON KitJSON
	if err := json.Unmarshal(data, &kitJSON); err != nil {
		return err
	}
	
	// Set basic fields
	k.ID = kitJSON.ID
	k.Name = kitJSON.Name
	k.Description = kitJSON.Description
	k.CreatedBy = kitJSON.CreatedBy
	k.UpdatedBy = kitJSON.UpdatedBy
	
	// Parse dates - handle both date-only (YYYY-MM-DD) and full datetime formats
	if kitJSON.PurchasedDate != nil && *kitJSON.PurchasedDate != "" {
		dateStr := strings.TrimSpace(*kitJSON.PurchasedDate)
		if dateStr != "" {
			var parsedDate time.Time
			var err error
			
			// Try date-only format first
			if len(dateStr) == 10 {
				parsedDate, err = time.Parse("2006-01-02", dateStr)
			} else {
				// Try full datetime format
				parsedDate, err = time.Parse(time.RFC3339, dateStr)
			}
			
			if err == nil {
				k.PurchasedDate = &parsedDate
			}
		}
	}
	
	if kitJSON.ReplenishedDate != nil && *kitJSON.ReplenishedDate != "" {
		dateStr := strings.TrimSpace(*kitJSON.ReplenishedDate)
		if dateStr != "" {
			var parsedDate time.Time
			var err error
			
			// Try date-only format first
			if len(dateStr) == 10 {
				parsedDate, err = time.Parse("2006-01-02", dateStr)
			} else {
				// Try full datetime format
				parsedDate, err = time.Parse(time.RFC3339, dateStr)
			}
			
			if err == nil {
				k.ReplenishedDate = &parsedDate
			}
		}
	}
	
	return nil
}

// Sample represents a water testing session
type Sample struct {
	BaseModel
	PoolID         uint      `gorm:"not null" json:"pool_id"`
	SampleDateTime time.Time `gorm:"column:sample_date_time;not null" json:"sample_datetime"`
	UserID         uint      `gorm:"not null" json:"user_id"`
	KitID          uint      `gorm:"not null" json:"kit_id"`
	Notes          string    `gorm:"type:text" json:"notes"`
	
	// Relationships
	Pool         *Pool         `gorm:"foreignKey:PoolID" json:"pool,omitempty"`
	User         *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Kit          *Kit          `gorm:"foreignKey:KitID" json:"kit,omitempty"`
	Measurements *Measurements `gorm:"foreignKey:SampleID" json:"measurements,omitempty"`
	Indices      *Indices      `gorm:"foreignKey:SampleID" json:"indices,omitempty"`
}

// SampleJSON is a helper struct for JSON unmarshaling with string datetime
type SampleJSON struct {
	ID             uint                   `json:"id"`
	PoolID         uint                   `json:"pool_id"`
	SampleDateTime string                 `json:"sample_datetime"`
	UserID         uint                   `json:"user_id"`
	KitID          uint                   `json:"kit_id"`
	Notes          string                 `json:"notes"`
	CreatedAt      *string                `json:"created_at,omitempty"`
	UpdatedAt      *string                `json:"updated_at,omitempty"`
	CreatedBy      uint                   `json:"created_by"`
	UpdatedBy      uint                   `json:"updated_by"`
	Measurements   map[string]interface{} `json:"measurements,omitempty"`
}

// UnmarshalJSON custom unmarshaler for Sample to handle datetime-local format
func (s *Sample) UnmarshalJSON(data []byte) error {
	var sampleJSON SampleJSON
	if err := json.Unmarshal(data, &sampleJSON); err != nil {
		return err
	}
	
	// Set basic fields
	s.ID = sampleJSON.ID
	s.PoolID = sampleJSON.PoolID
	s.UserID = sampleJSON.UserID
	s.KitID = sampleJSON.KitID
	s.Notes = sampleJSON.Notes
	s.CreatedBy = sampleJSON.CreatedBy
	s.UpdatedBy = sampleJSON.UpdatedBy
	
	// Parse the datetime - handle multiple formats
	if sampleJSON.SampleDateTime != "" {
		dateStr := strings.TrimSpace(sampleJSON.SampleDateTime)
		var parsedTime time.Time
		var err error
		
		// Try datetime-local format first (YYYY-MM-DDTHH:MM)
		if len(dateStr) == 16 && strings.Count(dateStr, "T") == 1 {
			parsedTime, err = time.Parse("2006-01-02T15:04", dateStr)
		} else if len(dateStr) == 19 && strings.Count(dateStr, "T") == 1 {
			// Try datetime-local with seconds (YYYY-MM-DDTHH:MM:SS)
			parsedTime, err = time.Parse("2006-01-02T15:04:05", dateStr)
		} else {
			// Try full RFC3339 format
			parsedTime, err = time.Parse(time.RFC3339, dateStr)
		}
		
		if err != nil {
			return fmt.Errorf("failed to parse sample_datetime '%s': %v", dateStr, err)
		}
		
		s.SampleDateTime = parsedTime
	}
	
	// Parse measurements if provided
	if sampleJSON.Measurements != nil {
		measurements := &Measurements{}
		// Initialize BaseModel fields to ensure they're not set from JSON
		measurements.BaseModel = BaseModel{
			ID:        0, // Let database generate
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			CreatedBy: 0,
			UpdatedBy: 0,
		}
		
		// Helper function to convert interface{} to float64
		getFloat := func(key string) float64 {
			if val, ok := sampleJSON.Measurements[key]; ok {
				if fval, ok := val.(float64); ok {
					return fval
				}
			}
			return 0
		}
		
		// Helper function to convert interface{} to *float64
		getFloatPtr := func(key string) *float64 {
			if val, ok := sampleJSON.Measurements[key]; ok {
				if fval, ok := val.(float64); ok {
					return &fval
				}
			}
			return nil
		}
		
		// Helper function to convert interface{} to *string
		getStringPtr := func(key string) *string {
			if val, ok := sampleJSON.Measurements[key]; ok {
				if sval, ok := val.(string); ok && sval != "" {
					return &sval
				}
			}
			return nil
		}
		
		// Parse all measurement fields
		measurements.FC = getFloat("fc")
		measurements.TC = getFloat("tc")
		measurements.PH = getFloat("ph")
		measurements.TA = getFloat("ta")
		measurements.CH = getFloat("ch")
		measurements.Temperature = getFloat("temperature")
		measurements.CYA = getFloatPtr("cya")
		measurements.Salinity = getFloatPtr("salinity")
		measurements.TDS = getFloatPtr("tds")
		measurements.Appearance = getStringPtr("appearance")
		measurements.Maintenance = getStringPtr("maintenance")
		
		s.Measurements = measurements
	}
	
	return nil
}

// Measurements stores water chemistry measurements
type Measurements struct {
	BaseModel
	SampleID     uint     `gorm:"not null;uniqueIndex" json:"sample_id"`
	FC           float64  `gorm:"not null" json:"fc"`           // Free Chlorine (ppm)
	TC           float64  `gorm:"not null" json:"tc"`           // Total Chlorine (ppm)
	PH           float64  `gorm:"not null" json:"ph"`           // pH (0-14 scale)
	TA           float64  `gorm:"not null" json:"ta"`           // Total Alkalinity (ppm)
	CH           float64  `gorm:"not null" json:"ch"`           // Calcium Hardness (ppm)
	CYA          *float64 `json:"cya,omitempty"`                // Cyanuric Acid (ppm)
	Temperature  float64  `gorm:"not null" json:"temperature"`  // Temperature (°F)
	Salinity     *float64 `json:"salinity,omitempty"`           // Salinity (ppm)
	TDS          *float64 `json:"tds,omitempty"`                // Total Dissolved Solids (mg/l)
	Appearance   *string  `json:"appearance,omitempty"`         // Water appearance notes
	Maintenance  *string  `json:"maintenance,omitempty"`        // Maintenance notes
}

// Indices stores calculated water balance indices
type Indices struct {
	BaseModel
	SampleID uint     `gorm:"not null;uniqueIndex" json:"sample_id"`
	LSI      *float64 `json:"lsi,omitempty"` // Langelier Saturation Index
	RSI      *float64 `json:"rsi,omitempty"` // Ryznar Stability Index
	Comment  *string  `json:"comment,omitempty"` // Notes about estimation/missing parameters
}

// BeforeCreate hook to set audit fields
func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if userID, ok := tx.Statement.Context.Value("user_id").(uint); ok {
		m.CreatedBy = userID
		m.UpdatedBy = userID
	}
	return nil
}

// BeforeUpdate hook to set audit fields
func (m *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	if userID, ok := tx.Statement.Context.Value("user_id").(uint); ok {
		m.UpdatedBy = userID
	}
	return nil
}