package services

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"

	"github.com/uploadparty/app/backend/internal/models"
)

type PluginService struct{ DB *gorm.DB }

func NewPluginService(db *gorm.DB) *PluginService { return &PluginService{DB: db} }

type UpsertPluginInput struct {
	Name     string          `json:"name"`
	Vendor   string          `json:"vendor"`
	Version  string          `json:"version"`
	Format   string          `json:"format"`
	Metadata json.RawMessage `json:"metadata"`
}

// ensureProjectOwned checks that the project belongs to the user and returns it.
func (s *PluginService) ensureProjectOwned(userID, projectID uint) (*models.Project, error) {
	var p models.Project
	if err := s.DB.Where("user_id = ? AND id = ?", userID, projectID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// UpsertByName creates or updates a plugin for a project identified by name (unique per project).
func (s *PluginService) UpsertByName(userID, projectID uint, in UpsertPluginInput) (*models.Plugin, error) {
	if in.Name == "" {
		return nil, errors.New("name required")
	}
	if _, err := s.ensureProjectOwned(userID, projectID); err != nil {
		return nil, err
	}
	var pl models.Plugin
	err := s.DB.Where("project_id = ? AND name = ?", projectID, in.Name).First(&pl).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		pl = models.Plugin{ProjectID: projectID, Name: in.Name}
	}
	pl.Vendor = in.Vendor
	pl.Version = in.Version
	pl.Format = in.Format
	if in.Metadata != nil {
		pl.Metadata = in.Metadata
	}
	if pl.ID == 0 {
		if err := s.DB.Create(&pl).Error; err != nil {
			return nil, err
		}
	} else {
		if err := s.DB.Save(&pl).Error; err != nil {
			return nil, err
		}
	}
	return &pl, nil
}

func (s *PluginService) ListByProject(userID, projectID uint) ([]models.Plugin, error) {
	if _, err := s.ensureProjectOwned(userID, projectID); err != nil {
		return nil, err
	}
	var items []models.Plugin
	if err := s.DB.Where("project_id = ?", projectID).Order("name asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
