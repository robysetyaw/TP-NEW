package manager

import (
	"sync"
	"trackprosto/repository"
)

type RepoManager interface {
	GetUserRepo() repository.UserRepository
	GetMeatRepo() repository.MeatRepository
}

type repoManager struct {
	infraManager InfraManager
	userRepo     repository.UserRepository
	meatRepo     repository.MeatRepository
}

var onceLoadUserRepo sync.Once
var onceLoadMeatRepo sync.Once

func (rm *repoManager) GetUserRepo() repository.UserRepository {
	onceLoadUserRepo.Do(func() {
		rm.userRepo = repository.NewUserRepository(rm.infraManager.GetDB())
	})
	return rm.userRepo
}

func (rm *repoManager) GetMeatRepo() repository.MeatRepository {
	onceLoadMeatRepo.Do(func() {
		rm.meatRepo = repository.NewMeatRepository(rm.infraManager.GetDB())
	})
	return rm.meatRepo
}

func NewRepoManager(infraManager InfraManager) RepoManager {
	return &repoManager{
		infraManager: infraManager,
	}
}
