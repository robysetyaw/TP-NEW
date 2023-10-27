package usecase

import (
	"fmt"
	"time"
	model "trackprosto/models"
	"trackprosto/repository"
)

type UserUseCase interface {
	CreateUser(user *model.User) error
	UpdateUser(user *model.User, username string) error
	GetUserByID(id string) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
	DeleteUser(id string) error
	GetUserByUsername(username string) (*model.User, error)
}

type userUseCase struct {
	userRepository repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepository: userRepo,
	}
}

func (uc *userUseCase) CreateUser(user *model.User) error {
	// Implement any business logic or validation before creating the user
	// You can also perform data manipulation or enrichment if needed

	existingUser, err := uc.userRepository.GetByUsername(user.Username)
	if err != nil {
		return fmt.Errorf("failed to check username existence: %v", err)
	}
	if existingUser != nil {
		return fmt.Errorf("username already exists")
	}

	user.IsActive = true
	user.CreatedAt = time.Now()
	user.CreatedBy = "admin"

	err = uc.userRepository.CreateUser(user)
	if err != nil {
		// Handle any repository errors or perform error logging
		return err
	}

	return nil
}

func (uc *userUseCase) UpdateUser(user *model.User, username string) error {
	// Implement any business logic or validation before updating the user
	// You can also perform data manipulation or enrichment if needed
	if user.Username != username {
		existingUser, err := uc.userRepository.GetByUsername(user.Username)
		if err != nil {
			return fmt.Errorf("failed to check username existence: %v", err)
		}
		if existingUser != nil {
			return fmt.Errorf("username already exists")
		}
	}
	err := uc.userRepository.UpdateUser(user)
	if err != nil {
		// Handle any repository errors or perform error logging
		return err
	}

	return nil
}

func (uc *userUseCase) GetUserByID(id string) (*model.User, error) {
	user, err := uc.userRepository.GetUserByID(id)
	if err != nil {
		// Handle any repository errors or perform error logging
		return nil, err
	}

	// Perform any additional data processing or transformation if needed

	return user, nil
}

func (uc *userUseCase) GetUserByUsername(username string) (*model.User, error) {
	user, err := uc.userRepository.GetByUsername(username)
	if err != nil {
		// Handle any repository errors or perform error logging
		return nil, err
	}

	// Perform any additional data processing or transformation if needed

	return user, nil
}

func (uc *userUseCase) GetAllUsers() ([]*model.User, error) {
	users, err := uc.userRepository.GetAllUsers()
	if err != nil {
		// Handle any repository errors or perform error logging
		return nil, err
	}

	// Perform any additional data processing or transformation if needed

	return users, nil
}

func (uc *userUseCase) DeleteUser(username string) error {
	// Implement any business logic or validation before deleting the user
	existingUser, err := uc.userRepository.GetByUsername(username)
	if err != nil {
		return fmt.Errorf("failed to check username existence: %v", err)
	}
	if existingUser == nil {
		return fmt.Errorf("user Not Found")
	}
	err = uc.userRepository.DeleteUser(username)
	if err != nil {
		// Handle any repository errors or perform error logging
		return err
	}

	return nil
}
