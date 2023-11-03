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
	GetAllMeats() ([]*model.MeatWithStock, error)
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

func (mc *meatUseCase) GetAllMeats() ([]*model.MeatWithStock, error) {
	meats, err := mc.meatRepository.GetAllMeats()
	if err != nil {
		log.WithField("error", err).Error("Failed to get all meats")
		return nil, err
	}
	todayDate := time.Now().Format("2006-01-02")
	var meatsWithStocks []*model.MeatWithStock
	for _, meat := range meats {
		stockIn, stockOut, err := mc.txRepository.CalculateMeatStockByDate(meat.ID, todayDate)
		if err != nil {
			log.WithField("error", err).Error("Failed to calculate meat stock")
			return nil, err
		}
		meatWithStock := &model.MeatWithStock{
			Meat:     meat,
			StockIn:  stockIn,
			StockOut: stockOut,
		}
		meatsWithStocks = append(meatsWithStocks, meatWithStock)
	}

	return meatsWithStocks, nil
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
	existingMeat, err := uc.meatRepository.GetMeatByName(meat.Name)
	if err != nil {
		log.WithField("error", err).Error("Failed to check meat name existence")
		return fmt.Errorf("failed to check meat name existence: %v", err)
	}
	if existingMeat != nil {
		log.WithField("meatName", meat.Name).Error("Meat name already exists")
		return fmt.Errorf("meat name already exists")
	}
	err = uc.meatRepository.UpdateMeat(meat)
	if err != nil {
		log.WithField("error", err).Error("Failed to update meat")
		return err
	}

	return nil
}
