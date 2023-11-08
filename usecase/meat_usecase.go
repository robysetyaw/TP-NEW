package usecase

import (
	"fmt"
	"time"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/repository"

	log "github.com/sirupsen/logrus"
)

type MeatUseCase interface {
	CreateMeat(meat *model.Meat) error
	GetMeatById(string) (*model.Meat, error)
	GetAllMeats(page int, itemsPerPage int) ([]*model.MeatWithStock, int, error)
	GetMeatByName(string) (*model.Meat, error)
	UpdateMeat(meat *model.Meat) error
	DeleteMeat(string) error
}

type meatUseCase struct {
	meatRepository repository.MeatRepository
	txRepository   repository.TransactionRepository
}

func NewMeatUseCase(meatRepo repository.MeatRepository, txRepository repository.TransactionRepository) MeatUseCase {
	return &meatUseCase{
		meatRepository: meatRepo,
		txRepository:   txRepository,
	}
}

func (ms *meatUseCase) CreateMeat(meat *model.Meat) error {
	isExist, _ := ms.meatRepository.GetMeatByName(meat.Name)
	if isExist != nil {
		log.WithField("meatName", meat.Name).Error("Meat name already exists")
		return utils.ErrMeatNameAlreadyExist
	}
	meat.IsActive = true
	err := ms.meatRepository.CreateMeat(meat)
	if err != nil {
		log.WithField("error", err).Error("Failed to create meat")
		return err
	}

	return nil
}

func (mc *meatUseCase) GetAllMeats(page int, itemsPerPage int) ([]*model.MeatWithStock, int, error) {
	meats, totalPages, err := mc.meatRepository.GetAllMeats(page, itemsPerPage)
	if err != nil {
		log.WithField("error", err).Error("Failed to get all meats")
		return nil, 0, err
	}
	todayDate := time.Now().Format("2006-01-02")
	var meatsWithStocks []*model.MeatWithStock
	for _, meat := range meats {
		stockIn, stockOut, err := mc.txRepository.CalculateMeatStockByDate(meat.ID, todayDate)
		if err != nil {
			log.WithField("error", err).Error("Failed to calculate meat stock")
			return nil, 0, err
		}
		meatWithStock := &model.MeatWithStock{
			Meat:     meat,
			StockIn:  stockIn,
			StockOut: stockOut,
		}
		meatsWithStocks = append(meatsWithStocks, meatWithStock)
	}

	return meatsWithStocks, totalPages, nil
}

func (mc *meatUseCase) GetMeatByName(name string) (*model.Meat, error) {
	meat, err := mc.meatRepository.GetMeatByName(name)
	if err != nil {
		log.WithField("error", err).Error("Failed to get meat by name")
		return nil, err
	}

	return meat, nil
}

func (mc *meatUseCase) GetMeatById(id string) (*model.Meat, error) {
	meat, err := mc.meatRepository.GetMeatByName(id)
	if err != nil {
		log.WithField("error", err).Error("Failed to get meat by ID")
		return nil, err
	}

	return meat, nil
}

func (mc *meatUseCase) DeleteMeat(id string) error {
	// Implement any business logic or validation before deleting the meat
	existingMeat, err := mc.meatRepository.GetMeatByName(id)
	if err != nil {
		log.WithField("error", err).Error("Failed to check meat name existence")
		return fmt.Errorf("failed to check meat name existence: %v", err)
	}
	if existingMeat != nil {
		log.WithField("meatName", id).Error("Meat name already exists")
		return fmt.Errorf("meat name already exists")
	}
	err = mc.meatRepository.DeleteMeat(id)
	if err != nil {
		log.WithField("error", err).Error("Failed to delete meat")
		return err
	}

	return nil
}

func (uc *meatUseCase) UpdateMeat(meat *model.Meat) error {
	// Implement any business logic or validation before updating the meat
	// You can also perform data manipulation or enrichment if needed
	currentMeatValue, err := uc.meatRepository.GetMeatByID(meat.ID)
	if err != nil {
		log.WithField("error", err).Error("Failed to get meat by ID")
		return fmt.Errorf("failed to get meat by ID: %v", err)
	}
	if currentMeatValue == nil {
		log.WithField("meatID", meat.ID).Error("Meat not found")
		return utils.ErrMeatNotFound
	}
	existingMeat, _ := uc.meatRepository.GetMeatByName(meat.Name)
	if existingMeat != nil && existingMeat.ID != meat.ID {
		log.WithField("meatName", meat.Name).Error("Meat name already exists")
		return utils.ErrMeatNameAlreadyExist
	}
	meat.CreatedBy = currentMeatValue.CreatedBy
	meat.CreatedAt = currentMeatValue.CreatedAt
	meat.Name = utils.NonEmpty(meat.Name, currentMeatValue.Name)
	meat.Stock = utils.NonZero(meat.Stock, currentMeatValue.Stock)
	meat.Price = utils.NonZero(meat.Price, currentMeatValue.Price)
	meat.IsActive = currentMeatValue.IsActive
	meat.UpdatedAt = time.Now()
	err = uc.meatRepository.UpdateMeat(meat)
	if err != nil {
		log.WithField("error", err).Error("Failed to update meat")
		return err
	}
	return nil
}
