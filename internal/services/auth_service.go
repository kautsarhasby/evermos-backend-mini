package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"evermos-hasby/ecommerce/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	UserRepository interfaces.IUserRepository
	DB             *gorm.DB
	Redis          *redis.Client
}

func sessionKey(userID uint, sessionID string) string {
	return fmt.Sprintf("session:%d:%s", userID, sessionID)
}

func (s *AuthService) Login(ctx context.Context, request models.LoginRequest) (models.UserResponse, error) {
	var response models.UserResponse
	if s.Redis == nil {
		return response, errors.New("Redis is not configured")
	}

	user, err := s.UserRepository.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return models.UserResponse{}, errors.Wrap(err, "Failed to load user")
	}
	if user == nil {
		return models.UserResponse{}, errors.New("Username or password invalid")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return models.UserResponse{}, errors.Wrap(err, "Username or password invalid")
	}

	sessionID := uuid.NewString()
	token, err := utils.GenerateTokenWithID(ctx, user.ID, user.Username, user.Role, "token", sessionID)
	if err != nil {
		return response, errors.Wrap(err, "failed to generate token")
	}
	refreshToken, err := utils.GenerateTokenWithID(ctx, user.ID, user.Username, user.Role, "refreshToken", sessionID)
	if err != nil {
		return response, errors.Wrap(err, "failed to generate refresh token")
	}

	session := models.RedisSession{UserID: user.ID, Token: token, RefreshToken: refreshToken}
	data, err := json.Marshal(session)
	if err != nil {
		return response, err
	}
	if err := s.Redis.Set(ctx, sessionKey(user.ID, sessionID), data, utils.TokenType["refreshToken"]).Err(); err != nil {
		return response, errors.Wrap(err, "Failed to create new session")
	}

	response.UserID = user.ID
	response.Username = user.Username
	response.Role = user.Role
	response.Token = token
	response.RefreshToken = refreshToken

	return response, nil
}

func (s *AuthService) Register(ctx context.Context, request *models.User) (models.UserResponse, error) {
	if request.Username == "" || request.Email == "" || request.Phone == "" || request.Password == "" {
		return models.UserResponse{}, errors.New("Username, email, phone, and password are required")
	}
	user, err := s.UserRepository.GetUserByUsername(ctx, request.Username)
	if err != nil {
		utils.Logger.WithError(err).WithField("username", request.Username).Error("Failed to check existing user")
		return models.UserResponse{}, errors.Wrap(err, "Failed to load user")
	}
	if user != nil {
		utils.Logger.WithField("username", request.Username).Warn("Registration failed: user already exists")
		return models.UserResponse{}, errors.New("User Already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Logger.WithError(err).Error("Critical: Failed to hash password")
		return models.UserResponse{}, errors.Wrap(err, "Failed to hash password")
	}

	userModel := &models.User{Username: request.Username,
		Email:       request.Email,
		Password:    string(hashedPassword),
		Phone:       request.Phone,
		DateOfBirth: request.DateOfBirth,
		Gender:      request.Gender,
		Occupation:  request.Occupation,
		ProvinceID:  request.ProvinceID,
		CityID:      request.CityID,
		Role:        "user"}

	err = s.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(userModel).Error; err != nil {
			return err
		}

		storeModel := &models.Store{
			Name:   userModel.Username + " Store",
			UserID: userModel.ID,
		}

		if err := tx.Create(storeModel).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return models.UserResponse{}, errors.Wrap(err, "Failed to Create User")
	}
	return models.UserResponse{Username: userModel.Username, Role: userModel.Role}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string, claims models.TokenClaims) (models.RefreshTokenResponse, error) {
	if s.Redis == nil {
		return models.RefreshTokenResponse{}, errors.New("Redis is not configured")
	}
	key := sessionKey(claims.UserID, claims.ID)
	var session models.RedisSession
	data, err := s.Redis.Get(ctx, key).Bytes()
	if err != nil || json.Unmarshal(data, &session) != nil || session.RefreshToken != refreshToken {
		return models.RefreshTokenResponse{}, errors.New("Invalid refresh session")
	}
	token, err := utils.GenerateTokenWithID(ctx, claims.UserID, claims.Username, claims.Role, "token", claims.ID)
	if err != nil {
		return models.RefreshTokenResponse{}, errors.Wrap(err, "Failed to generate token")
	}
	session.Token = token
	data, err = json.Marshal(session)
	if err != nil {
		return models.RefreshTokenResponse{}, err
	}
	if claims.ExpiresAt == nil {
		return models.RefreshTokenResponse{}, errors.New("Refresh token expiry is missing")
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return models.RefreshTokenResponse{}, errors.New("Refresh token expired")
	}
	if err := s.Redis.Set(ctx, key, data, ttl).Err(); err != nil {
		return models.RefreshTokenResponse{}, errors.Wrap(err, "failed to update session")
	}
	return models.RefreshTokenResponse{Token: token}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	if s.Redis == nil {
		return errors.New("Redis is not configured")
	}
	claims, err := utils.ValidateToken(ctx, token)
	if err != nil {
		return err
	}
	return s.Redis.Del(ctx, sessionKey(claims.UserID, claims.ID)).Err()
}

func (s *AuthService) GetProfile(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.UserRepository.GetUserById(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gorm.ErrRecordNotFound
	}
	user.Password = ""
	return user, nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uint, request *models.User) (*models.User, error) {
	user, err := s.UserRepository.GetUserById(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gorm.ErrRecordNotFound
	}
	if request.Username != "" {
		user.Username = request.Username
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Phone != "" {
		user.Phone = request.Phone
	}
	if !request.DateOfBirth.IsZero() {
		user.DateOfBirth = request.DateOfBirth
	}
	if request.Gender != "" {
		user.Gender = request.Gender
	}
	if request.Occupation != "" {
		user.Occupation = request.Occupation
	}
	if request.ProvinceID != "" {
		user.ProvinceID = request.ProvinceID
	}
	if request.CityID != "" {
		user.CityID = request.CityID
	}
	if err := s.UserRepository.UpdateUserProfile(ctx, user); err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}
