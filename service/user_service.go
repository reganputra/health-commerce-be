package service

import (
	"errors"
	"fmt"
	"health-store/models"
	"health-store/repositories"

	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic for users
type UserService struct {
	userRepo repositories.UserRepositoryInterface
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepositoryInterface) *UserService {
	return &UserService{userRepo: userRepo}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(req models.UserRegisterRequest) (*models.User, error) {
	// Check if username already exists
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Username:      req.Username,
		Password:      string(hashedPassword),
		Email:         req.Email,
		Dob:           req.Dob,
		Gender:        req.Gender,
		Address:       req.Address,
		City:          req.City,
		ContactNumber: req.ContactNumber,
		Role:          "customer", // Default role
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser authenticates a user login
func (s *UserService) AuthenticateUser(req models.UserLoginRequest) (*models.User, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("username not found: %s", req.Username)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("password mismatch for user: %s", req.Username)
	}

	return user, nil
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// GetAllUsers gets all users
func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.FindAll()
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}
