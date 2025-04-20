package mock

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// MockUserRepository is a mock implementation of user repository for testing
type MockUserRepository struct {
	users       map[uint]*entity.User
	userByEmail map[string]*entity.User
	lastID      uint
}

// NewMockUserRepository creates a new instance of MockUserRepository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:       make(map[uint]*entity.User),
		userByEmail: make(map[string]*entity.User),
		lastID:      0,
	}
}

// Create adds a user to the repository
func (r *MockUserRepository) Create(user *entity.User) error {
	// Check if email already exists
	if _, exists := r.userByEmail[user.Email]; exists {
		return errors.New("user with this email already exists")
	}

	// Assign ID
	r.lastID++
	user.ID = r.lastID

	// Store user
	r.users[user.ID] = user
	r.userByEmail[user.Email] = user

	return nil
}

// GetByID retrieves a user by ID
func (r *MockUserRepository) GetByID(id uint) (*entity.User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetByEmail retrieves a user by email
func (r *MockUserRepository) GetByEmail(email string) (*entity.User, error) {
	user, exists := r.userByEmail[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// Update updates a user
func (r *MockUserRepository) Update(user *entity.User) error {
	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	// Update user
	r.users[user.ID] = user
	r.userByEmail[user.Email] = user

	return nil
}

// Delete removes a user
func (r *MockUserRepository) Delete(id uint) error {
	user, exists := r.users[id]
	if !exists {
		return errors.New("user not found")
	}

	// Remove user from both maps
	delete(r.userByEmail, user.Email)
	delete(r.users, id)

	return nil
}

// List retrieves users with pagination
func (r *MockUserRepository) List(offset, limit int) ([]*entity.User, error) {
	result := make([]*entity.User, 0)
	count := 0
	skip := offset

	// Iterate through users and apply pagination
	for _, user := range r.users {
		if skip > 0 {
			skip--
			continue
		}

		result = append(result, user)
		count++

		if count >= limit {
			break
		}
	}

	return result, nil
}
