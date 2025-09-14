package services

import (
	"errors"
	"time"
	"uploadparty/internal/middlewares"
	"uploadparty/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db        *gorm.DB
	jwtSecret string
}

type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=30"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,max=50"`
	LastName  string `json:"last_name" validate:"required,max=50"`
	Age       int    `json:"age" validate:"required,min=13,max=19"`
	School    string `json:"school" validate:"required,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email or username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := models.User{
		ID:              uuid.New(),
		Username:        req.Username,
		Email:           req.Email,
		PasswordHash:    string(hashedPassword),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Age:             req.Age,
		School:          req.School,
		IsVerified:      false,
		IsActive:        true,
		TotalPoints:     0,
		Level:           1,
		NILicenseStatus: "pending",
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.generateJWT(user.ID.String(), user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	// Hide password hash
	user.PasswordHash = ""

	return &AuthResponse{
		User:  &user,
		Token: token,
	}, nil
}

func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Generate JWT token
	token, err := s.generateJWT(user.ID.String(), user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	// Hide password hash
	user.PasswordHash = ""

	return &AuthResponse{
		User:  &user,
		Token: token,
	}, nil
}

func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// Hide password hash
	user.PasswordHash = ""
	return &user, nil
}

func (s *AuthService) generateJWT(userID, username, email string) (string, error) {
	claims := middlewares.JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
