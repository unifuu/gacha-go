package services

import (
	"gacha/models"
	"sync"
)

// UserService handles user management
type UserService struct {
	users map[string]*models.User
	mu    sync.RWMutex
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	service := &UserService{
		users: make(map[string]*models.User),
	}

	// Initialize default user
	service.users["default"] = &models.User{
		ID:        1,
		Username:  "default",
		Currency:  10000,
		Inventory: []models.Character{},
		PityCount: 0,
	}

	return service
}

// GetUser retrieves a user by username
func (s *UserService) GetUser(username string) *models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users[username]
}

// GetDefaultUser retrieves the default user
func (s *UserService) GetDefaultUser() *models.User {
	return s.GetUser("default")
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(username string, user *models.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[username] = user
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username string) *models.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[username]; exists {
		return s.users[username]
	}

	user := &models.User{
		ID:        len(s.users) + 1,
		Username:  username,
		Currency:  1000,
		Inventory: []models.Character{},
		PityCount: 0,
	}

	s.users[username] = user
	return user
}
