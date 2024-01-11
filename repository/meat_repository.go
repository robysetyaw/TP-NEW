package repository

import (
	"errors"
	"fmt"
	model "trackprosto/models"

	"gorm.io/gorm"
)

type MeatRepository interface {
	CreateMeat(meat *model.Meat) error
	GetMeatByID(string) (*model.Meat, error)
	GetAllMeats(page int, itemsPerPage int) ([]*model.Meat, int, error)
	GetMeatByName(string) (*model.Meat, error)
	UpdateMeat(meat *model.Meat) error
	DeleteMeat(string) error
	ReduceStock(meatID string, qty float64) error
	IncreaseStock(meatID string, qty float64) error
}

type meatRepository struct {
	db *gorm.DB
}

func NewMeatRepository(db *gorm.DB) MeatRepository {
	return &meatRepository{db: db}
}

func (mr *meatRepository) CreateMeat(meat *model.Meat) error {
	return mr.db.Create(&meat).Error
}

func (r *meatRepository) GetAllMeats(page int, itemsPerPage int) ([]*model.Meat, int, error) {
	var meats []*model.Meat

	if page < 1 {
		page = 1
	}

	var totalCount int64
	if err := r.db.Model(&model.Meat{}).Where("is_active = true").Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	totalPages := int((totalCount + int64(itemsPerPage) - 1) / int64(itemsPerPage))

	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * itemsPerPage

	if err := r.db.Where("is_active = ?", true).Order("name ASC").
		Offset(offset).Limit(itemsPerPage).
		Order("created_at desc").Find(&meats).Error; err != nil {
		return nil, totalPages, err
	}
	return meats, totalPages, nil
}

func (r *meatRepository) GetMeatByName(name string) (*model.Meat, error) {
	var meat model.Meat
	if err := r.db.Where("name = ? AND is_active = ?", name, true).First(&meat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("not Found")
		}
		return nil, err
	}
	return &meat, nil
}

func (r *meatRepository) GetMeatByID(id string) (*model.Meat, error) {
	var meat model.Meat
	if err := r.db.First(&meat, "id = ? AND is_active = ?", id, true).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &meat, nil
}

func (r *meatRepository) DeleteMeat(id string) error {
	return r.db.Model(&model.Meat{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *meatRepository) UpdateMeat(meat *model.Meat) error {
	return r.db.Save(&meat).Error
}

func (r *meatRepository) ReduceStock(meatID string, qty float64) error {
	return r.db.Model(&model.Meat{}).Where("id = ?", meatID).UpdateColumn("stock", gorm.Expr("stock - ?", qty)).Error
}

func (r *meatRepository) IncreaseStock(meatID string, qty float64) error {
	return r.db.Model(&model.Meat{}).Where("id = ?", meatID).UpdateColumn("stock", gorm.Expr("stock + ?", qty)).Error
}
