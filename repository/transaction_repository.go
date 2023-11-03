package repository

import (
	"errors"
	"time"
	model "trackprosto/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransactionHeader(header *model.TransactionHeader) (*model.TransactionHeader, error)
	GetTransactionByID(id string) (*model.TransactionHeader, error)
	GetTransactionByRangeDate(startDate time.Time, endDate time.Time) ([]*model.TransactionHeader, error)
	GetAllTransactions(page int, itemsPerPage int) ([]*model.TransactionHeader, int, error)
	DeleteTransaction(id string) error
	CountTransactions() (int, error)
	GetByInvoiceNumber(invoice_number string) (*model.TransactionHeader, error)
	UpdateStatusInvoicePaid(id string) error
	UpdateStatusPaymentAmount(id string, total float64) error
	GetTransactionByRangeDateWithTxType(startDate time.Time, endDate time.Time, tx_type string) ([]*model.TransactionHeader, error)
	GetTransactionByRangeDateWithTxTypeAndPaid(startDate time.Time, endDate time.Time, tx_type, payment_status string) ([]*model.TransactionHeader, error)
	GetTransactionsByDateAndType(startDate time.Time, endDate time.Time, txType string) ([]*model.TransactionHeader, error)
	GetAllTransactionsByCustomerUsername(customer_id string) ([]*model.TransactionHeader, error)
	getCustomerDebt(customer_id string) (float64, error)
	getTransactionDebt(id string) (float64, error)
	CalculateMeatStockByDate(meatID string, startDate string) (stockIn float64, stockOut float64, err error)
	UpdateCustomerDebt(id string, additionalDebt float64) error
	GetDB() *gorm.DB
	UpdateDebtTransaction(id string, total float64) error
}

type transactionRepository struct {
	db *gorm.DB
}

func (repo *transactionRepository) GetDB() *gorm.DB {
	return repo.db
}
func (repo *transactionRepository) CalculateMeatStockByDate(meatID string, startDate string) (float64, float64, error) {
	var stockIn, stockOut float64

	// Calculate stockIn
	stockIn, err := repo.calculateMeatStockByType(meatID, startDate, "in")
	if err != nil {
		return 0, 0, err
	}

	// Calculate stockOut
	stockOut, err = repo.calculateMeatStockByType(meatID, startDate, "out")
	if err != nil {
		return 0, 0, err
	}

	return stockIn, stockOut, nil
}

func (repo *transactionRepository) calculateMeatStockByType(meatID string, startDate string, txType string) (float64, error) {
	var transactionDetails []*model.TransactionDetail

	err := repo.db.Joins("JOIN transaction_headers ON transaction_headers.id = transaction_details.transaction_id").
		Where("transaction_headers.date >= ? AND transaction_details.meat_id = ? AND transaction_headers.tx_type = ?", startDate, meatID, txType).
		Find(&transactionDetails).Error

	if err != nil {
		return 0, err
	}

	totalQty := 0.
	for _, detail := range transactionDetails {
		totalQty += detail.Qty
	}

	return totalQty, nil
}

// GetAllTransactionsByCustomerUsername implements TransactionRepository.
func (*transactionRepository) GetAllTransactionsByCustomerUsername(customer_id string) ([]*model.TransactionHeader, error) {
	panic("unimplemented")
}

// getCustomerDebt implements TransactionRepository.
func (*transactionRepository) getCustomerDebt(customer_id string) (float64, error) {
	panic("unimplemented")
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (repo *transactionRepository) CreateTransactionHeader(header *model.TransactionHeader) (*model.TransactionHeader, error) {
	now := time.Now()
	header.CreatedAt = now
	header.UpdatedAt = now
	header.IsActive = true

	if err := repo.db.Create(header).Error; err != nil {
		return nil, err
	}

	return header, nil
}

func (repo *transactionRepository) GetTransactionByID(id string) (*model.TransactionHeader, error) {
	var transaction model.TransactionHeader

	if err := repo.db.Preload("TransactionDetails").Where("id = ? AND is_active = true", id).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &transaction, nil
}

func (repo *transactionRepository) GetAllTransactions(page int, itemsPerPage int) ([]*model.TransactionHeader, int, error) {
	var transactions []*model.TransactionHeader

	if page < 1 {
		page = 1
	}

	var totalCount int64
	if err := repo.db.Model(&model.TransactionHeader{}).Where("is_active = true").Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	totalPages := int((totalCount + int64(itemsPerPage) - 1) / int64(itemsPerPage))

	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * itemsPerPage

	err := repo.db.Preload("TransactionDetails").Where("is_active = true").
		Offset(offset).Limit(itemsPerPage).
		Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return nil, totalPages, err
	}

	return transactions, totalPages, nil
}


