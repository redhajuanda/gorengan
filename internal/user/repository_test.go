// +build all repository

package user

import (
	"context"
	"testing"
	"time"

	"github.com/redhajuanda/gorengan/internal/domain"
	"github.com/redhajuanda/gorengan/internal/test"
	"github.com/stretchr/testify/assert"
)

var userDataTests = []domain.User{
	{
		ID:        domain.GenerateID(),
		FirstName: "Redha",
		LastName:  "Juanda",
		Email:     "redhajuanda@gmail.com",
		Password:  "password",
		Address:   "Address",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        domain.GenerateID(),
		FirstName: "John",
		LastName:  "Mick",
		Email:     "johnmick@gmail.com",
		Password:  "password",
		Address:   "Address",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

func TestCreateUser(t *testing.T) {
	db := test.GetTestDB(t)
	repo := NewRepository(db)

	for _, user := range userDataTests {
		err := repo.Create(context.Background(), user)
		assert.NoError(t, err)
	}
}

func TestGetOneUser(t *testing.T) {
	db := test.GetTestDB(t)
	repo := NewRepository(db)

	for _, user := range userDataTests {
		userGot, err := repo.Get(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.GetID(), userGot.GetID())
	}
}

func TestQueryUser(t *testing.T) {
	db := test.GetTestDB(t)
	repo := NewRepository(db)

	usersGot, err := repo.Query(context.Background(), 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, len(userDataTests), len(usersGot))
}

func TestUpdateUser(t *testing.T) {
	db := test.GetTestDB(t)
	repo := NewRepository(db)

	for _, user := range userDataTests {
		user.FirstName = "Update"
		err := repo.Update(context.Background(), user)
		assert.NoError(t, err)

		userGot, err := repo.Get(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.FirstName, userGot.FirstName)
	}
}

func TestCountUser(t *testing.T) {
	db := test.GetTestDB(t)
	repo := NewRepository(db)

	count, err := repo.Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, len(userDataTests), count)
}

func TestDeleteUser(t *testing.T) {
	db := test.GetTestDB(t)
	repo := NewRepository(db)

	err := repo.Delete(context.Background(), userDataTests[0].ID)
	assert.NoError(t, err)

	_, err = repo.Get(context.Background(), userDataTests[0].ID)
	assert.Error(t, err)

	count, err := repo.Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, len(userDataTests)-1, count)
}
