package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"service-app/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, input user.CreateUserInput) (user.User, error) {
	const q = `
INSERT INTO users (name, email)
VALUES ($1, $2)
RETURNING id, name, email, created_at, updated_at`

	var out user.User
	if err := r.db.QueryRowContext(ctx, q, input.Name, input.Email).Scan(&out.ID, &out.Name, &out.Email, &out.CreatedAt, &out.UpdatedAt); err != nil {
		return user.User{}, fmt.Errorf("create user: %w", err)
	}
	return out, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (user.User, error) {
	const q = `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`

	var out user.User
	if err := r.db.QueryRowContext(ctx, q, id).Scan(&out.ID, &out.Name, &out.Email, &out.CreatedAt, &out.UpdatedAt); err != nil {
		return user.User{}, err
	}
	return out, nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int32) ([]user.User, error) {
	const q = `SELECT id, name, email, created_at, updated_at FROM users ORDER BY id DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]user.User, 0)
	for rows.Next() {
		var item user.User
		if err := rows.Scan(&item.ID, &item.Name, &item.Email, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, item)
	}
	return users, rows.Err()
}

func (r *UserRepository) Update(ctx context.Context, id int64, input user.UpdateUserInput) (user.User, error) {
	const q = `
UPDATE users
SET name = $2, email = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, name, email, created_at, updated_at`

	var out user.User
	if err := r.db.QueryRowContext(ctx, q, id, input.Name, input.Email).Scan(&out.ID, &out.Name, &out.Email, &out.CreatedAt, &out.UpdatedAt); err != nil {
		return user.User{}, err
	}
	return out, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM users WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("user not found")
	}
	return nil
}
