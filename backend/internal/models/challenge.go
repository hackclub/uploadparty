package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Challenge struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"not null" validate:"required,max=100"`
	Description string    `json:"description" gorm:"type:text;not null"`
	Rules       string    `json:"rules" gorm:"type:text"`

	// Challenge parameters
	Genre               string `json:"genre"`
	BPMRange            string `json:"bpm_range"` // e.g., "120-140"
	KeyRequirement      string `json:"key_requirement"`
	MaxDuration         int    `json:"max_duration"`         // in seconds
	RequiredInstruments string `json:"required_instruments"` // comma-separated

	// Timing
	StartDate time.Time `json:"start_date" gorm:"not null"`
	EndDate   time.Time `json:"end_date" gorm:"not null"`

	// Rewards
	PointsReward    int        `json:"points_reward" gorm:"default:0"`
	BadgeReward     *uuid.UUID `json:"badge_reward,omitempty"`
	NILicenseReward bool       `json:"ni_license_reward" gorm:"default:false"`
	CustomReward    string     `json:"custom_reward"`

	// Status
	IsActive bool   `json:"is_active" gorm:"default:true"`
	Status   string `json:"status" gorm:"default:upcoming"` // upcoming, active, judging, completed

	// Engagement
	ParticipantCount int `json:"participant_count" gorm:"default:0"`
	SubmissionCount  int `json:"submission_count" gorm:"default:0"`

	// Judging
	JudgingCriteria string `json:"judging_criteria" gorm:"type:text"`
	MaxWinners      int    `json:"max_winners" gorm:"default:1"`

	// Relationships
	Submissions []ChallengeSubmission `json:"submissions,omitempty" gorm:"foreignKey:ChallengeID"`
	Winners     []ChallengeWinner     `json:"winners,omitempty" gorm:"foreignKey:ChallengeID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type ChallengeSubmission struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ChallengeID uuid.UUID `json:"challenge_id" gorm:"not null;index"`
	UserID      uuid.UUID `json:"user_id" gorm:"not null;index"`
	BeatID      uuid.UUID `json:"beat_id" gorm:"not null;index"`

	// Submission details
	SubmissionNotes string `json:"submission_notes" gorm:"type:text"`
	IsQualified     bool   `json:"is_qualified" gorm:"default:true"`

	// Scoring (for judging)
	CreativityScore  float64 `json:"creativity_score" gorm:"default:0"`
	TechnicalScore   float64 `json:"technical_score" gorm:"default:0"`
	OriginalityScore float64 `json:"originality_score" gorm:"default:0"`
	TotalScore       float64 `json:"total_score" gorm:"default:0"`

	// Relationships
	Challenge Challenge `json:"challenge" gorm:"foreignKey:ChallengeID"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Beat      Beat      `json:"beat" gorm:"foreignKey:BeatID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChallengeWinner struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ChallengeID  uuid.UUID `json:"challenge_id" gorm:"not null"`
	UserID       uuid.UUID `json:"user_id" gorm:"not null"`
	SubmissionID uuid.UUID `json:"submission_id" gorm:"not null"`
	Position     int       `json:"position" gorm:"not null"` // 1st, 2nd, 3rd place

	// Prize details
	PointsAwarded    int        `json:"points_awarded" gorm:"default:0"`
	BadgeAwarded     *uuid.UUID `json:"badge_awarded,omitempty"`
	NILicenseAwarded bool       `json:"ni_license_awarded" gorm:"default:false"`
	CustomPrize      string     `json:"custom_prize"`

	// Relationships
	Challenge  Challenge           `json:"challenge" gorm:"foreignKey:ChallengeID"`
	User       User                `json:"user" gorm:"foreignKey:UserID"`
	Submission ChallengeSubmission `json:"submission" gorm:"foreignKey:SubmissionID"`

	CreatedAt time.Time `json:"created_at"`
}

type Score struct {
	ID       uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID   uuid.UUID  `json:"user_id" gorm:"not null;index"`
	Points   int        `json:"points" gorm:"not null"`
	Source   string     `json:"source" gorm:"not null"` // beat_upload, challenge_win, collaboration, etc.
	SourceID *uuid.UUID `json:"source_id,omitempty"`    // ID of the beat, challenge, etc.
	Reason   string     `json:"reason"`

	User User `json:"user" gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"`
}
