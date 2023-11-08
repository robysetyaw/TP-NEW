package usecase

import (
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/repository"

	"github.com/sirupsen/logrus"
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
	isExist,err := cu.companyRepo.GetCompanyByName(company.CompanyName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if isExist != nil {
		logrus.WithField("error", err).Error("Company name already exists")
		return utils.ErrCompanyNameAlreadyExist
	}
	err = cu.companyRepo.CreateCompany(company)
	return err 
}

func (cu *companyUseCase) UpdateCompany(company *model.Company) error {
	currentCompany, err := cu.companyRepo.GetCompanyById(company.ID)
	
	if err != nil {
		logrus.Error(err)
		return err
	}
	if currentCompany == nil {
		logrus.WithField("error", err).Error("Failed to get company by ID")
		return utils.ErrCompanyNotFound
	}

	company.Address = utils.NonEmpty(company.Address, currentCompany.Address)
	company.CompanyName = utils.NonEmpty(company.CompanyName, currentCompany.CompanyName)
	company.Email = utils.NonEmpty(company.Email, currentCompany.Email)
	company.PhoneNumber = utils.NonEmpty(company.PhoneNumber, currentCompany.PhoneNumber)
	company.IsActive = currentCompany.IsActive
	company.CreatedAt = currentCompany.CreatedAt
	company.CreatedBy = currentCompany.CreatedBy

	err = cu.companyRepo.UpdateCompany(company)
	return err
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
