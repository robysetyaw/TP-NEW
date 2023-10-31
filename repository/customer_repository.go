package repository

import (
	model "trackprosto/models"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	CreateCustomer(*model.CustomerModel) error
	UpdateCustomer(*model.CustomerModel) error
	GetCustomerById(string) (*model.CustomerModel, error)
	GetCustomerByName(string) (*model.CustomerModel, error)
	GetAllCustomer() ([]*model.CustomerModel, error)
	DeleteCustomer(string) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

func (repo *customerRepository) CreateCustomer(customer *model.CustomerModel) error {
	return repo.db.Create(customer).Error
}

func (repo *customerRepository) UpdateCustomer(customer *model.CustomerModel) error {
	return repo.db.Save(customer).Error
}

func (repo *customerRepository) GetCustomerById(id string) (*model.CustomerModel, error) {
	var customer model.CustomerModel
	if err := repo.db.First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (repo *customerRepository) GetCustomerByName(name string) (*model.CustomerModel, error) {
	var customer model.CustomerModel
	if err := repo.db.First(&customer, "fullname = ?", name).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (repo *customerRepository) GetAllCustomer() ([]*model.CustomerModel, error) {
	var customers []*model.CustomerModel
	if err := repo.db.Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func (repo *customerRepository) DeleteCustomer(id string) error {
	return repo.db.Delete(&model.CustomerModel{}, "id = ?", id).Error
}
