package services

import (
	"errors"
	"fmt"
	"time"
	"uploadparty/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RewardService struct {
	db *gorm.DB
}

type CreateRewardRequest struct {
	UserID      string            `json:"user_id" validate:"required"`
	Type        models.RewardType `json:"type" validate:"required"`
	Title       string            `json:"title" validate:"required"`
	Description string            `json:"description"`
	Value       string            `json:"value"`
	SourceType  string            `json:"source_type"`
	SourceID    *string           `json:"source_id,omitempty"`
	AdminNotes  string            `json:"admin_notes,omitempty"`
	ExpiresAt   *time.Time        `json:"expires_at,omitempty"`
}

type ApproveRewardRequest struct {
	AdminNotes string `json:"admin_notes,omitempty"`
}

type RejectRewardRequest struct {
	AdminNotes string `json:"admin_notes" validate:"required"`
}

type ImportLicensesRequest struct {
	Licenses    []NILicenseImport `json:"licenses" validate:"required,min=1"`
	ProductName string            `json:"product_name" validate:"required"`
	ProductCode string            `json:"product_code"`
	BatchID     string            `json:"batch_id" validate:"required"`
}

type NILicenseImport struct {
	LicenseKey string `json:"license_key" validate:"required"`
}

func NewRewardService(db *gorm.DB) *RewardService {
	return &RewardService{db: db}
}

// Admin Functions

func (s *RewardService) GetPendingRewards(page, limit int) ([]models.Reward, int64, error) {
	var rewards []models.Reward
	var total int64

	// Count total pending rewards
	if err := s.db.Model(&models.Reward{}).
		Where("status = ?", models.RewardStatusPending).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated rewards with user info
	offset := (page - 1) * limit
	if err := s.db.Preload("User").
		Where("status = ?", models.RewardStatusPending).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&rewards).Error; err != nil {
		return nil, 0, err
	}

	return rewards, total, nil
}

func (s *RewardService) GetAllRewards(status string, page, limit int) ([]models.Reward, int64, error) {
	var rewards []models.Reward
	var total int64

	query := s.db.Model(&models.Reward{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated rewards
	offset := (page - 1) * limit
	if err := query.Preload("User").Preload("Approver").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&rewards).Error; err != nil {
		return nil, 0, err
	}

	return rewards, total, nil
}

func (s *RewardService) ApproveReward(rewardID, adminID string, req *ApproveRewardRequest) error {
	adminUUID, err := uuid.Parse(adminID)
	if err != nil {
		return errors.New("invalid admin ID")
	}

	rewardUUID, err := uuid.Parse(rewardID)
	if err != nil {
		return errors.New("invalid reward ID")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		var reward models.Reward
		if err := tx.Where("id = ? AND status = ?", rewardUUID, models.RewardStatusPending).
			First(&reward).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("reward not found or already processed")
			}
			return err
		}

		now := time.Now()
		updates := map[string]interface{}{
			"status":      models.RewardStatusApproved,
			"approved_by": adminUUID,
			"approved_at": &now,
		}

		if req.AdminNotes != "" {
			updates["admin_notes"] = req.AdminNotes
		}

		if err := tx.Model(&reward).Updates(updates).Error; err != nil {
			return err
		}

		// Handle specific reward types
		if err := s.handleRewardApproval(tx, &reward); err != nil {
			return err
		}

		// Log admin action
		return s.logAdminAction(tx, adminUUID, "approve_reward", &rewardUUID, "reward",
			fmt.Sprintf("Approved reward %s for user %s", reward.Title, reward.UserID))
	})
}

func (s *RewardService) RejectReward(rewardID, adminID string, req *RejectRewardRequest) error {
	adminUUID, err := uuid.Parse(adminID)
	if err != nil {
		return errors.New("invalid admin ID")
	}

	rewardUUID, err := uuid.Parse(rewardID)
	if err != nil {
		return errors.New("invalid reward ID")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		var reward models.Reward
		if err := tx.Where("id = ? AND status = ?", rewardUUID, models.RewardStatusPending).
			First(&reward).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("reward not found or already processed")
			}
			return err
		}

		now := time.Now()
		updates := map[string]interface{}{
			"status":      models.RewardStatusRejected,
			"approved_by": adminUUID,
			"rejected_at": &now,
			"admin_notes": req.AdminNotes,
		}

		if err := tx.Model(&reward).Updates(updates).Error; err != nil {
			return err
		}

		// Log admin action
		return s.logAdminAction(tx, adminUUID, "reject_reward", &rewardUUID, "reward",
			fmt.Sprintf("Rejected reward %s for user %s: %s", reward.Title, reward.UserID, req.AdminNotes))
	})
}

