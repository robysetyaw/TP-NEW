package repository

import (
	"errors"
	model "trackprosto/models"

	"gorm.io/gorm"
)

type CompanyRepository interface {
	CreateCompany(*model.Company) error
	UpdateCompany(*model.Company) error
	GetCompanyById(string) (*model.Company, error)
	GetCompanyByName(string) (*model.Company, error)
	GetAllCompany() ([]*model.Company, error)
	DeleteCompany(string) error
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{
		db: db,
	}
}

func (repo *companyRepository) CreateCompany(company *model.Company) error {
	return repo.db.Create(company).Error
}

func (repo *companyRepository) UpdateCompany(company *model.Company) error {
	return repo.db.Save(company).Error
}

func (repo *companyRepository) GetCompanyById(id string) (*model.Company, error) {
	var company model.Company
	if err := repo.db.First(&company, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // company not found
		}
		return nil, err
	}
	return &company, nil
}

func (repo *companyRepository) GetCompanyByName(companyName string) (*model.Company, error) {
	var company model.Company
	if err := repo.db.First(&company, "company_name = ?", companyName).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // company not found
		}
		return nil, err
	}
	return &company, nil
}

func (repo *companyRepository) GetAllCompany() ([]*model.Company, error) {
	var companies []*model.Company
	if err := repo.db.Where("is_active = ?", true).Order("created_at desc").Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

func (repo *companyRepository) DeleteCompany(id string) error {
	return repo.db.Delete(&model.Company{}, "id = ?", id).Error
}
