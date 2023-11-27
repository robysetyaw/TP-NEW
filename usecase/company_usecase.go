package usecase

import (
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/models/dto"
	"trackprosto/repository"

	"github.com/sirupsen/logrus"
)

type CompanyUseCase interface {
	CreateCompany(*model.Company) error
	UpdateCompany(*dto.CompanyRequest) (*dto.CompanyResponse, error)
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
	isExist, err := cu.companyRepo.GetCompanyByName(company.CompanyName)
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

func (cu *companyUseCase) UpdateCompany(companyRequest *dto.CompanyRequest) (*dto.CompanyResponse, error) {
	currentCompany, err := cu.companyRepo.GetCompanyById(companyRequest.ID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if currentCompany == nil {
		return nil, utils.ErrCompanyNotFound
	}
	if companyRequest.CompanyName != "" {
		isExist, err := cu.companyRepo.GetCompanyByName(companyRequest.CompanyName)
		if err != nil {
			return nil, err
		}
		if isExist != nil {
			return nil, utils.ErrCompanyNameAlreadyExist
		}
	}

	currentCompany.Address = utils.NonEmpty(companyRequest.Address, currentCompany.Address)
	currentCompany.CompanyName = utils.NonEmpty(companyRequest.CompanyName, currentCompany.CompanyName)
	currentCompany.Email = utils.NonEmpty(companyRequest.Email, currentCompany.Email)
	currentCompany.PhoneNumber = utils.NonEmpty(companyRequest.PhoneNumber, currentCompany.PhoneNumber)
	err = cu.companyRepo.UpdateCompany(currentCompany)
	if err != nil {
		return nil, err
	}

	companyResponse := &dto.CompanyResponse{
		ID:          currentCompany.ID,
		CompanyName: currentCompany.CompanyName,
		Address:     currentCompany.Address,
		Email:       currentCompany.Email,
		PhoneNumber: currentCompany.PhoneNumber,
	}

	return companyResponse, nil
}

func (cu *companyUseCase) GetCompanyById(id string) (*model.Company, error) {
	company, err := cu.companyRepo.GetCompanyById(id)
	if company == nil {
		return nil, utils.ErrCompanyNotFound
	}
	if err != nil {
		return nil, err
	}
	return company, nil
}

func (cu *companyUseCase) GetAllCompany() ([]*model.Company, error) {
	return cu.companyRepo.GetAllCompany()
}

func (cu *companyUseCase) DeleteCompany(id string) error {
	currentCompany, err := cu.companyRepo.GetCompanyById(id)
	if currentCompany == nil {
		return utils.ErrCompanyNotFound
	}
	if err != nil {
		return err
	}
	currentCompany.IsActive = false
	err = cu.companyRepo.UpdateCompany(currentCompany)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
