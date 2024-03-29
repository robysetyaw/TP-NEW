package usecase

import (
	"time"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/repository"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type CreditPaymentUseCase interface {
	CreateCreditPayment(payment *model.CreditPayment) (*model.CreditPaymentResponse, error)
	GetCreditPayments() ([]*model.CreditPayment, error)
	GetCreditPaymentByID(id string) (*model.CreditPayment, error)
	UpdateCreditPayment(payment *model.CreditPayment) error
	GetCreditPaymentsByInvoiceNumber(inv_number string) ([]*model.CreditPayment, error)
}
type creditPaymentUseCase struct {
	creditPaymentRepo    repository.CreditPaymentRepository
	transactionRepo      repository.TransactionRepository
	dailyExpenditureRepo repository.DailyExpenditureRepository
}

func NewCreditPaymentUseCase(creditPaymentRepo repository.CreditPaymentRepository, transactionRepo repository.TransactionRepository, dailyExpenditureRepo repository.DailyExpenditureRepository) CreditPaymentUseCase {
	return &creditPaymentUseCase{
		creditPaymentRepo: creditPaymentRepo,
		transactionRepo:   transactionRepo,
		dailyExpenditureRepo: dailyExpenditureRepo,
	}
}

func (uc *creditPaymentUseCase) CreateCreditPayment(payment *model.CreditPayment) (*model.CreditPaymentResponse, error) {

	tx := uc.transactionRepo.GetDB().Begin() // Start a transaction
	defer tx.Rollback()

	transaction, err := uc.transactionRepo.GetByInvoiceNumber(payment.InvoiceNumber)
	if err != nil {
		log.WithFields(log.Fields{
			"invoiceNumber": payment.InvoiceNumber,
			"error":         err,
		}).Error("Failed to get transaction by invoice number")
		return nil, err
	}
	if transaction == nil {
		log.WithField("invoiceNumber", payment.InvoiceNumber).Error("Invoice not found")
		return nil, utils.ErrInvoiceNumberNotExist
	}
	if transaction.PaymentStatus == "paid" {
		log.WithField("invoiceNumber", payment.InvoiceNumber).Error("Invoice has already been paid.")
		return nil, utils.ErrInvoiceAlreadyPaid
	}

	CountCreditPayments, err := uc.creditPaymentRepo.CountCreditPayments(payment.InvoiceNumber)
	if err != nil {
		log.WithFields(log.Fields{
			"invoiceNumber": payment.InvoiceNumber,
			"error":         err,
		}).Error("Failed to count credit payment")
		return nil, err
	}
	totalcount := CountCreditPayments

	createdat := time.Now()
	todayDate := time.Now().Format("2006-01-02")
	payment.ID = uuid.NewString()
	payment.PaymentDate = todayDate
	payment.CreatedAt = createdat
	payment.UpdatedAt = createdat
	payment.UpdatedBy = payment.CreatedBy
	payment.Notes = utils.NumberToOrdinal(totalcount+1) + " Installment"

	totalCredit, err := uc.creditPaymentRepo.GetTotalCredit(payment.InvoiceNumber)
	if err != nil {
		tx.Rollback() // Rollback in case of error
		log.WithFields(log.Fields{
			"invoiceNumber": payment.InvoiceNumber,
			"error":         err,
		}).Error("Failed to get total credit")
		return nil, err
	}
	totalAmountAfterCredit := totalCredit + payment.Amount
	if totalAmountAfterCredit >= transaction.Total {
		payment.Notes = "Settled"
		err = uc.transactionRepo.UpdateStatusInvoicePaid(transaction.ID)
		if err != nil {
			tx.Rollback() // Rollback in case of error
			log.WithFields(log.Fields{
				"invoiceNumber": payment.InvoiceNumber,
				"error":         err,
			}).Error("Failed to update invoice payment status")
			return nil, err
		}
	}
	if totalAmountAfterCredit > transaction.Total {
		tx.Rollback() // Rollback in case of error
		return nil, utils.ErrAmountGreaterThanTotal
	}

	if transaction.TxType == "in" {
		uc.dailyExpenditureRepo.CreateDailyExpenditure(&model.DailyExpenditure{
			ID:         uuid.NewString(),
			Date:       todayDate,
			DeNote:     payment.InvoiceNumber,
			Amount:     payment.Amount,
			IsActive:   true,
			CreatedAt:  createdat,
			UpdatedAt:  createdat,
			CreatedBy:  payment.CreatedBy,
			UpdatedBy:  payment.CreatedBy,
			Description: payment.Notes,
		})
	}

	err = uc.creditPaymentRepo.CreateCreditPayment(payment)
	if err != nil {
		tx.Rollback() // Rollback in case of error
		log.WithFields(log.Fields{
			"invoiceNumber": payment.InvoiceNumber,
			"error":         err,
		}).Error("Failed to create credit payment")
		return nil, err
	}

	newDebt := transaction.Total - totalAmountAfterCredit
	uc.transactionRepo.UpdateStatusPaymentAmount(transaction.ID, totalAmountAfterCredit)
	uc.transactionRepo.UpdateDebtTransaction(transaction.ID, newDebt)

	// Update the transaction after the payment amount update
	transaction.PaymentAmount = totalCredit
	transaction.Debt = newDebt

	creditPaymentResponse := &model.CreditPaymentResponse{
		Transaction:   transaction,
		CreditPayment: payment,
	}

	if err := tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{
			"invoiceNumber": payment.InvoiceNumber,
			"error":         err,
		}).Error("Failed to complete the transaction")
		return nil, err
	}

	return creditPaymentResponse, nil
}

func (uc *creditPaymentUseCase) GetCreditPayments() ([]*model.CreditPayment, error) {
	payments, err := uc.creditPaymentRepo.GetAllCreditPayments()
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (uc *creditPaymentUseCase) GetCreditPaymentByID(id string) (*model.CreditPayment, error) {
	payment, err := uc.creditPaymentRepo.GetCreditPaymentByID(id)
	if payment == nil {
		return nil, utils.ErrCreditPaymentNotFound
	}
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (uc *creditPaymentUseCase) UpdateCreditPayment(payment *model.CreditPayment) error {

	err := uc.creditPaymentRepo.UpdateCreditPayment(payment)
	if err != nil {
		return err
	}

	return nil
}

func (uc *creditPaymentUseCase) GetCreditPaymentsByInvoiceNumber(inv_number string) ([]*model.CreditPayment, error) {
	payments, err := uc.creditPaymentRepo.GetCreditPaymentsByInvoiceNumber(inv_number)
	if len(payments) == 0 {
		return nil, utils.ErrCreditPaymentNotFound
	}
	if err != nil {
		return nil, err
	}

	return payments, nil
}
