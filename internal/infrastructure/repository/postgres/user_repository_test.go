package postgres_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/infrastructure/repository/postgres"
	"github.com/zenfulcode/commercify/testutil/db"
)

func TestUserRepository_Integration(t *testing.T) {
	// Skip if short tests flag is set
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	testDB := db.SetupTestDB(t)
	defer testDB.Close()

	// Create repository
	userRepo := postgres.NewUserRepository(testDB.DB)

	// Test user data
	email := "test@example.com"
	password := "password123"
	firstName := "Test"
	lastName := "User"
	role := entity.RoleUser

	// Test Create
	t.Run("Create", func(t *testing.T) {
		user, err := entity.NewUser(email, password, firstName, lastName, role)
		assert.NoError(t, err)

		err = userRepo.Create(user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
	})

	// Test GetByID
	t.Run("GetByID", func(t *testing.T) {
		// First create a user
		user, _ := entity.NewUser("getbyid@example.com", password, firstName, lastName, role)
		_ = userRepo.Create(user)

		// Get the user by ID
		fetchedUser, err := userRepo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, fetchedUser.ID)
		assert.Equal(t, user.Email, fetchedUser.Email)
		assert.Equal(t, user.FirstName, fetchedUser.FirstName)
		assert.Equal(t, user.LastName, fetchedUser.LastName)
		assert.Equal(t, user.Role, fetchedUser.Role)

		// Test getting non-existent user
		_, err = userRepo.GetByID(9999)
		assert.Error(t, err)
	})

	// Test GetByEmail
	t.Run("GetByEmail", func(t *testing.T) {
		// First create a user
		uniqueEmail := "unique@example.com"
		user, _ := entity.NewUser(uniqueEmail, password, firstName, lastName, role)
		_ = userRepo.Create(user)

		// Get the user by email
		fetchedUser, err := userRepo.GetByEmail(uniqueEmail)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, fetchedUser.ID)
		assert.Equal(t, uniqueEmail, fetchedUser.Email)

		// Test getting non-existent email
		_, err = userRepo.GetByEmail("nonexistent@example.com")
		assert.Error(t, err)
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		// First create a user
		user, _ := entity.NewUser("update@example.com", password, firstName, lastName, role)
		_ = userRepo.Create(user)

		// Update user
		user.FirstName = "Updated"
		user.LastName = "Name"
		err := userRepo.Update(user)
		assert.NoError(t, err)

		// Verify update
		updatedUser, _ := userRepo.GetByID(user.ID)
		assert.Equal(t, "Updated", updatedUser.FirstName)
		assert.Equal(t, "Name", updatedUser.LastName)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		// First create a user
		user, _ := entity.NewUser("delete@example.com", password, firstName, lastName, role)
		_ = userRepo.Create(user)

		// Delete user
		err := userRepo.Delete(user.ID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = userRepo.GetByID(user.ID)
		assert.Error(t, err)
	})

	// Test List
	t.Run("List", func(t *testing.T) {
		// Clean database first to ensure consistent results
		testDB.Clean()

		// Create multiple users
		for i := 0; i < 5; i++ {
			email := fmt.Sprintf("list%d@example.com", i)
			user, _ := entity.NewUser(email, password, firstName, lastName, role)
			_ = userRepo.Create(user)
		}

		// Test pagination
		users, err := userRepo.List(0, 3)
		assert.NoError(t, err)
		assert.Len(t, users, 3)

		users, err = userRepo.List(3, 3)
		assert.NoError(t, err)
		assert.Len(t, users, 2)
	})
}
