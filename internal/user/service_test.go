// +build all service

package user

import (
	"context"
	"database/sql"
	"testing"

	"github.com/redhajuanda/gorengan/internal/domain"
	"github.com/redhajuanda/gorengan/pkg/log"
	"github.com/stretchr/testify/assert"
)

var serviceTest Service

func createNewServiceTest(t *testing.T) Service {
	if serviceTest != nil {
		return serviceTest
	}
	logger, _ := log.NewForTest()
	serviceTest = NewService(&mockRepository{}, logger)
	return serviceTest
}

func TestServiceCreateUser(t *testing.T) {
	service := createNewServiceTest(t)

	var inputRequests = []CreateUserRequest{
		{
			FirstName: "Redha",
			LastName:  "Redha",
			Email:     "Redha@sdfdsf.vo",
			Password:  "Redha",
			Address:   "Redha",
		},
	}

	for _, inputRequest := range inputRequests {
		_, err := service.Create(context.Background(), inputRequest)
		assert.NoError(t, err)
	}
}

func TestServiceGetUser(t *testing.T) {
	service := createNewServiceTest(t)

	users, err := service.Query(context.Background(), 0, 0)
	assert.NoError(t, err)

	user, err := service.Get(context.Background(), users[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, user, users[0])
}

func TestServiceQueryUser(t *testing.T) {
	service := createNewServiceTest(t)

	users, err := service.Query(context.Background(), 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users))
}

func TestServiceCountUser(t *testing.T) {
	service := createNewServiceTest(t)

	count, err := service.Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestServiceDeleteUser(t *testing.T) {
	service := createNewServiceTest(t)
	users, err := service.Query(context.Background(), 0, 0)
	assert.NoError(t, err)

	_, err = service.Delete(context.Background(), users[0].ID)
	assert.NoError(t, err)

	count, err := service.Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

type mockRepository struct {
	users []domain.User
}

// Get returns the user with the specified user ID.
func (m mockRepository) Get(ctx context.Context, id string) (domain.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return domain.User{}, sql.ErrNoRows
}

// Count returns the number of users.
func (m mockRepository) Count(ctx context.Context) (int, error) {
	return len(m.users), nil
}

// Query returns the list of users with the given offset and limit.
func (m mockRepository) Query(ctx context.Context, offset, limit int) ([]domain.User, error) {
	return m.users, nil
}

// Create saves a new user in the storage.
func (m *mockRepository) Create(ctx context.Context, user domain.User) error {
	m.users = append(m.users, user)
	return nil
}

// Update updates the user with given ID in the storage.
func (m *mockRepository) Update(ctx context.Context, user domain.User) error {
	for i, item := range m.users {
		if item.ID == user.ID {
			m.users[i] = user
			break
		}
	}
	return nil
}

// Delete removes the user with given ID from the storage.
func (m *mockRepository) Delete(ctx context.Context, id string) error {
	for i, user := range m.users {
		if user.ID == id {
			m.users[i] = m.users[len(m.users)-1]
			m.users = m.users[:len(m.users)-1]
			break
		}
	}
	return nil
}
