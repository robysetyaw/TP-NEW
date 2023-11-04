package usecase

import (
	"fmt"
	"time"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TransactionUseCase interface {
	CreateTransaction(transaction *model.TransactionHeader) (*model.TransactionHeader, error)
	GetAllTransactions(page int, itemsPerPage int) ([]*model.TransactionHeader, int, error)
	GetTransactionByID(id string) (*model.TransactionHeader, error)
	DeleteTransaction(id string) error
	GetTransactionByInvoiceNumber(inv_number string) (*model.TransactionHeader, error)
}

type transactionUseCase struct {
	transactionRepo   repository.TransactionRepository
	customerRepo      repository.CustomerRepository
	meatRepo          repository.MeatRepository
	companyRepo       repository.CompanyRepository
	creditPaymentRepo repository.CreditPaymentRepository
}

// CreateTransaction implements TransactionUseCase.
func (uc *transactionUseCase) CreateTransaction(transaction *model.TransactionHeader) (*model.TransactionHeader, error) {
	// Generate invoice number
	tx := uc.transactionRepo.GetDB().Begin()
	defer tx.Rollback()
	today := time.Now().Format("20060102")
	todayDate := time.Now().Format("2006-01-02")
	number, err := uc.transactionRepo.CountTransactions()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to count transactions")
		return nil, err
	}

	customer, err := uc.customerRepo.GetCustomerById(transaction.CustomerID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to get customer by ID")
		return nil, fmt.Errorf("failed to get customer by id: %w", err)
	}
	company, err := uc.companyRepo.GetCompanyById(customer.CompanyId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to get company by ID")
		return nil, fmt.Errorf("failed to get company by id: %w", err)
	}
	invoiceNumberFormat := "MJP-%s-%04d"

	if transaction.TxType == "out" {
		invoiceNumberFormat = "INV-%s-%04d"
	}

	invoiceNumber := fmt.Sprintf(invoiceNumberFormat, today, number)
	transaction.ID = uuid.NewString()
	transaction.Date = todayDate
	transaction.Name = customer.FullName
	transaction.InvoiceNumber = invoiceNumber
	transaction.Address = customer.Address
	transaction.PhoneNumber = customer.PhoneNumber
	transaction.Company = company.CompanyName
	transaction.UpdatedBy = transaction.CreatedBy
	transaction.PaymentStatus = "paid"

	for _, detail := range transaction.TransactionDetails {
		meat, err := uc.meatRepo.GetMeatByName(detail.MeatName)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":     err,
				"meat_name": detail.MeatName,
			}).Error("Failed to get meat by name")
			return nil, err
		}
		if meat == nil {
			logrus.WithFields(logrus.Fields{
				"meat_name": detail.MeatName,
			}).Error("Meat not found by name")
			return nil, fmt.Errorf("meat name %s not found", detail.MeatName)
		}
		detail.ID = uuid.NewString()
		detail.MeatID = meat.ID
		detail.TransactionID = transaction.ID
		detail.IsActive = true

		if detail.Qty >= meat.Stock {
			logrus.WithFields(logrus.Fields{
				"meat_name": detail.MeatName,
			}).Error("Insufficient stock for meat")
			return nil, fmt.Errorf("insufficient stock for %s", detail.MeatName)
		}

		if transaction.TxType == "in" {
			err := uc.meatRepo.IncreaseStock(meat.ID, detail.Qty)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error":     err,
					"meat_id":   meat.ID,
					"meat_name": detail.MeatName,
				}).Error("Failed to increase meat stock")
				return nil, err
			}
		}
		if transaction.TxType == "out" {
			err = uc.meatRepo.ReduceStock(meat.ID, detail.Qty)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error":     err,
					"meat_id":   meat.ID,
					"meat_name": detail.MeatName,
				}).Error("Failed to reduce meat stock")
				return nil, err
			}
		}
	}
	transaction.CalulatedTotal()
	newTotal := uc.UpdateTotalTransaction(transaction)

	if transaction.PaymentAmount > newTotal {
		logrus.WithFields(logrus.Fields{
			"payment_amount": transaction.PaymentAmount,
			"new_total":      newTotal,
		}).Error("Amount greater than total")
		return nil, utils.ErrAmountGreaterThanTotal
	}

	if newTotal > transaction.PaymentAmount {
		transaction.PaymentStatus = "unpaid"
		transaction.Debt = newTotal - transaction.PaymentAmount
	}

	uc.creditPaymentRepo.CreateCreditPayment(&model.CreditPayment{
		ID:            uuid.New().String(),
		InvoiceNumber: transaction.InvoiceNumber,
		Amount:        transaction.PaymentAmount,
		PaymentDate:   transaction.Date,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.CreatedAt,
		CreatedBy:     transaction.CreatedBy,
		UpdatedBy:     transaction.CreatedBy,
	})

	// Create transaction header
	result, err := uc.transactionRepo.CreateTransactionHeader(transaction)
	if err != nil {
		tx.Rollback()
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to create transaction header")
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	err = uc.transactionRepo.UpdateCustomerDebt(transaction.CustomerID, transaction.Debt)

	if err != nil {
		tx.Rollback()
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to update customer debt")
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to commit transaction")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"invoice_number": transaction.InvoiceNumber,
		"username":       transaction.CreatedBy,
	}).Info("Transaction created successfully")
	return result, nil
}

func (uc *transactionUseCase) GetAllTransactions(page int, itemsPerPage int) ([]*model.TransactionHeader, int, error) {
	transactions, totalPages, err := uc.transactionRepo.GetAllTransactions(page, itemsPerPage)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to get all transactions")
		return nil, 0, err
	}
	logrus.WithFields(logrus.Fields{
		"page":           page,
		"items_per_page": itemsPerPage,
		"total_pages":    totalPages,
	}).Info("Retrieved all transactions successfully")
	return transactions, totalPages, nil
}

// notUse
func (uc *transactionUseCase) GetTransactionByID(id string) (*model.TransactionHeader, error) {
	transaction, err := uc.transactionRepo.GetTransactionByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

func (uc *transactionUseCase) DeleteTransaction(id string) error {
	return uc.transactionRepo.DeleteTransaction(id)
}

func (uc *transactionUseCase) UpdateTotalTransaction(transaction *model.TransactionHeader) float64 {
	var newTotal float64
	for _, detail := range transaction.TransactionDetails {
		detail.Total = detail.Price * detail.Qty
		newTotal = newTotal + detail.Total
	}

	return newTotal
}

func (uc *transactionUseCase) GetTransactionByInvoiceNumber(inv_number string) (*model.TransactionHeader, error) {
	transaction, err := uc.transactionRepo.GetByInvoiceNumber(inv_number)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

func NewTransactionUseCase(transactionRepo repository.TransactionRepository, customerRepo repository.CustomerRepository, meatRepo repository.MeatRepository, companyRepo repository.CompanyRepository, creditPaymentRepo repository.CreditPaymentRepository) TransactionUseCase {
	return &transactionUseCase{
		transactionRepo:   transactionRepo,
		customerRepo:      customerRepo,
		meatRepo:          meatRepo,
		companyRepo:       companyRepo,
		creditPaymentRepo: creditPaymentRepo,
	}
}
