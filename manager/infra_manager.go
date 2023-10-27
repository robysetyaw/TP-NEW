package manager

import (
	"fmt"
	"sync"
	"trackprosto/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type InfraManager interface {
	GetDB() *gorm.DB
}
type infraManager struct {
	db  *gorm.DB
	cfg config.Config
}

var onceLoadDB sync.Once

func (i *infraManager) initDb() {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", i.cfg.Host, i.cfg.Port, i.cfg.User, i.cfg.Password, i.cfg.Name)
	onceLoadDB.Do(func() {
		db, err := gorm.Open(postgres.Open(psqlconn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		i.db = db
	})
	fmt.Println("DB Connected using GORM")
}

func (i *infraManager) GetDB() *gorm.DB {
	return i.db
}

func NewInfraManager(config config.Config) InfraManager {
	infra := infraManager{
		cfg: config,
	}
	infra.initDb()
	return &infra
}
