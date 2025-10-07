package services

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/uploadparty/app/internal/models"
)

type UserService struct {
	DB        *gorm.DB
	JWTSecret string
}

func NewUserService(db *gorm.DB, secret string) *UserService {
	return &UserService{DB: db, JWTSecret: secret}
}

func (s *UserService) Register(email, username, password string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.ToLower(strings.TrimSpace(username))
	if email == "" || username == "" || len(password) < 6 {
		return nil, errors.New("invalid input")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &models.User{Email: email, Username: username, PasswordHash: string(hash), DisplayName: username, Public: true}
	if err := s.DB.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) Authenticate(emailOrUsername, password string) (string, *models.User, error) {
	var u models.User
	q := s.DB.Where("email = ?", strings.ToLower(emailOrUsername)).Or("username = ?", strings.ToLower(emailOrUsername))
	if err := q.First(&u).Error; err != nil {
		return "", nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", nil, errors.New("invalid credentials")
	}
	claims := jwt.MapClaims{
		"sub": u.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", nil, err
	}
	return signed, &u, nil
}

func (s *UserService) FindPublicByHandle(handle string) (*models.User, error) {
	var u models.User
	if err := s.DB.Where("username = ? AND public = ?", strings.ToLower(handle), true).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// SyncAuth0User creates or updates a user from Auth0 data
func (s *UserService) SyncAuth0User(auth0ID, email, username, displayName, picture string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.ToLower(strings.TrimSpace(username))

	if auth0ID == "" || email == "" || username == "" {
		return nil, errors.New("auth0_id, email, and username are required")
	}

	var user models.User

	// Try to find existing user by Auth0 ID
	err := s.DB.Where("auth0_id = ?", auth0ID).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		// User doesn't exist, create new one
		user = models.User{
			Auth0ID:     auth0ID,
			Email:       email,
			Username:    username,
			DisplayName: displayName,
			Picture:     picture,
			Public:      true,
		}

		if err := s.DB.Create(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	} else if err != nil {
		return nil, err
	}

	// User exists, update profile info
	updates := map[string]interface{}{
		"email":        email,
		"display_name": displayName,
		"picture":      picture,
	}

	if err := s.DB.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByAuth0ID finds a user by their Auth0 ID
func (s *UserService) FindByAuth0ID(auth0ID string) (*models.User, error) {
	var u models.User
	if err := s.DB.Where("auth0_id = ?", auth0ID).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
