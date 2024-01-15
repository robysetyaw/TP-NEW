package manager

import (
	"sync"
	"trackprosto/repository"
)

type RepoManager interface {
	GetUserRepo() repository.UserRepository
	GetMeatRepo() repository.MeatRepository
	GetCustomerRepo() repository.CustomerRepository
	GetCompanyRepo() repository.CompanyRepository
	GetTransactionRepo() repository.TransactionRepository
	GetCreditPaymentRepo() repository.CreditPaymentRepository
	GetDailyExpenditureRepo() repository.DailyExpenditureRepository
}

type repoManager struct {
	infraManager         InfraManager
	customerRepo         repository.CustomerRepository
	userRepo             repository.UserRepository
	meatRepo             repository.MeatRepository
	companyRepo          repository.CompanyRepository
	transactionRepo      repository.TransactionRepository
	creditPaymentRepo    repository.CreditPaymentRepository
	dailyExpenditureRepo repository.DailyExpenditureRepository
}

// GetCustomerRepo implements RepoManager.

var onceLoadUserRepo sync.Once
var onceLoadMeatRepo sync.Once
var onceLoadCustomerRepo sync.Once
var onceLoadCompanyRepo sync.Once
var onceLoadTxRepo sync.Once
var onceLoadCreditPaymentRepo sync.Once
var onceLoadDailyExpenditureRepo sync.Once

func (rm *repoManager) GetDailyExpenditureRepo() repository.DailyExpenditureRepository {
	onceLoadDailyExpenditureRepo.Do(func() {
		rm.dailyExpenditureRepo = repository.NewDailyExpenditureRepository(rm.infraManager.GetDB())
	})
	return rm.dailyExpenditureRepo
}

func (rm *repoManager) GetUserRepo() repository.UserRepository {
	onceLoadUserRepo.Do(func() {
		rm.userRepo = repository.NewUserRepository(rm.infraManager.GetDB())
	})
	return rm.userRepo
}

func (rm *repoManager) GetCustomerRepo() repository.CustomerRepository {
	onceLoadCustomerRepo.Do(func() {
		rm.customerRepo = repository.NewCustomerRepository(rm.infraManager.GetDB())
	})
	return rm.customerRepo
}

func (rm *repoManager) GetMeatRepo() repository.MeatRepository {
	onceLoadMeatRepo.Do(func() {
		rm.meatRepo = repository.NewMeatRepository(rm.infraManager.GetDB())
	})
	return rm.meatRepo
}

func (rm *repoManager) GetTransactionRepo() repository.TransactionRepository {
	onceLoadTxRepo.Do(func() {
		rm.transactionRepo = repository.NewTransactionRepository(rm.infraManager.GetDB())
	})
	return rm.transactionRepo
}

func (rm *repoManager) GetCreditPaymentRepo() repository.CreditPaymentRepository {
	onceLoadCreditPaymentRepo.Do(func() {
		rm.creditPaymentRepo = repository.NewCreditPaymentRepository(rm.infraManager.GetDB())
	})
	return rm.creditPaymentRepo
}

func (rm *repoManager) GetCompanyRepo() repository.CompanyRepository {
	onceLoadCompanyRepo.Do(func() {
		rm.companyRepo = repository.NewCompanyRepository(rm.infraManager.GetDB())
	})
	return rm.companyRepo
}

func NewRepoManager(infraManager InfraManager) RepoManager {
	return &repoManager{
		infraManager: infraManager,
	}
}
