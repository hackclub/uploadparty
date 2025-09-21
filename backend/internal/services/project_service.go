package services

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/uploadparty/app/backend/internal/models"
)

type ProjectService struct{ DB *gorm.DB }

func NewProjectService(db *gorm.DB) *ProjectService { return &ProjectService{DB: db} }

type UpsertProjectInput struct {
	Title           string          `json:"title"`
	DAW             string          `json:"daw"`
	PluginVersion   string          `json:"pluginVersion"`
	DurationSeconds int             `json:"durationSeconds"`
	Metadata        json.RawMessage `json:"metadata"`
	Public          *bool           `json:"public"`
}

func (s *ProjectService) UpsertByTitle(userID uint, in UpsertProjectInput) (*models.Project, error) {
	if in.Title == "" {
		return nil, errors.New("title required")
	}
	var p models.Project
	err := s.DB.Where("user_id = ? AND title = ?", userID, in.Title).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		p = models.Project{UserID: userID, Title: in.Title}
	}
	p.DAW = in.DAW
	p.PluginVersion = in.PluginVersion
	p.DurationSeconds = in.DurationSeconds
	if in.Metadata != nil {
		p.Metadata = in.Metadata
	}
	if in.Public != nil {
		p.Public = *in.Public
	}
	if p.ID == 0 {
		if err := s.DB.Create(&p).Error; err != nil {
			return nil, err
		}
	} else {
		if err := s.DB.Save(&p).Error; err != nil {
			return nil, err
		}
	}
	return &p, nil
}

func (s *ProjectService) MarkComplete(userID, projectID uint) (*models.Project, error) {
	var p models.Project
	if err := s.DB.Where("user_id = ? AND id = ?", userID, projectID).First(&p).Error; err != nil {
		return nil, err
	}
	now := time.Now()
	p.Status = models.StatusComplete
	p.CompletedAt = &now
	if err := s.DB.Save(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *ProjectService) ListPublicByUser(userID uint) ([]models.Project, error) {
	var ps []models.Project
	if err := s.DB.Where("user_id = ? AND public = ?", userID, true).Preload("Plugins").Order("created_at desc").Find(&ps).Error; err != nil {
		return nil, err
	}
	return ps, nil
}

func (s *ProjectService) ListByUser(userID uint) ([]models.Project, error) {
	var ps []models.Project
	if err := s.DB.Where("user_id = ?", userID).Preload("Plugins").Order("created_at desc").Find(&ps).Error; err != nil {
		return nil, err
	}
	return ps, nil
}
