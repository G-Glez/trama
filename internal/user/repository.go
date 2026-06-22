package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	usergen "trama/internal/gen/user"
)

type userQuerier interface {
	CreateUser(ctx context.Context, arg usergen.CreateUserParams) error
	GetUser(ctx context.Context, id string) (usergen.User, error)
	GetUserByEmail(ctx context.Context, email string) (usergen.User, error)
	ListUsers(ctx context.Context) ([]usergen.User, error)
	UpdateUser(ctx context.Context, arg usergen.UpdateUserParams) error
	DeleteUser(ctx context.Context, id string) error
}

type UserSQLRepository struct {
	querier userQuerier
}

func NewUserRepository(querier userQuerier) *UserSQLRepository {
	return &UserSQLRepository{querier: querier}
}

func (r *UserSQLRepository) Create(ctx context.Context, u User) (User, error) {
	u.ID = UserID{uuid.New()}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	err := r.querier.CreateUser(ctx, usergen.CreateUserParams{
		ID:             u.ID.String(),
		Username:       u.Username,
		Email:          u.Email,
		PasswordHash:   u.PasswordHash,
		DefaultFaction: toNullString(u.DefaultFaction),
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	})
	if err != nil {
		return User{}, fmt.Errorf("database error: %w", err)
	}

	return u, nil
}

func (r *UserSQLRepository) Get(ctx context.Context, id UserID) (User, error) {
	u, err := r.querier.GetUser(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("database error: %w", err)
	}

	return toDomainUser(u), nil
}

func (r *UserSQLRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	u, err := r.querier.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("database error: %w", err)
	}

	return toDomainUser(u), nil
}

func (r *UserSQLRepository) List(ctx context.Context) ([]User, error) {
	items, err := r.querier.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	list := make([]User, len(items))
	for i, u := range items {
		list[i] = toDomainUser(u)
	}
	return list, nil
}

func (r *UserSQLRepository) Update(ctx context.Context, u User) error {
	err := r.querier.UpdateUser(ctx, usergen.UpdateUserParams{
		Username:       u.Username,
		Email:          u.Email,
		PasswordHash:   u.PasswordHash,
		DefaultFaction: toNullString(u.DefaultFaction),
		UpdatedAt:      u.UpdatedAt,
		ID:             u.ID.String(),
	})
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (r *UserSQLRepository) Delete(ctx context.Context, id UserID) error {
	err := r.querier.DeleteUser(ctx, id.String())
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func toDomainUser(u usergen.User) User {
	uid := uuid.MustParse(u.ID)
	var defFac uuid.UUID
	if u.DefaultFaction.Valid {
		defFac = uuid.MustParse(u.DefaultFaction.String)
	}

	return User{
		ID:             UserID{uid},
		Username:       u.Username,
		Email:          u.Email,
		PasswordHash:   u.PasswordHash,
		DefaultFaction: defFac,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

func toNullString(f uuid.UUID) sql.NullString {
	if f == uuid.Nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: f.String(), Valid: true}
}
