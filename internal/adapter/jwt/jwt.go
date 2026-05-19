package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const AuthHeader = "Authorization"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type Config struct {
	Secret        string `yaml:"secret"`
	ExpireMinutes int    `yaml:"expire_minutes"`
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

var globalConfig *Config

func Init(cfg *Config) {
	globalConfig = cfg
}

func GenerateToken(userID int64, email, role string) (string, error) {
	if globalConfig == nil {
		return "", errors.New("jwt config is not initialized")
	}

	expirationTime := time.Now().Add(time.Duration(globalConfig.ExpireMinutes) * time.Minute)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(globalConfig.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	if globalConfig == nil {
		return nil, errors.New("jwt config is not initialized")
	}

	claims := &Claims{UserID: 0, Role: "", Email: "", RegisteredClaims: jwt.RegisteredClaims{}}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(globalConfig.Secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func ExtractClaims(tokenString string) (*Claims, error) {
	if globalConfig == nil {
		return nil, errors.New("jwt config is not initialized")
	}

	if after, ok := strings.CutPrefix(tokenString, "Bearer "); ok {
		tokenString = after
	}

	claims := &Claims{UserID: 0, Role: "", Email: "", RegisteredClaims: jwt.RegisteredClaims{}}

	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	if _, ok := token.Claims.(*Claims); !ok {
		return nil, fmt.Errorf("%w: failed to cast claims", ErrInvalidToken)
	}

	return claims, nil
}
