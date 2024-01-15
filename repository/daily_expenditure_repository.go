package repository

import (
	"fmt"
	"time"
	model "trackprosto/models"
	"gorm.io/gorm"
)

type DailyExpenditureRepository interface {
	CreateDailyExpenditure(expenditure *model.DailyExpenditure) error
	UpdateDailyExpenditure(expenditure *model.DailyExpenditure) error
	GetDailyExpenditureByID(id string) (*model.DailyExpenditure, error)
	GetAllDailyExpenditures() ([]*model.DailyExpenditure, error)
	DeleteDailyExpenditure(id string) error
	GetTotalExpenditureByDateRange(startDate time.Time, endDate time.Time) (float64, error)
	// GetExpendituresByDateRange(startDate time.Time, endDate time.Time) ([]*model.DailyExpenditureReport, error)
	GetLastNotaNumber(date string) (int, error)
}

type dailyExpenditureRepository struct {
	db *gorm.DB
}

func NewDailyExpenditureRepository(db *gorm.DB) DailyExpenditureRepository {
	return &dailyExpenditureRepository{
		db: db,
	}
}

func (repo *dailyExpenditureRepository) GetTotalExpenditureByDateRange(startDate time.Time, endDate time.Time) (float64, error) {
	var total float64
	result := repo.db.Model(&model.DailyExpenditure{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("created_at BETWEEN ? AND ? AND is_active = ?", startDate, endDate, true).
		Scan(&total)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to get total expenditure by date range: %w", result.Error)
	}
	return total, nil
}

func (repo *dailyExpenditureRepository) CreateDailyExpenditure(expenditure *model.DailyExpenditure) error {
	expenditure.CreatedAt = time.Now()
	expenditure.UpdatedAt = time.Now()
	expenditure.IsActive = true

	result := repo.db.Create(expenditure)
	if result.Error != nil {
		return fmt.Errorf("failed to create daily expenditure: %w", result.Error)
	}

	return nil
}

func (repo *dailyExpenditureRepository) UpdateDailyExpenditure(expenditure *model.DailyExpenditure) error {
	expenditure.UpdatedAt = time.Now()

	result := repo.db.Model(&model.DailyExpenditure{}).
		Where("id = ? AND is_active = ?", expenditure.ID, true).
		Updates(map[string]interface{}{
			"amount":      expenditure.Amount,
			"description": expenditure.Description,
			"is_active":   expenditure.IsActive,
			"updated_at":  expenditure.UpdatedAt,
			"updated_by":  expenditure.UpdatedBy,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to update daily expenditure: %w", result.Error)
	}

	return nil
}

func (repo *dailyExpenditureRepository) GetDailyExpenditureByID(id string) (*model.DailyExpenditure, error) {
	var expenditure model.DailyExpenditure

	result := repo.db.Model(&model.DailyExpenditure{}).
		Where("id = ? AND is_active = ?", id, true).
		First(&expenditure)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Daily expenditure not found
		}
		return nil, fmt.Errorf("failed to get daily expenditure by ID: %w", result.Error)
	}

	return &expenditure, nil
}

func (repo *dailyExpenditureRepository) GetAllDailyExpenditures() ([]*model.DailyExpenditure, error) {
	var expenditures []*model.DailyExpenditure

	result := repo.db.Model(&model.DailyExpenditure{}).
		Where("is_active = ?", true).
		Find(&expenditures)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all daily expenditures: %w", result.Error)
	}

	return expenditures, nil
}

func (repo *dailyExpenditureRepository) DeleteDailyExpenditure(id string) error {
	result := repo.db.Model(&model.DailyExpenditure{}).
		Where("id = ?", id).
		Update("is_active", false)

	if result.Error != nil {
		return fmt.Errorf("failed to delete daily expenditure: %w", result.Error)
	}

	return nil
}

// func (repo *dailyExpenditureRepository) GetExpendituresByDateRange(startDate time.Time, endDate time.Time) ([]*model.DailyExpenditureReport, error) {
// 	var expenditures []*model.DailyExpenditureReport

// 	result := repo.db.Model(&model.DailyExpenditure{}).
// 		Select("daily_expenditures.id, daily_expenditures.user_id, users.username, daily_expenditures.amount, daily_expenditures.description, daily_expenditures.created_at, daily_expenditures.updated_at, daily_expenditures.date").
// 		Joins("JOIN users ON daily_expenditures.user_id = users.id").
// 		Where("DATE(daily_expenditures.created_at) >= ? AND DATE(daily_expenditures.created_at) <= ? AND daily_expenditures.is_active = ?", startDate, endDate, true).
// 		Scan(&expenditures)

// 	if result.Error != nil {
// 		return nil, fmt.Errorf("failed to get expenditures by date range: %w", result.Error)
// 	}

// 	return expenditures, nil
// }

func (repo *dailyExpenditureRepository) GetLastNotaNumber(date string) (int, error) {
	var count int64
	result := repo.db.Model(&model.DailyExpenditure{}).
		Where("DATE(created_at) = ?", date).
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to get last nota number: %w", result.Error)
	}

	return int(count) + 1, nil
}
