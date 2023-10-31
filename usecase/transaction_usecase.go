package usecase

import (
	"fmt"
	"time"
	model "trackprosto/models"
	"trackprosto/repository"

	"github.com/google/uuid"
)

type TransactionUseCase interface {
	CreateTransaction(transaction *model.TransactionHeader) (*model.TransactionHeader, error)
	GetAllTransactions() ([]*model.TransactionHeader, error)
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

	// TODO countTransaction still return 0
	today := time.Now().Format("20060102")
	todayDate := time.Now().Format("2006-01-02")
	number, err := uc.transactionRepo.CountTransactions()
	if err != nil {
		return nil, err
	}

	customer, err := uc.customerRepo.GetCustomerByName(transaction.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer by name: %w", err)
	}

	company, err := uc.companyRepo.GetCompanyById(customer.CompanyId)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer by name: %w", err)
	}

	// TODO cek invoicenumebr
	invoiceNumber := fmt.Sprintf("INV-%s-%04d", today, number)
	transaction.ID = uuid.NewString()
	transaction.Date = todayDate
	transaction.InvoiceNumber = invoiceNumber
	transaction.CustomerID = customer.Id
	// transaction.Name = transaction.Name
	transaction.Address = customer.Address
	transaction.PhoneNumber = customer.PhoneNumber
	transaction.Company = company.CompanyName
	transaction.CreatedBy = "admin"
	transaction.UpdatedBy = "admin"
	transaction.PaymentStatus = "paid"

	for _, detail := range transaction.TransactionDetails {
		meat, err := uc.meatRepo.GetMeatByName(detail.MeatName)
		if err != nil {
			return nil, err
		}
		if meat == nil {
			return nil, fmt.Errorf("meat name %s not found", meat.Name)
		}
		detail.ID = uuid.NewString()
		detail.MeatID = meat.ID
		detail.TransactionID = transaction.ID

		if detail.Qty >= meat.Stock {
			return nil, fmt.Errorf("insufficient stock for %s", detail.MeatName)
		}

		if transaction.TxType == "in" {
			err := uc.meatRepo.IncreaseStock(meat.ID, detail.Qty)
			if err != nil {
				return nil, err
			}
		}
		if transaction.TxType == "out" {
			err = uc.meatRepo.ReduceStock(meat.ID, detail.Qty)
			if err != nil {
				return nil, err
			}
		}

	}
	transaction.CalulatedTotal()
	newTotal := uc.UpdateTotalTransaction(transaction)

	if transaction.PaymentAmount > newTotal {
		return nil, fmt.Errorf("amount is large than total transaction")
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
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	err = uc.transactionRepo.UpdateCustomerDebt(transaction.CustomerID, transaction.Debt)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (uc *transactionUseCase) GetAllTransactions() ([]*model.TransactionHeader, error) {
	return uc.transactionRepo.GetAllTransactions()
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
