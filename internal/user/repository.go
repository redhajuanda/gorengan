package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/redhajuanda/gorengan/internal/domain"
)

// Repository encapsulates the logic to access users from the data source.
type Repository interface {
	// Get returns the user with the specified user ID.
	Get(ctx context.Context, id string) (domain.User, error)
	// Count returns the number of users.
	Count(ctx context.Context) (int, error)
	// Query returns the list of users with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]domain.User, error)
	// Create saves a new user in the storage.
	Create(ctx context.Context, user domain.User) error
	// Update updates the user with given ID in the storage.
	Update(ctx context.Context, user domain.User) error
	// Delete removes the user with given ID from the storage.
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new user repository
func NewRepository(db *sql.DB) Repository {
	return repository{db}
}

// Get returns the user with the specified user ID.
func (r repository) Get(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, first_name, last_name, email, password, address, created_at, updated_at FROM users WHERE id=?")
	if err != nil {
		return domain.User{}, err
	}
	row := stmt.QueryRowContext(ctx, id)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Address, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// Count returns the number of users.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	stmt, err := r.db.PrepareContext(ctx, "SELECT COUNT(*) as count FROM users")
	if err != nil {
		return 0, err
	}
	rows := stmt.QueryRowContext(ctx)
	rows.Scan(&count)
	return count, nil
}

// Query returns the list of users with the given offset and limit.
func (r repository) Query(ctx context.Context, offset, limit int) ([]domain.User, error) {
	var users []domain.User
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, first_name, last_name, email, password, address, created_at, updated_at FROM users LIMIT ?, ?")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.QueryContext(ctx, offset, limit)
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Address, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Create saves a new user in the storage.
func (r repository) Create(ctx context.Context, user domain.User) error {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users (id, first_name, last_name, email, password, address, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("Error preparing statement: %v", err)
	}
	_, err = stmt.ExecContext(ctx, user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.Address, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("Error exec query: %v", err)
	}
	return nil
}

// Update updates the user with given ID in the storage.
func (r repository) Update(ctx context.Context, user domain.User) error {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE users SET first_name=?, last_name=?, email=?, password=?, address=?, created_at=?, updated_at=? WHERE id=?")
	if err != nil {
		return fmt.Errorf("Error preparing statement: %v", err)
	}
	_, err = stmt.ExecContext(ctx, user.FirstName, user.LastName, user.Email, user.Password, user.Address, user.CreatedAt, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("Error exec query: %v", err)
	}
	return nil
}

// Delete removes the user with given ID from the storage.
func (r repository) Delete(ctx context.Context, id string) error {
	stmt, err := r.db.PrepareContext(ctx, "DELETE FROM users WHERE id=?")
	if err != nil {
		return fmt.Errorf("Error preparing statement: %v", err)
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("Error exec query: %v", err)
	}
	return nil
}
