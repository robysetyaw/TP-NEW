package usecase

import (
	"time"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/repository"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

type CreditPaymentUseCase interface {
	CreateCreditPayment(payment *model.CreditPayment) (*model.CreditPaymentResponse, error)
	GetCreditPayments() ([]*model.CreditPayment, error)
	GetCreditPaymentByID(id string) (*model.CreditPayment, error)
	UpdateCreditPayment(payment *model.CreditPayment) error
	GetCreditPaymentsByInvoiceNumber(inv_number string) ([]*model.CreditPayment, error)
}
type creditPaymentUseCase struct {
	creditPaymentRepo repository.CreditPaymentRepository
	transactionRepo   repository.TransactionRepository
}

func NewCreditPaymentUseCase(creditPaymentRepo repository.CreditPaymentRepository, transactionRepo repository.TransactionRepository) CreditPaymentUseCase {
	return &creditPaymentUseCase{
		creditPaymentRepo: creditPaymentRepo,
		transactionRepo:   transactionRepo,
	}
}

func (uc *creditPaymentUseCase) CreateCreditPayment(payment *model.CreditPayment) (*model.CreditPaymentResponse, error) {
	// Validasi atau logika bisnis sebelum membuat pembayaran kredit
	// ...
	log.Info("Mengawali pembuatan pembayaran kredit")
	tx := uc.transactionRepo.GetDB().Begin() // Mulai transaction
	defer tx.Rollback()

	transaction, err := uc.transactionRepo.GetByInvoiceNumber(payment.InvoiceNumber)
	if err != nil {
		log.Error("Gagal mengambil transaksi berdasarkan nomor faktur:", err)
		return nil, err
	}
	if transaction == nil {
		log.Error("Transaksi tidak ditemukan berdasarkan nomor faktur:", payment.InvoiceNumber)
		// log.WithFields(log.Fields{
		// 	"transaction": payment.InvoiceNumber,
		// 	"nom":         payment.Amount,
		// }).Info("Transaksi tidak ditemukan berdasarkan nomor faktur")

		return nil, utils.ErrInvoiceNumberNotExist
	}
	if transaction.PaymentStatus == "paid" {
		log.Error("Faktur sudah dibayar sebelumnya.")
		return nil, utils.ErrInvoiceAlreadyPaid
	}
	createdat := time.Now()
	todayDate := time.Now().Format("2006-01-02")
	payment.ID = uuid.NewString()
	payment.PaymentDate = todayDate
	payment.CreatedAt = createdat
	payment.UpdatedAt = createdat
	payment.CreatedBy = "admin"
	payment.UpdatedBy = "admin"

	err = uc.creditPaymentRepo.CreateCreditPayment(payment)
	if err != nil {
		log.Error("Gagal membuat pembayaran kredit:", err)
		tx.Rollback() // Rollback jika terjadi kesalahan
		return nil, err
	}
	totalCredit, err := uc.creditPaymentRepo.GetTotalCredit(payment.InvoiceNumber)
	if err != nil {
		log.Error("Gagal mengambil total kredit:", err)
		tx.Rollback() // Rollback jika terjadi kesalahan
		return nil, err
	}
	newDebt := transaction.Total - totalCredit
	uc.transactionRepo.UpdateStatusPaymentAmount(transaction.ID, totalCredit)
	uc.transactionRepo.UpdateDebtTransaction(transaction.ID, newDebt)

	if totalCredit >= transaction.Total {
		err = uc.transactionRepo.UpdateStatusInvoicePaid(transaction.ID)
		if err != nil {
			log.Error("Gagal memperbarui status faktur menjadi sudah dibayar:", err)
			tx.Rollback() // Rollback jika terjadi kesalahan
			return nil, err
		}
	}
	// Perbarui transaksi setelah pembaruan payment amount
	transaction.PaymentAmount = totalCredit
	transaction.Debt = newDebt

	creditPaymentResponse := &model.CreditPaymentResponse{
		Transaction:   transaction,
		CreditPayment: payment,
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("Gagal melakukan commit transaksi:", err)
		return nil, err
	}

	log.Info("Pembayaran kredit berhasil dibuat.")
	return creditPaymentResponse, nil
}

func (uc *creditPaymentUseCase) GetCreditPayments() ([]*model.CreditPayment, error) {
	payments, err := uc.creditPaymentRepo.GetAllCreditPayments()
	if err != nil {
		return nil, err
	}
	log.Info("Pembayaran kredit selesai")
	return payments, nil
}

func (uc *creditPaymentUseCase) GetCreditPaymentByID(id string) (*model.CreditPayment, error) {
	payment, err := uc.creditPaymentRepo.GetCreditPaymentByID(id)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (uc *creditPaymentUseCase) UpdateCreditPayment(payment *model.CreditPayment) error {
	// Validasi atau logika bisnis sebelum memperbarui pembayaran kredit
	// ...

	err := uc.creditPaymentRepo.UpdateCreditPayment(payment)
	if err != nil {
		return err
	}

	return nil
}

func (uc *creditPaymentUseCase) GetCreditPaymentsByInvoiceNumber(inv_number string) ([]*model.CreditPayment, error) {
	payments, err := uc.creditPaymentRepo.GetCreditPaymentsByInvoiceNumber(inv_number)
	if err != nil {
		return nil, err
	}

	return payments, nil
}
