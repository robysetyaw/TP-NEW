package usecase

import (
	model "trackprosto/models"
	"trackprosto/repository"

	"github.com/sirupsen/logrus"
)

type CustomerUsecase interface {
	CreateCustomer(customer *model.CustomerModel) (*model.CustomerModel, error)
	UpdateCustomer(customer *model.CustomerModel) error
	GetCustomerById(id string) (*model.CustomerModel, error)
	GetCustomerByName(name string) (*model.CustomerModel, error)
	GetAllCustomers(page int, itemsPerPage int) ([]*model.CustomerModel, int, error)
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

func (uc *customerUsecase) CreateCustomer(customer *model.CustomerModel) (*model.CustomerModel, error) {
	customer, err := uc.customerRepo.CreateCustomer(customer)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return customer, nil
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

func (uc *customerUsecase) GetAllCustomers(page int, itemsPerPage int) ([]*model.CustomerModel, int, error) {
	customers, totalPages, err  := uc.customerRepo.GetAllCustomer(page, itemsPerPage)
	if err != nil {
		logrus.Error(err)
		return nil, 0, err
	}
	return customers, totalPages, nil
}

func (uc *customerUsecase) DeleteCustomer(id string) error {
	return uc.customerRepo.DeleteCustomer(id)
}
