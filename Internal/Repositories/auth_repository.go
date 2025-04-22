package Repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/berik/GoTune/Domain/Entities"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *Entities.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, address, phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx, query, user.Username, user.Email, user.PasswordHash, user.Address, user.Phone,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create instrument: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id int) (*Entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, address, phone
		FROM users
		WHERE id = $1
	`

	var user Entities.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Address, &user.Phone,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("instrument not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get instrument: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByUsername(ctx context.Context, username string) (*Entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, address, phone
		FROM users
		WHERE username = $1
	`

	var user Entities.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Address, &user.Phone,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("instrument not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get instrument by username: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*Entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, address, phone
		FROM users
		WHERE email = $1
	`

	var user Entities.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Address, &user.Phone,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("instrument not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get instrument by email: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *Entities.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3, address = $4, phone = $5
		WHERE id = $6
	`

	result, err := r.db.ExecContext(
		ctx, query, user.Username, user.Email, user.PasswordHash, user.Address, user.Phone, user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update instrument: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("instrument with id %d not found", user.ID)
	}

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete instrument: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("instrument with id %d not found", id)
	}

	return nil
}

func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*Entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, address, phone
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*Entities.User
	for rows.Next() {
		var user Entities.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Address, &user.Phone,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan instrument: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users rows: %w", err)
	}

	return users, nil
}
