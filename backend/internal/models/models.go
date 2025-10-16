package models

import (
	"time"

	"gorm.io/datatypes"
)

type RSVP struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Email     string `gorm:"uniqueIndex;size:255;not null" json:"email"`
	EmailSent bool   `gorm:"default:false" json:"emailSent"`
	UserID    *uint  `gorm:"index" json:"userId,omitempty"`
	User      *User  `gorm:"constraint:OnDelete:SET NULL" json:"user,omitempty"`
}

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Auth0 fields for social login
	Auth0ID string `gorm:"uniqueIndex;size:255" json:"auth0Id,omitempty"`
	Email   string `gorm:"uniqueIndex;size:255" json:"email"`

	// Username and password for legacy auth (optional with Auth0)
	Username     string `gorm:"uniqueIndex;size:50" json:"username"`
	PasswordHash string `json:"-"`

	// Profile fields
	DisplayName string `gorm:"size:100" json:"displayName"`
	Picture     string `gorm:"size:500" json:"picture,omitempty"`
	Bio         string `gorm:"size:280" json:"bio"`
	Public      bool   `json:"public"`
}

type ProjectStatus string

const (
	StatusInProgress ProjectStatus = "in_progress"
	StatusComplete   ProjectStatus = "complete"
)

type Project struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	UserID uint `json:"userId"`
	User   User `gorm:"constraint:OnDelete:CASCADE" json:"-"`

	Title           string         `gorm:"size:200" json:"title"`
	DAW             string         `gorm:"size:100" json:"daw"`
	PluginVersion   string         `gorm:"size:50" json:"pluginVersion"`
	DurationSeconds int            `json:"durationSeconds"`
	Metadata        datatypes.JSON `json:"metadata"`
	Status          ProjectStatus  `gorm:"size:20;default:in_progress" json:"status"`
	CompletedAt     *time.Time     `json:"completedAt"`
	Public          bool           `json:"public"`

	Plugins []Plugin `json:"plugins,omitempty"`
}

// Plugin represents a plugin instance attached to a project.
// We keep simple identifying fields and arbitrary metadata from the VST.
// A project can have many plugins.
type Plugin struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	ProjectID uint    `gorm:"index:idx_project_name,unique" json:"projectId"`
	Project   Project `gorm:"constraint:OnDelete:CASCADE" json:"-"`

	Name     string         `gorm:"size:120;index:idx_project_name,unique" json:"name"`
	Vendor   string         `gorm:"size:120" json:"vendor"`
	Version  string         `gorm:"size:50" json:"version"`
	Format   string         `gorm:"size:20" json:"format"` // e.g., VST3, AU, AAX
	Metadata datatypes.JSON `json:"metadata"`
}