func (repo *transactionRepository) DeleteTransaction(id string) error {
	transaction := model.TransactionHeader{ID: id}

	if err := repo.db.Model(&transaction).Update("is_active", false).Error; err != nil {
		return err
	}

	return nil
}

func (repo *transactionRepository) GetTransactionByRangeDate(startDate time.Time, endDate time.Time) ([]*model.TransactionHeader, error) {
	var transactions []*model.TransactionHeader

	err := repo.db.Preload("TransactionDetails").Where("created_at BETWEEN ? AND ? AND is_active = true", startDate, endDate).Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (repo *transactionRepository) GetByInvoiceNumber(invoice_number string) (*model.TransactionHeader, error) {
	var transaction model.TransactionHeader

	err := repo.db.Preload("TransactionDetails").Where("inv_number = ? AND is_active = true", invoice_number).First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &transaction, nil
}

func (repo *transactionRepository) UpdateStatusInvoicePaid(id string) error {
	return repo.db.Model(&model.TransactionHeader{}).Where("id = ?", id).Update("payment_status", "paid").Error
}

func (repo *transactionRepository) UpdateCustomerDebt(id string, additionalDebt float64) error {
	var customer model.CustomerModel
	// Mengambil data pelanggan berdasarkan ID
	if err := repo.db.Where("id = ?", id).First(&customer).Error; err != nil {
		return err
	}
	// Mengambil total utang saat ini
	currentDebt := customer.Debt

	// Menambahkan utang tambahan
	newDebt := currentDebt + additionalDebt
	if err := repo.db.Model(&customer).Update("debt", newDebt).Error; err != nil {
		return err
	}
	return nil
}

func (repo *transactionRepository) UpdateStatusPaymentAmount(id string, total float64) error {
	return repo.db.Model(&model.TransactionHeader{}).Where("id = ?", id).Update("payment_amount", total).Error
}

func (repo *transactionRepository) UpdateDebtTransaction(id string, total float64) error {
	return repo.db.Model(&model.TransactionHeader{}).Where("id = ?", id).Update("debt", total).Error
}

func (repo *transactionRepository) GetTransactionByRangeDateWithTxType(startDate time.Time, endDate time.Time, tx_type string) ([]*model.TransactionHeader, error) {
	var transactions []*model.TransactionHeader

	err := repo.db.Preload("TransactionDetails").Where("created_at BETWEEN ? AND ? AND is_active = true AND tx_type = ?", startDate, endDate, tx_type).Order("created_at ASC").Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (repo *transactionRepository) CountTransactions() (int, error) {
	var count int64
	today := time.Now()
	err := repo.db.Where("date = ? ", today.Format("2006-01-02")).Model(&model.TransactionHeader{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *transactionRepository) GetTransactionByRangeDateWithTxTypeAndPaid(startDate time.Time, endDate time.Time, tx_type, payment_status string) ([]*model.TransactionHeader, error) {
	var transactions []*model.TransactionHeader

	err := repo.db.Preload("TransactionDetails").Where("created_at BETWEEN ? AND ? AND is_active = true AND tx_type = ? AND payment_status = ?", startDate, endDate, tx_type, payment_status).Order("created_at ASC").Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (repo *transactionRepository) GetTransactionsByDateAndType(startDate time.Time, endDate time.Time, txType string) ([]*model.TransactionHeader, error) {
	var transactions []*model.TransactionHeader

	err := repo.db.Preload("TransactionDetails").Where("date BETWEEN ? AND ? AND is_active = true AND tx_type = ?", startDate, endDate, txType).Order("created_at ASC").Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (repo *transactionRepository) getTransactionDebt(id string) (float64, error) {
	var transaction model.TransactionHeader

	err := repo.db.Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return 0, err
	}

	return transaction.Total - transaction.PaymentAmount, nil
}
