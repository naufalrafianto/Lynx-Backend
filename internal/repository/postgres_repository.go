package repository

import (
	"auth-service/internal/cache"
	"auth-service/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type cachedUserRepository struct {
	db    *sql.DB
	cache cache.CacheService
}

func NewCachedUserRepository(db *sql.DB, cache cache.CacheService) domain.UserRepository {
	return &cachedUserRepository{
		db:    db,
		cache: cache,
	}
}

const (
	userCacheDuration = 15 * time.Minute
	userByIDKey       = "user:id:%s"
	userByEmailKey    = "user:email:%s"
)

func (r *cachedUserRepository) GetByID(id string) (*domain.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf(userByIDKey, id)

	var user domain.User
	err := r.cache.GetOrSet(ctx, cacheKey, &user, userCacheDuration, func() (interface{}, error) {
		user := &domain.User{}
		query := `
            SELECT id, email, name, is_active, created_at, updated_at
            FROM users 
            WHERE id = $1
        `
		err := r.db.QueryRow(query, id).Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, domain.ErrUserNotFound
			}
			return nil, fmt.Errorf("database error: %w", err)
		}
		return user, nil
	})

	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *cachedUserRepository) GetByEmail(email string) (*domain.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf(userByEmailKey, email)
	var user domain.User
	err := r.cache.GetOrSet(ctx, cacheKey, &user, userCacheDuration, func() (interface{}, error) {
		user := &domain.User{}
		fmt.Println(user)

		query := `
            SELECT id, email, name, password, is_active, created_at, updated_at
            FROM users 
            WHERE email = $1
        `
		err := r.db.QueryRow(query, email).Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Password,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, domain.ErrUserNotFound
			}
			return nil, fmt.Errorf("database error: %w", err)
		}
		fmt.Println(user)
		return user, nil
	})

	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *cachedUserRepository) Create(user *domain.User) error {
	ctx := context.Background()

	// Check if user already exists
	_, err := r.GetByEmail(user.Email)
	if err == nil {
		return domain.ErrUserAlreadyExists
	}
	if err != domain.ErrUserNotFound {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	query := `
        INSERT INTO users (id, email, name, password, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `

	if err := r.db.QueryRow(
		query,
		user.ID,
		user.Email,
		user.Name,
		user.Password,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Invalidate any existing cache entries
	_ = r.cache.Delete(ctx, fmt.Sprintf(userByIDKey, user.ID))
	_ = r.cache.Delete(ctx, fmt.Sprintf(userByEmailKey, user.Email))

	return nil
}

func (r *cachedUserRepository) UpdateActive(id string, isActive bool) error {
	ctx := context.Background()

	query := `
        UPDATE users 
        SET is_active = $1, updated_at = CURRENT_TIMESTAMP 
        WHERE id = $2
        RETURNING email
    `

	var email string
	err := r.db.QueryRow(query, isActive, id).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Invalidate cache entries
	_ = r.cache.Delete(ctx, fmt.Sprintf(userByIDKey, id))
	_ = r.cache.Delete(ctx, fmt.Sprintf(userByEmailKey, email))

	return nil
}

func (r *cachedUserRepository) SoftDelete(id string) error {
	ctx := context.Background()

	query := `
        UPDATE users 
        SET is_deleted = true,
            deleted_at = CURRENT_TIMESTAMP,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1 AND NOT is_deleted
        RETURNING email
    `

	var email string
	err := r.db.QueryRow(query, id).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("failed to soft delete user: %w", err)
	}

	// Invalidate cache entries
	_ = r.cache.Delete(ctx, fmt.Sprintf(userByIDKey, id))
	_ = r.cache.Delete(ctx, fmt.Sprintf(userByEmailKey, email))

	return nil
}

func (r *cachedUserRepository) HardDelete(id string) error {
	ctx := context.Background()

	// First get the email for cache invalidation
	var email string
	err := r.db.QueryRow("SELECT email FROM users WHERE id = $1", id).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("failed to get user email: %w", err)
	}

	// Then delete the user
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	// Invalidate cache entries
	_ = r.cache.Delete(ctx, fmt.Sprintf(userByIDKey, id))
	_ = r.cache.Delete(ctx, fmt.Sprintf(userByEmailKey, email))

	return nil
}
