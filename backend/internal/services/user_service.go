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
