package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Beat struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"not null;index"`
	Title       string    `json:"title" gorm:"not null" validate:"required,max=100"`
	Description string    `json:"description" gorm:"type:text"`
	Genre       string    `json:"genre" gorm:"not null" validate:"required"`
	BPM         int       `json:"bpm" validate:"min=60,max=200"`
	Key         string    `json:"key"`
	Duration    int       `json:"duration"` // in seconds

	// File storage
	AudioURL    string `json:"audio_url" gorm:"not null"`
	AudioS3Key  string `json:"audio_s3_key" gorm:"not null"`
	FileSize    int64  `json:"file_size"`
	AudioFormat string `json:"audio_format"` // mp3, wav, etc.

	// Metadata
	Tags       string `json:"tags"` // comma-separated
	IsPublic   bool   `json:"is_public" gorm:"default:true"`
	IsApproved bool   `json:"is_approved" gorm:"default:false"`

	// Engagement metrics
	PlayCount     int `json:"play_count" gorm:"default:0"`
	LikeCount     int `json:"like_count" gorm:"default:0"`
	DownloadCount int `json:"download_count" gorm:"default:0"`
	ShareCount    int `json:"share_count" gorm:"default:0"`

	// Gamification
	Points int `json:"points" gorm:"default:0"`

	// Relationships
	User                 User                  `json:"user" gorm:"foreignKey:UserID"`
	Collaborations       []Collaboration       `json:"collaborations,omitempty" gorm:"foreignKey:BeatID"`
	Comments             []Comment             `json:"comments,omitempty" gorm:"foreignKey:BeatID"`
	Likes                []Like                `json:"likes,omitempty" gorm:"foreignKey:BeatID"`
	ChallengeSubmissions []ChallengeSubmission `json:"challenge_submissions,omitempty" gorm:"foreignKey:BeatID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Comment struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null;index"`
	BeatID    uuid.UUID `json:"beat_id" gorm:"not null;index"`
	Content   string    `json:"content" gorm:"type:text;not null" validate:"required,max=500"`
	Timestamp int       `json:"timestamp"` // timestamp in audio (seconds) for timestamped feedback

	User User `json:"user" gorm:"foreignKey:UserID"`
	Beat Beat `json:"beat" gorm:"foreignKey:BeatID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Like struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID uuid.UUID `json:"user_id" gorm:"not null"`
	BeatID uuid.UUID `json:"beat_id" gorm:"not null"`

	User User `json:"user" gorm:"foreignKey:UserID"`
	Beat Beat `json:"beat" gorm:"foreignKey:BeatID"`

	CreatedAt time.Time `json:"created_at"`
}

type Collaboration struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BeatID       uuid.UUID `json:"beat_id" gorm:"not null;index"`
	UserID       uuid.UUID `json:"user_id" gorm:"not null;index"` // collaborator
	Role         string    `json:"role" gorm:"not null"`          // producer, mixer, vocalist, etc.
	Status       string    `json:"status" gorm:"default:pending"` // pending, accepted, rejected
	Contribution string    `json:"contribution" gorm:"type:text"`

	Beat Beat `json:"beat" gorm:"foreignKey:BeatID"`
	User User `json:"user" gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