func (s *RewardService) CreateReward(adminID string, req *CreateRewardRequest) (*models.Reward, error) {
	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	adminUUID, err := uuid.Parse(adminID)
	if err != nil {
		return nil, errors.New("invalid admin ID")
	}

	reward := models.Reward{
		UserID:      userUUID,
		Type:        req.Type,
		Status:      models.RewardStatusPending,
		Title:       req.Title,
		Description: req.Description,
		Value:       req.Value,
		SourceType:  req.SourceType,
		AdminNotes:  req.AdminNotes,
		ExpiresAt:   req.ExpiresAt,
	}

	if req.SourceID != nil {
		sourceUUID, err := uuid.Parse(*req.SourceID)
		if err != nil {
			return nil, errors.New("invalid source ID")
		}
		reward.SourceID = &sourceUUID
	}

	if err := s.db.Create(&reward).Error; err != nil {
		return nil, err
	}

	// Log admin action
	s.logAdminAction(s.db, adminUUID, "create_reward", &reward.ID, "reward",
		fmt.Sprintf("Created reward %s for user %s", reward.Title, reward.UserID))

	return &reward, nil
}

// NI License Management

func (s *RewardService) ImportLicenses(adminID string, req *ImportLicensesRequest) error {
	adminUUID, err := uuid.Parse(adminID)
	if err != nil {
		return errors.New("invalid admin ID")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, licenseImport := range req.Licenses {
			license := models.NILicense{
				LicenseKey:  licenseImport.LicenseKey,
				ProductName: req.ProductName,
				ProductCode: req.ProductCode,
				BatchID:     req.BatchID,
				ImportedBy:  adminUUID,
				IsAssigned:  false,
			}

			if err := tx.Create(&license).Error; err != nil {
				return fmt.Errorf("failed to import license %s: %w", license.LicenseKey, err)
			}
		}

		// Log admin action
		return s.logAdminAction(tx, adminUUID, "import_licenses", nil, "ni_license",
			fmt.Sprintf("Imported %d licenses for %s (batch: %s)", len(req.Licenses), req.ProductName, req.BatchID))
	})
}

func (s *RewardService) GetAvailableLicenses(productName string, page, limit int) ([]models.NILicense, int64, error) {
	var licenses []models.NILicense
	var total int64

	query := s.db.Model(&models.NILicense{}).Where("is_assigned = ?", false)
	if productName != "" {
		query = query.Where("product_name = ?", productName)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&licenses).Error; err != nil {
		return nil, 0, err
	}

	return licenses, total, nil
}

func (s *RewardService) GetAssignedLicenses(page, limit int) ([]models.NILicense, int64, error) {
	var licenses []models.NILicense
	var total int64

	query := s.db.Model(&models.NILicense{}).Where("is_assigned = ?", true)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("User").
		Order("assigned_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&licenses).Error; err != nil {
		return nil, 0, err
	}

	return licenses, total, nil
}

// User Functions

func (s *RewardService) GetUserRewards(userID string, status string, page, limit int) ([]models.Reward, int64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, errors.New("invalid user ID")
	}

	var rewards []models.Reward
	var total int64

	query := s.db.Model(&models.Reward{}).Where("user_id = ?", userUUID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&rewards).Error; err != nil {
		return nil, 0, err
	}

	return rewards, total, nil
}

// Helper Functions

func (s *RewardService) handleRewardApproval(tx *gorm.DB, reward *models.Reward) error {
	switch reward.Type {
	case models.RewardTypeNILicense:
		return s.assignNILicense(tx, reward)
	case models.RewardTypePoints:
		return s.addPointsToUser(tx, reward)
	case models.RewardTypeBadge:
		return s.assignBadgeToUser(tx, reward)
	}
	return nil
}

func (s *RewardService) assignNILicense(tx *gorm.DB, reward *models.Reward) error {
	// Find an available license
	var license models.NILicense
	if err := tx.Where("is_assigned = ? AND product_name = ?", false, reward.Value).
		First(&license).Error; err != nil {
		return errors.New("no available licenses found for this product")
	}

	// Assign license to user
	now := time.Now()
	updates := map[string]interface{}{
		"is_assigned": true,
		"assigned_to": reward.UserID,
		"assigned_at": &now,
	}

	if err := tx.Model(&license).Updates(updates).Error; err != nil {
		return err
	}

	// Update reward with license key
	return tx.Model(reward).Updates(map[string]interface{}{
		"value":        license.LicenseKey,
		"status":       models.RewardStatusDelivered,
		"delivered_at": &now,
	}).Error
}

func (s *RewardService) addPointsToUser(tx *gorm.DB, reward *models.Reward) error {
	// This would integrate with your points system
	// For now, just mark as delivered
	now := time.Now()
	return tx.Model(reward).Updates(map[string]interface{}{
		"status":       models.RewardStatusDelivered,
		"delivered_at": &now,
	}).Error
}

func (s *RewardService) assignBadgeToUser(tx *gorm.DB, reward *models.Reward) error {
	// This would integrate with your badge system
	// For now, just mark as delivered
	now := time.Now()
	return tx.Model(reward).Updates(map[string]interface{}{
		"status":       models.RewardStatusDelivered,
		"delivered_at": &now,
	}).Error
}

func (s *RewardService) logAdminAction(tx *gorm.DB, adminID uuid.UUID, action string, targetID *uuid.UUID, targetType, details string) error {
	adminAction := models.AdminAction{
		AdminID:    adminID,
		Action:     action,
		TargetID:   targetID,
		TargetType: targetType,
		Details:    details,
	}

	return tx.Create(&adminAction).Error
}
