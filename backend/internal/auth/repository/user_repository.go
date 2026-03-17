package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"iot-platform/backend/internal/auth/domain"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByUsername(ctx context.Context, tenantID, username string) (*domain.User, error)
	UpdateLastLogin(ctx context.Context, userID string) error
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	metadata, _ := json.Marshal(user.Metadata)

	query := `
		INSERT INTO users (id, tenant_id, username, email, password_hash, full_name, phone, status, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.TenantID, user.Username, user.Email, user.Password,
		user.FullName, user.Phone, user.Status, metadata, user.CreatedAt, user.UpdatedAt,
	)

	return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, username, email, password_hash, full_name, phone, status, last_login_at, metadata, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	var metadata []byte
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.TenantID, &user.Username, &user.Email, &user.Password,
		&user.FullName, &user.Phone, &user.Status, &lastLoginAt, &metadata,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	if len(metadata) > 0 {
		json.Unmarshal(metadata, &user.Metadata)
	}

	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, tenantID, username string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, username, email, password_hash, full_name, phone, status, last_login_at, metadata, created_at, updated_at
		FROM users
		WHERE tenant_id = $1 AND username = $2
	`

	user := &domain.User{}
	var metadata []byte
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, tenantID, username).Scan(
		&user.ID, &user.TenantID, &user.Username, &user.Email, &user.Password,
		&user.FullName, &user.Phone, &user.Status, &lastLoginAt, &metadata,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	if len(metadata) > 0 {
		json.Unmarshal(metadata, &user.Metadata)
	}

	return user, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	query := `UPDATE users SET last_login_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}
