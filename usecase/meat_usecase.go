package usecase

import (
	"fmt"
	model "trackprosto/models"
	"trackprosto/repository"
)

type MeatUseCase interface {
	CreateMeat(meat *model.Meat) error
	GetMeatById(string) (*model.Meat, error)
	GetAllMeats() ([]*model.Meat, error)
	GetMeatByName(string) (*model.Meat, error)
	UpdateMeat(meat *model.Meat) error
	DeleteMeat(string) error
}

type meatUseCase struct {
	meatRepository repository.MeatRepository
}

func NewMeatUseCase(meatRepo repository.MeatRepository) MeatUseCase {
	return &meatUseCase{meatRepository: meatRepo}
}

func (ms *meatUseCase) CreateMeat(meat *model.Meat) error {

	isExist, _ := ms.meatRepository.GetMeatByName(meat.Name)
	if isExist != nil {
		return fmt.Errorf("meatname already exists")
	}
	meat.IsActive = true
	err := ms.meatRepository.CreateMeat(meat)
	if err != nil {
		return err
	}

	return nil
}

func (mc *meatUseCase) GetAllMeats() ([]*model.Meat, error) {
	meats, err := mc.meatRepository.GetAllMeats()
	if err != nil {
		// Handle any repository errors or perform error logging
		return nil, err
	}

	// Perform any additional data processing or transformation if needed

	return meats, nil
}

func (mc *meatUseCase) GetMeatByName(name string) (*model.Meat, error) {
	meat, err := mc.meatRepository.GetMeatByName(name)
	if err != nil {
		// Handle any repository errors or perform error logging
		return nil, err
	}

	// Perform any additional data processing or transformation if needed

	return meat, nil
}

func (mc *meatUseCase) GetMeatById(id string) (*model.Meat, error) {
	meat, err := mc.meatRepository.GetMeatByName(id)
	if err != nil {
		// Handle any repository errors or perform error logging
		return nil, err
	}

	// Perform any additional data processing or transformation if needed

	return meat, nil
}

func (mc *meatUseCase) DeleteMeat(id string) error {
	// Implement any business logic or validation before deleting the meat
	existingMeat, err := mc.meatRepository.GetMeatByName(id)
	if err != nil {
		return fmt.Errorf("failed to check meatname existence: %v", err)
	}
	if existingMeat != nil {
		return fmt.Errorf("meatname already exists")
	}
	err = mc.meatRepository.DeleteMeat(id)
	if err != nil {
		// Handle any repository errors or perform error logging
		return nil
	}

	// Perform any additional data processing or transformation if needed

	return nil
}

func (uc *meatUseCase) UpdateMeat(meat *model.Meat) error {
	// Implement any business logic or validation before updating the meat
	// You can also perform data manipulation or enrichment if needed
	existingMeat, err := uc.meatRepository.GetMeatByName(meat.Name)
	if err != nil {
		return fmt.Errorf("failed to check meatname existence: %v", err)
	}
	if existingMeat != nil {
		return fmt.Errorf("meatname already exists")
	}
	err = uc.meatRepository.UpdateMeat(meat)
	if err != nil {
		// Handle any repository errors or perform error logging
		return err
	}

	return nil
}
