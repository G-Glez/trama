package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserExists      = errors.New("username already exists")
)

type TokenProvider interface {
	Generate(claims Claims) (string, error)
	Validate(tokenString string) (Claims, error)
}

type UserRepository interface {
	Create(ctx context.Context, user User) error
	GetByUsername(ctx context.Context, username string) (User, error)
}

// -----------------------------------------------------------------------------------
// Service provides user registration, login, and token validation.
// -----------------------------------------------------------------------------------
type Service struct {
	users         UserRepository
	tokenProvider TokenProvider
}

func NewService(userRepository UserRepository, tokens TokenProvider) *Service {
	return &Service{
		users:         userRepository,
		tokenProvider: tokens,
	}
}

// -----------------------------------------------------------------------------------
// RegisterInput contains the credentials for creating a new account.
// -----------------------------------------------------------------------------------
type RegisterInput struct {
	Username string
	Password string
}

// -----------------------------------------------------------------------------------
// Register creates a new user and persists it via the UserRepository.
// OnError: ErrUserExists
// -----------------------------------------------------------------------------------
func (s *Service) Register(ctx context.Context, input RegisterInput) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user := User{
		Username:     input.Username,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
	}

	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, ErrUserExists) {
			return ErrUserExists
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// -----------------------------------------------------------------------------------
type LoginInput struct {
	Username string
	Password string
}
type LoginOutput struct {
	Token string
}

// -----------------------------------------------------------------------------------
// Login authenticates a user by fetching the stored hash from DynamoDB and comparing
// it with the provided password. On success it generates a JWT and returns it.
// Output: LoginResponse with token and user data.
// OnError: ErrUserNotFound, ErrInvalidPassword
// -----------------------------------------------------------------------------------
func (s *Service) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	user, err := s.users.GetByUsername(ctx, input.Username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return LoginOutput{}, ErrUserNotFound
		}
		return LoginOutput{}, fmt.Errorf("get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return LoginOutput{}, ErrInvalidPassword
	}

	claims := claimsFromUser(user)
	token, err := s.tokenProvider.Generate(claims)
	if err != nil {
		return LoginOutput{}, err
	}

	return LoginOutput{Token: token}, nil
}

// -----------------------------------------------------------------------------------
// ValidateToken parses and validates a JWT token string.
// Output: Claims
// -----------------------------------------------------------------------------------
func (s *Service) ValidateToken(tokenString string) (Claims, error) {
	return s.tokenProvider.Validate(tokenString)
}

func claimsFromUser(user User) Claims {
	return Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Username,
		},
	}
}
