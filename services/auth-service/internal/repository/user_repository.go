package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"iot-platform/services/auth-service/internal/domain"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByUsername(ctx context.Context, tenantID, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, tenantID, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, tenantID string, offset, limit int) ([]*domain.User, int64, error)
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

func (r *userRepository) GetByEmail(ctx context.Context, tenantID, email string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, username, email, password_hash, full_name, phone, status, last_login_at, metadata, created_at, updated_at
		FROM users
		WHERE tenant_id = $1 AND email = $2
	`

	user := &domain.User{}
	var metadata []byte
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, tenantID, email).Scan(
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

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	metadata, _ := json.Marshal(user.Metadata)

	query := `
		UPDATE users
		SET full_name = $1, phone = $2, status = $3, metadata = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.ExecContext(ctx, query,
		user.FullName, user.Phone, user.Status, metadata, user.UpdatedAt, user.ID,
	)

	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) List(ctx context.Context, tenantID string, offset, limit int) ([]*domain.User, int64, error) {
	// 查询总数
	var total int64
	countQuery := `SELECT COUNT(*) FROM users WHERE tenant_id = $1`
	err := r.db.QueryRowContext(ctx, countQuery, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	query := `
		SELECT id, tenant_id, username, email, full_name, phone, status, last_login_at, metadata, created_at, updated_at
		FROM users
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		user := &domain.User{}
		var metadata []byte
		var lastLoginAt sql.NullTime

		err := rows.Scan(
			&user.ID, &user.TenantID, &user.Username, &user.Email,
			&user.FullName, &user.Phone, &user.Status, &lastLoginAt, &metadata,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		}

		if len(metadata) > 0 {
			json.Unmarshal(metadata, &user.Metadata)
		}

		users = append(users, user)
	}

	return users, total, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	query := `UPDATE users SET last_login_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}
