package dto

import "time"

type CompanyRequest struct {
	ID          string    `json:"id" gorm:"primary_key"`
	CompanyName string    `json:"company_name"`
	Address     string    `json:"address"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
}

type CompanyResponse struct {
	ID          string `json:"id" gorm:"primary_key"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}