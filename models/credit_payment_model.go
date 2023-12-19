package model

import "time"

type CreditPayment struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	InvoiceNumber string    `json:"inv_number" gorm:"column:inv_number"`
	PaymentDate   string    `json:"payment_date"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	UpdatedBy     string    `json:"updated_by"`
	Notes         string    `json:"notes"`
}

type CreditPaymentResponse struct {
	Transaction   *TransactionHeader
	CreditPayment *CreditPayment
}
