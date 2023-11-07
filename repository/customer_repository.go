package repository

import (
	model "trackprosto/models"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	CreateCustomer(*model.CustomerModel) (*model.CustomerModel, error)
	UpdateCustomer(*model.CustomerModel) error
	GetCustomerById(string) (*model.CustomerModel, error)
	GetCustomerByName(string) (*model.CustomerModel, error)
	GetAllCustomer(page int, itemsPerPage int) ([]*model.CustomerModel, int, error)
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

func (repo *customerRepository) CreateCustomer(customer *model.CustomerModel) (*model.CustomerModel, error) {
	err := repo.db.Create(customer).Error
	if err != nil {
		return nil, err
	}
	return customer, nil
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

func (repo *customerRepository) GetAllCustomer(page int, itemsPerPage int) ([]*model.CustomerModel, int, error) {
	var customers []*model.CustomerModel
	if page < 1 {
		page = 1
	}

	var totalCount int64
	if err := repo.db.Model(&model.Meat{}).Where("is_active = true").Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	totalPages := int((totalCount + int64(itemsPerPage) - 1) / int64(itemsPerPage))

	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * itemsPerPage
	if err := repo.db.Where("is_active = ?", true).Offset(offset).Limit(itemsPerPage).
		Order("created_at desc").Find(&customers).Error; err != nil {
		return nil, 0, err
	}
	return customers, totalPages, nil
}

func (repo *customerRepository) DeleteCustomer(id string) error {
	return repo.db.Delete(&model.CustomerModel{}, "id = ?", id).Error
}
