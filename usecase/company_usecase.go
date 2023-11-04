package usecase

import (
	model "trackprosto/models"
	"trackprosto/repository"
)

type CompanyUseCase interface {
	CreateCompany(*model.Company) error
	UpdateCompany(*model.Company) error
	GetCompanyById(string) (*model.Company, error)
	GetAllCompany() ([]*model.Company, error)
	DeleteCompany(string) error
}

type companyUseCase struct {
	companyRepo repository.CompanyRepository
}

func NewCompanyUseCase(companyRepo repository.CompanyRepository) CompanyUseCase {
	return &companyUseCase{
		companyRepo: companyRepo,
	}
}

func (cu *companyUseCase) CreateCompany(company *model.Company) error {
	// Perform any business logic or validation before creating the company
	// ...

	return cu.companyRepo.CreateCompany(company)
}

func (cu *companyUseCase) UpdateCompany(company *model.Company) error {
	// Perform any business logic or validation before updating the company
	// ...

	return cu.companyRepo.UpdateCompany(company)
}

func (cu *companyUseCase) GetCompanyById(id string) (*model.Company, error) {
	return cu.companyRepo.GetCompanyById(id)
}

func (cu *companyUseCase) GetAllCompany() ([]*model.Company, error) {
	return cu.companyRepo.GetAllCompany()
}

func (cu *companyUseCase) DeleteCompany(id string) error {
	return cu.companyRepo.DeleteCompany(id)
}
