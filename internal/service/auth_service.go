package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gitsang/order/internal/config"
	"github.com/gitsang/order/internal/model"
	"github.com/gitsang/order/internal/repository"
)

type AuthService struct {
	jwtConfig config.JWTConfig
	userRepo  *repository.UserRepository
}

func NewAuthService(jwtConfig config.JWTConfig, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		jwtConfig: jwtConfig,
		userRepo:  userRepo,
	}
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Duration(s.jwtConfig.Expiration) * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) Register(username, password, name, phone string) (*model.User, error) {
	existingUser, _ := s.userRepo.FindByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	if phone != "" {
		existingUser, _ = s.userRepo.FindByPhone(phone)
		if existingUser != nil {
			return nil, errors.New("phone already exists")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Name:     name,
		Phone:    phone,
		Role:     "customer",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) ValidateToken(tokenString string) (uuid.UUID, string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil {
		return uuid.Nil, "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, "", "", errors.New("invalid token")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, "", "", errors.New("invalid user_id claim")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, "", "", err
	}

	username, _ := claims["username"].(string)
	role, _ := claims["role"].(string)

	return userID, username, role, nil
}
