package model

import "time"

type CustomerModel struct {
	Id          string    `json:"-" gorm:"primaryKey"`
	FullName    string    `json:"fullname" binding:"required"`
	Address     string    `json:"address"`
	CompanyId   string    `json:"company_id"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
	Debt        float64   `json:"debt"`
}

func (CustomerModel) TableName() string {
	return "customers"
}
