package interfaces

import (
	"context"
	"evermos-hasby/ecommerce/internal/models"
)

type IUserRepository interface {
	ListUsers(ctx context.Context) ([]models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserById(ctx context.Context, id uint) (*models.User, error)
	UpdateUserProfile(ctx context.Context, user *models.User) error
	UpdateManagedUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uint) error
}

type IAccountService interface {
	List(ctx context.Context, adminID uint) ([]models.User, error)
	Get(ctx context.Context, adminID, id uint) (*models.User, error)
	Create(ctx context.Context, adminID uint, user *models.User) error
	Update(ctx context.Context, adminID, id uint, user *models.User) (*models.User, error)
	Delete(ctx context.Context, adminID, id uint) error
}

type IUserService interface {
	Logout(ctx context.Context, token string) error
	Register(ctx context.Context, request *models.User) (models.UserResponse, error)
	Login(ctx context.Context, request models.LoginRequest) (models.UserResponse, error)
	Refresh(ctx context.Context, refreshToken string, claims models.TokenClaims) (models.RefreshTokenResponse, error)
	GetProfile(ctx context.Context, userID uint) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uint, request *models.User) (*models.User, error)
}
