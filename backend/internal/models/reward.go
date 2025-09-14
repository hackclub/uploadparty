package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RewardType string

const (
	RewardTypeNILicense RewardType = "ni_license"
	RewardTypeBadge     RewardType = "badge"
	RewardTypePoints    RewardType = "points"
	RewardTypeCustom    RewardType = "custom"
)

type RewardStatus string

const (
	RewardStatusPending   RewardStatus = "pending"
	RewardStatusApproved  RewardStatus = "approved"
	RewardStatusDelivered RewardStatus = "delivered"
	RewardStatusRejected  RewardStatus = "rejected"
	RewardStatusExpired   RewardStatus = "expired"
)

type Reward struct {
	ID     uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID uuid.UUID    `json:"user_id" gorm:"not null;index"`
	Type   RewardType   `json:"type" gorm:"not null"`
	Status RewardStatus `json:"status" gorm:"default:pending"`

	// Reward Details
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description" gorm:"type:text"`
	Value       string `json:"value"` // License key, badge ID, point amount, etc.

	// Source Information
	SourceType string     `json:"source_type"` // challenge, achievement, manual, etc.
	SourceID   *uuid.UUID `json:"source_id,omitempty"`

	// Admin Management
	ApprovedBy  *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	RejectedAt  *time.Time `json:"rejected_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`

	// Notes and Metadata
	AdminNotes string `json:"admin_notes,omitempty" gorm:"type:text"`
	UserNotes  string `json:"user_notes,omitempty" gorm:"type:text"`
	Metadata   string `json:"metadata,omitempty" gorm:"type:json"`

	// Relationships
	User     User  `json:"user" gorm:"foreignKey:UserID"`
	Approver *User `json:"approver,omitempty" gorm:"foreignKey:ApprovedBy"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type RewardClaim struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RewardID  uuid.UUID `json:"reward_id" gorm:"not null;index"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null;index"`
	ClaimData string    `json:"claim_data" gorm:"type:json"`     // User-provided claim information
	Status    string    `json:"status" gorm:"default:submitted"` // submitted, processing, completed, failed

	// Relationships
	Reward Reward `json:"reward" gorm:"foreignKey:RewardID"`
	User   User   `json:"user" gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NILicense represents a Native Instruments license
type NILicense struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LicenseKey  string     `json:"license_key" gorm:"uniqueIndex;not null"`
	ProductName string     `json:"product_name" gorm:"not null"` // e.g., "Komplete 15"
	ProductCode string     `json:"product_code"`
	IsAssigned  bool       `json:"is_assigned" gorm:"default:false"`
	AssignedTo  *uuid.UUID `json:"assigned_to,omitempty"`
	AssignedAt  *time.Time `json:"assigned_at,omitempty"`
	ClaimedAt   *time.Time `json:"claimed_at,omitempty"`

	// Batch information for tracking
	BatchID    string    `json:"batch_id"`
	ImportedBy uuid.UUID `json:"imported_by"` // Admin who imported this license

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:AssignedTo"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// RewardTemplate for creating standardized rewards
type RewardTemplate struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name             string     `json:"name" gorm:"not null"`
	Type             RewardType `json:"type" gorm:"not null"`
	Title            string     `json:"title" gorm:"not null"`
	Description      string     `json:"description" gorm:"type:text"`
	DefaultValue     string     `json:"default_value"`
	IsActive         bool       `json:"is_active" gorm:"default:true"`
	RequiresApproval bool       `json:"requires_approval" gorm:"default:true"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AdminAction for audit trail
type AdminAction struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AdminID    uuid.UUID  `json:"admin_id" gorm:"not null;index"`
	Action     string     `json:"action" gorm:"not null"` // approve_reward, reject_reward, assign_license, etc.
	TargetID   *uuid.UUID `json:"target_id,omitempty"`
	TargetType string     `json:"target_type"` // reward, license, user, etc.
	Details    string     `json:"details" gorm:"type:text"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`

	// Relationships
	Admin User `json:"admin" gorm:"foreignKey:AdminID"`

	CreatedAt time.Time `json:"created_at"`
}
