package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
)

type JWTConfig struct {
	SecretKey      string
	AccessTokenTTL time.Duration
}

var ErrInvalidCredentials = errors.New("неправильный email или пароль")

type UserClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwtCfg   JWTConfig
}

func NewAuthService(
	userRepo repository.UserRepository,
	jwtCfg JWTConfig,
) AuthService {
	return &authService{userRepo: userRepo, jwtCfg: jwtCfg}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := checkPassword(user.Password, password); err != nil {
		return "", ErrInvalidCredentials
	}

	now := time.Now()
	claims := UserClaims{
		UserID: user.ID,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtCfg.AccessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signed, err := token.SignedString([]byte(s.jwtCfg.SecretKey))

	if err != nil {
		return "", err
	}

	return signed, nil
}
