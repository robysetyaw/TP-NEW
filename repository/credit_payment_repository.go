package repository

import (
	"fmt"

	"gorm.io/gorm"

	"trackprosto/delivery/utils"
	model "trackprosto/models"
)

type CreditPaymentRepository interface {
	CreateCreditPayment(payment *model.CreditPayment) error
	GetAllCreditPayments() ([]*model.CreditPayment, error)
	GetCreditPaymentByID(id string) (*model.CreditPayment, error)
	UpdateCreditPayment(payment *model.CreditPayment) error
	GetTotalCredit(inv_number string) (float64, error)
	GetCreditPaymentsByInvoiceNumber(inv_number string) ([]*model.CreditPayment, error)
	CountCreditPayments(invoiceNumber string) (int, error)
}

type creditPaymentRepository struct {
	db *gorm.DB
}

func NewCreditPaymentRepository(db *gorm.DB) CreditPaymentRepository {
	return &creditPaymentRepository{
		db: db,
	}
}

func (repo *creditPaymentRepository) CreateCreditPayment(payment *model.CreditPayment) error {
	return repo.db.Create(payment).Error
}

func (repo *creditPaymentRepository) GetAllCreditPayments() ([]*model.CreditPayment, error) {
	var payments []*model.CreditPayment
	if err := repo.db.Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (repo *creditPaymentRepository) GetCreditPaymentByID(id string) (*model.CreditPayment, error) {
	var payment model.CreditPayment
	if err := repo.db.Where("id = ?", id).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (repo *creditPaymentRepository) UpdateCreditPayment(payment *model.CreditPayment) error {
	return repo.db.Save(payment).Error
}

func (repo *creditPaymentRepository) GetTotalCredit(inv_number string) (float64, error) {
	var total float64
	if err := repo.db.Model(&model.CreditPayment{}).Where("inv_number = ?", inv_number).Select("SUM(amount)").Row().Scan(&total); err != nil {
		return 0, fmt.Errorf("failed to get total credit: %w", err)
	}
	return total, nil
}

func (repo *creditPaymentRepository) GetCreditPaymentsByInvoiceNumber(inv_number string) ([]*model.CreditPayment, error) {
	var payments []*model.CreditPayment
	if err := repo.db.Where("inv_number = ?", inv_number).Order("created_at desc").Find(&payments).Error; err != nil {
		if condition := err.Error(); condition == "record not found" {
			return nil, utils.ErrCreditPaymentNotFound
		} else {
			return nil, fmt.Errorf("failed to get credit payments: %w", err)	
		}
	}
	return payments, nil
}

func (repo *creditPaymentRepository) CountCreditPayments(invoiceNumber string) (int, error) {
	var count int64
	err := repo.db.Model(&model.CreditPayment{}).Where("inv_number = ?", invoiceNumber).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("gagal menghitung pembayaran kredit: %w", err)
	}
	return int(count), nil
}