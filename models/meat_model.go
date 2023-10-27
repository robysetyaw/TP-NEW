package model

import "time"

type Meat struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Stock     float64   `json:"stock"`
	Price     float64   `json:"price"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}
