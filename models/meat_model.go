package model

import "time"

type Meat struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Stock     float64   `json:"stock"`
	Price     float64   `json:"price"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}

type MeatWithStock struct {
	Meat     *Meat
	StockIn  float64
	StockOut float64
}
