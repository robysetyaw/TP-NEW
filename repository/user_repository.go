package repository

import (
	"errors"
	"time"
	model "trackprosto/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
	DeleteUser(id string) error
	GetByUsername(username string) (*model.User, error)
	CountUsers(username string) (int, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CountUsers(username string) (int, error) {
    var count int64
    if err := r.db.Model(&model.User{}).Where("is_active = true AND username = ?", username).Count(&count).Error; err != nil {
        return 0, err
    }
    return int(count), nil
}

func (r *userRepository) CreateUser(user *model.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	return r.db.Create(&user).Error
}

func (r *userRepository) UpdateUser(user *model.User) error {
	user.UpdatedAt = time.Now()
	return r.db.Save(&user).Error
}

func (r *userRepository) GetUserByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "username = ? AND is_active = true", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAllUsers() ([]*model.User, error) {
	var users []*model.User
	err := r.db.Where("is_active = true").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) DeleteUser(username string) error {
	return r.db.Model(&model.User{}).Where("username = ?", username).Update("is_active", false).Error
}
