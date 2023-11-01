package manager

import (
	"sync"
	"trackprosto/usecase"
)

type UsecaseManager interface {
	GetUserUsecase() usecase.UserUseCase
	GetLoginUsecase() usecase.LoginUseCase
	GetMeatUsecase() usecase.MeatUseCase
	GetTransactionUseCase() usecase.TransactionUseCase
	GetCreditPaymentUseCase() usecase.CreditPaymentUseCase
}

type usecaseManager struct {
	repoManager        RepoManager
	userUsecase        usecase.UserUseCase
	loginUsecase       usecase.LoginUseCase
	meatUsecase        usecase.MeatUseCase
	creditPaymentUseCase    usecase.CreditPaymentUseCase
	transactionUseCase usecase.TransactionUseCase
}

var onceLoadUserUsecase sync.Once
var onceLoadLoginUsecase sync.Once
var onceLoadMeatUsecase sync.Once
var onceLoadTxUsecase sync.Once
var onceLoadCreditPaymentUseCase sync.Once

func (um *usecaseManager) GetUserUsecase() usecase.UserUseCase {
	onceLoadUserUsecase.Do(func() {
		um.userUsecase = usecase.NewUserUseCase(um.repoManager.GetUserRepo())
	})
	return um.userUsecase
}

func (um *usecaseManager) GetCreditPaymentUseCase() usecase.CreditPaymentUseCase {
	onceLoadCreditPaymentUseCase.Do(func() {
		um.creditPaymentUseCase = usecase.NewCreditPaymentUseCase(um.repoManager.GetCreditPaymentRepo(), um.repoManager.GetTransactionRepo())
	})
	return um.creditPaymentUseCase
}

func (um *usecaseManager) GetLoginUsecase() usecase.LoginUseCase {
	onceLoadLoginUsecase.Do(func() {
		um.loginUsecase = usecase.NewLoginUseCase(um.repoManager.GetUserRepo())
	})
	return um.loginUsecase
}

func (mm *usecaseManager) GetMeatUsecase() usecase.MeatUseCase {
	onceLoadMeatUsecase.Do(func() {
		mm.meatUsecase = usecase.NewMeatUseCase(
			mm.repoManager.GetMeatRepo(),
			mm.repoManager.GetTransactionRepo())
	})
	return mm.meatUsecase
}

func (um *usecaseManager) GetTransactionUseCase() usecase.TransactionUseCase {
	onceLoadTxUsecase.Do(func() {
		um.transactionUseCase = usecase.NewTransactionUseCase(
			um.repoManager.GetTransactionRepo(),
			um.repoManager.GetCustomerRepo(),
			um.repoManager.GetMeatRepo(),
			um.repoManager.GetCompanyRepo(),
			um.repoManager.GetCreditPaymentRepo())
	})
	return um.transactionUseCase
}

func NewUsecaseManager(repoManager RepoManager) UsecaseManager {
	return &usecaseManager{
		repoManager: repoManager,
	}
}
