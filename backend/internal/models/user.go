package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username     string    `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=30"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	PasswordHash string    `json:"-" gorm:"not null"`
	FirstName    string    `json:"first_name" gorm:"not null" validate:"required,max=50"`
	LastName     string    `json:"last_name" gorm:"not null" validate:"required,max=50"`
	Age          int       `json:"age" gorm:"not null" validate:"required,min=13,max=19"`
	School       string    `json:"school" gorm:"not null" validate:"required,max=100"`
	Bio          string    `json:"bio" gorm:"type:text"`
	AvatarURL    string    `json:"avatar_url"`
	IsVerified   bool      `json:"is_verified" gorm:"default:false"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`

	// Gamification
	TotalPoints  int `json:"total_points" gorm:"default:0"`
	Level        int `json:"level" gorm:"default:1"`
	RankPosition int `json:"rank_position" gorm:"default:0"`

	// Native Instruments License
	NILicenseID     *string `json:"ni_license_id,omitempty"`
	NILicenseStatus string  `json:"ni_license_status" gorm:"default:pending"`

	// Relationships
	Beats          []Beat          `json:"beats,omitempty" gorm:"foreignKey:UserID"`
	Collaborations []Collaboration `json:"collaborations,omitempty" gorm:"foreignKey:UserID"`
	Scores         []Score         `json:"scores,omitempty" gorm:"foreignKey:UserID"`
	Badges         []UserBadge     `json:"badges,omitempty" gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserBadge struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID   uuid.UUID `json:"user_id" gorm:"not null"`
	BadgeID  uuid.UUID `json:"badge_id" gorm:"not null"`
	EarnedAt time.Time `json:"earned_at" gorm:"default:CURRENT_TIMESTAMP"`

	User  User  `json:"user" gorm:"foreignKey:UserID"`
	Badge Badge `json:"badge" gorm:"foreignKey:BadgeID"`
}

type Badge struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	IconURL     string    `json:"icon_url"`
	Points      int       `json:"points" gorm:"default:0"`
	Rarity      string    `json:"rarity" gorm:"default:common"` // common, rare, epic, legendary

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
