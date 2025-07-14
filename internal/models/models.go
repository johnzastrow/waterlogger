package models

import (
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

// Sample represents a water testing session
type Sample struct {
	BaseModel
	PoolID         uint      `gorm:"not null" json:"pool_id"`
	SampleDateTime time.Time `gorm:"not null" json:"sample_datetime"`
	UserID         uint      `gorm:"not null" json:"user_id"`
	KitID          uint      `gorm:"not null" json:"kit_id"`
	
	// Relationships
	Pool         *Pool         `gorm:"foreignKey:PoolID" json:"pool,omitempty"`
	User         *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Kit          *Kit          `gorm:"foreignKey:KitID" json:"kit,omitempty"`
	Measurements *Measurements `gorm:"foreignKey:SampleID" json:"measurements,omitempty"`
	Indices      *Indices      `gorm:"foreignKey:SampleID" json:"indices,omitempty"`
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