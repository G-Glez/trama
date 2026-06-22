package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	Create(ctx context.Context, u User) (User, error)
	Get(ctx context.Context, id UserID) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	List(ctx context.Context) ([]User, error)
	Update(ctx context.Context, u User) error
	Delete(ctx context.Context, id UserID) error
}

type UserService struct {
	repo userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, in RegisterInput) (UserOutput, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))

	if !strings.Contains(email, "@") {
		return UserOutput{}, ErrInvalidEmail
	}

	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return UserOutput{}, err
	}
	if err == nil && existing.Email != "" {
		return UserOutput{}, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserOutput{}, fmt.Errorf("hash error: %w", err)
	}

	u := User{
		Username:       in.Username,
		Email:          email,
		PasswordHash:   string(hash),
		DefaultFaction: in.DefaultFaction,
	}

	created, err := s.repo.Create(ctx, u)
	if err != nil {
		return UserOutput{}, err
	}

	return toUserOutput(created), nil
}

func (s *UserService) Login(ctx context.Context, in LoginInput) (UserOutput, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))

	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return UserOutput{}, fmt.Errorf("user %s: %w", email, ErrUserNotFound)
		}
		return UserOutput{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return UserOutput{}, ErrInvalidPassword
	}

	return toUserOutput(u), nil
}

func (s *UserService) Get(ctx context.Context, id uuid.UUID) (UserOutput, error) {
	u, err := s.repo.Get(ctx, UserID{id})
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return UserOutput{}, fmt.Errorf("user %s: %w", id, ErrUserNotFound)
		}
		return UserOutput{}, err
	}

	return toUserOutput(u), nil
}

func (s *UserService) Update(ctx context.Context, in UpdateUserInput) error {
	existing, err := s.repo.Get(ctx, UserID{in.ID})
	if err != nil {
		return err
	}

	u := User{
		ID:           existing.ID,
		Username:     in.Username,
		Email:        existing.Email,
		PasswordHash: existing.PasswordHash,
		UpdatedAt:    time.Now(),
	}

	if in.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash error: %w", err)
		}
		u.PasswordHash = string(hash)
	}

	u.DefaultFaction = in.DefaultFaction

	return s.repo.Update(ctx, u)
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, UserID{id})
}

func toUserOutput(u User) UserOutput {
	return UserOutput{
		ID:             u.ID.UUID,
		Username:       u.Username,
		Email:          u.Email,
		DefaultFaction: u.DefaultFaction,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}
