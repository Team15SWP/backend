package service

import (
	"context"
	"fmt"

	"study_buddy/internal/config"
	"study_buddy/internal/model"
	"study_buddy/pkg/constants"
	"study_buddy/pkg/errlist"
	"study_buddy/pkg/hash"

	"github.com/dgrijalva/jwt-go"
)

var _ Service = (*AuthService)(nil)

type AuthService struct {
	repo       UserProvider
	hashConfig *config.HashConfig
}

func NewAuthService(repo UserProvider, hashConfig *config.HashConfig) *AuthService {
	return &AuthService{
		repo:       repo,
		hashConfig: hashConfig,
	}
}

type Service interface {
	LogIn(ctx context.Context, email, password string) (string, error)
	SignUp(ctx context.Context, username, email, password string) (string, error)
}

type UserProvider interface {
	GetUserByEmailOrUsername(ctx context.Context, username string) (*model.UserData, error)
	CreateUser(ctx context.Context, username, email, password string) (*model.UserData, error)
}

func (a AuthService) LogIn(ctx context.Context, username, password string) (string, error) {
	user, err := a.repo.GetUserByEmailOrUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("[authService][LogIn][GetUserByEmailOrUsername]: %w", err)
	}
	if !hash.ComparePassword(password, user.Password) {
		return "", fmt.Errorf("[authService][LogIn][ComparePassword]: %w", errlist.ErrPasswordIsIncorrect)
	}

	tokenString, err := signToken(user, a.hashConfig.SigningKey)
	if err != nil {
		return "", fmt.Errorf("[authService][LogIn][SignToken]: %w", err)
	}
	return tokenString, nil
}

func (a AuthService) SignUp(ctx context.Context, username, email, password string) (string, error) {
	if _, err := a.repo.GetUserByEmailOrUsername(ctx, username); err == nil {
		return "", fmt.Errorf("[authService][SignUp][GetUserByEmailOrUsername]: %w", errlist.ErrUserExists)
	}
	if _, err := a.repo.GetUserByEmailOrUsername(ctx, email); err == nil {
		return "", fmt.Errorf("[authService][SignUp][GetUserByEmailOrUsername]: %w", errlist.ErrUserExists)
	}

	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("[authService][SignUp][HashPassword]: %w", err)
	}

	user, err := a.repo.CreateUser(ctx, username, email, hashedPassword)
	if err != nil {
		return "", fmt.Errorf("[authService][SignUp][CreateUser]: %w", err)
	}

	tokenString, err := signToken(user, a.hashConfig.SigningKey)
	if err != nil {
		return "", fmt.Errorf("[authService][SignUp][SignToken]: %w", err)
	}
	return tokenString, nil
}

func signToken(user *model.UserData, signingKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constants.UserID: user.ID,
		constants.Role:   user.Role,
		constants.Name:   user.Name,
		constants.Email:  user.Email,
	})

	secretKey := []byte(signingKey)
	tokenString, err := token.SignedString(secretKey)
	return tokenString, err
}
