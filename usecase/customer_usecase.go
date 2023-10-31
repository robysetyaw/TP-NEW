package usecase

import (
	model "trackprosto/models"
	"trackprosto/repository"
)

type CustomerUsecase interface {
	CreateCustomer(customer *model.CustomerModel) error
	UpdateCustomer(customer *model.CustomerModel) error
	GetCustomerById(id string) (*model.CustomerModel, error)
	GetCustomerByName(name string) (*model.CustomerModel, error)
	GetAllCustomers() ([]*model.CustomerModel, error)
	DeleteCustomer(id string) error
}

type customerUsecase struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerUsecase(cr repository.CustomerRepository) CustomerUsecase {
	return &customerUsecase{
		customerRepo: cr,
	}
}

func (uc *customerUsecase) CreateCustomer(customer *model.CustomerModel) error {
	return uc.customerRepo.CreateCustomer(customer)
}

func (uc *customerUsecase) UpdateCustomer(customer *model.CustomerModel) error {
	return uc.customerRepo.UpdateCustomer(customer)
}

func (uc *customerUsecase) GetCustomerById(id string) (*model.CustomerModel, error) {
	return uc.customerRepo.GetCustomerById(id)
}

func (uc *customerUsecase) GetCustomerByName(name string) (*model.CustomerModel, error) {
	return uc.customerRepo.GetCustomerByName(name)
}

func (uc *customerUsecase) GetAllCustomers() ([]*model.CustomerModel, error) {
	return uc.customerRepo.GetAllCustomer()
}

func (uc *customerUsecase) DeleteCustomer(id string) error {
	return uc.customerRepo.DeleteCustomer(id)
}
