package services

import (
	"errors"

	"zero-waste-kitchen/internal/models"
	"zero-waste-kitchen/internal/repositories"
	"zero-waste-kitchen/internal/utils"
)

type UserService interface {
	Register(user *models.User) error
	Login(email, password string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) error
	UpdateFCMToken(userID uint, token string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(user *models.User) error {
	user.Email = utils.NormalizeEmail(user.Email)

	if err := utils.ValidateStruct(user); err != nil {
		return err
	}

	existing, _ := s.repo.FindByEmail(user.Email)
	if existing != nil {
		return errors.New("email already registered")
	}

	// Hash the password before storing
	if err := user.HashPassword(); err != nil {
		return err
	}

	return s.repo.Create(user)
}

func (s *userService) Login(email, password string) (*models.User, error) {
	email = utils.NormalizeEmail(email)
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.ComparePassword(password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) UpdateUser(user *models.User) error {
	existing, err := s.repo.FindByID(user.ID)
	if err != nil {
		return err
	}

	existing.Name = user.Name

	// Only update password if it's provided
	if user.Password != "" {
		if err := user.HashPassword(); err != nil {
			return err
		}
		existing.Password = user.Password
	}

	return s.repo.Update(existing)
}

func (s *userService) UpdateFCMToken(userID uint, token string) error {
	return s.repo.UpdateFCMToken(userID, token)
}
