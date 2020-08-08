package auth

import (
	"context"
	"database/sql"

	"github.com/redhajuanda/gorengan/internal/domain"
)

// Repository encapsulates the logic to access users from the data source.
type Repository interface {
	// Get returns the user with the specified user ID.
	Login(ctx context.Context, email string) (domain.User, error)
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *sql.DB) Repository {
	return repository{db}
}

// Get returns the user with the specified user ID.
func (r repository) Login(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, first_name, last_name, email, password, address, created_at, updated_at FROM users WHERE email=?")
	if err != nil {
		return domain.User{}, err
	}
	row := stmt.QueryRowContext(ctx, email)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Address, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return domain.User{}, err
	}
	return user, nil
}
