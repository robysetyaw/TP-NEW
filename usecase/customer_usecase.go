package usecase

import (
	"trackprosto/delivery/utils"
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
	GetAllCustomerByCompanyId(page int, itemsPerPage int, company_id string) ([]*model.CustomerModel, int, error)
}

type customerUsecase struct {
	customerRepo repository.CustomerRepository
	companyRepo  repository.CompanyRepository
}

func NewCustomerUsecase(cr repository.CustomerRepository, cpr repository.CompanyRepository) CustomerUsecase {
	return &customerUsecase{
		customerRepo: cr,
		companyRepo:  cpr,
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
	currentCustomer, err := uc.GetCustomerById(customer.Id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if currentCustomer == nil {
		logrus.Error(utils.ErrCustomerNotFound)
		return utils.ErrCustomerNotFound
	}
	customer.FullName = utils.NonEmpty(customer.FullName, currentCustomer.FullName)
	customer.Address = utils.NonEmpty(customer.Address, currentCustomer.Address)
	customer.PhoneNumber = utils.NonEmpty(customer.PhoneNumber, currentCustomer.PhoneNumber)
	customer.CreatedAt = currentCustomer.CreatedAt
	customer.CreatedBy = currentCustomer.CreatedBy
	return uc.customerRepo.UpdateCustomer(customer)
}

func (uc *customerUsecase) GetCustomerById(id string) (*model.CustomerModel, error) {
	customer, err := uc.customerRepo.GetCustomerById(id)
	if customer == nil {
		return nil, utils.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (uc *customerUsecase) GetCustomerByName(name string) (*model.CustomerModel, error) {
	return uc.customerRepo.GetCustomerByName(name)
}

func (uc *customerUsecase) GetAllCustomers(page int, itemsPerPage int) ([]*model.CustomerModel, int, error) {
	customers, totalPages, err := uc.customerRepo.GetAllCustomer(page, itemsPerPage)
	if err != nil {
		logrus.Error(err)
		return nil, 0, err
	}
	return customers, totalPages, nil
}

func (uc *customerUsecase) DeleteCustomer(id string) error {
	customer, err := uc.GetCustomerById(id)
	if customer == nil {
		logrus.Error(utils.ErrCustomerNotFound)
		return utils.ErrCustomerNotFound
	}
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = uc.customerRepo.UpdateCustomer(customer)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (uc *customerUsecase) GetAllCustomerByCompanyId(page int, itemsPerPage int, company_id string) ([]*model.CustomerModel, int, error) {
	companies, err := uc.companyRepo.GetCompanyById(company_id)
	if companies == nil {
		return nil, 0, utils.ErrCompanyNotFound
	}
	if err != nil {
		return nil, 0, err
	}
	customers, totalPages, err := uc.customerRepo.GetAllCustomerByCompanyId(page, itemsPerPage, company_id)
	if customers == nil {
		return nil, 0, utils.ErrCustomerNotFound
	}
	if err != nil {
		return nil, 0, err
	}
	return customers, totalPages, nil
}