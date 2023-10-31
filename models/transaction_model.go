package model

import "time"

// TransactionHeader adalah representasi dari tabel transaction_headers di database.
type TransactionHeader struct {
	ID                 string               `json:"id" gorm:"primaryKey" gorm:"tableName=transaction_headers"`
	Date               string               `json:"date"`
	InvoiceNumber      string               `json:"invoice_number" gorm:"column:inv_number"`
	CustomerID         string               `json:"customer_id"`
	Name               string               `json:"name"`
	Address            string               `json:"address"`
	Company            string               `json:"company"`
	PhoneNumber        string               `json:"phone_number"`
	TxType             string               `json:"tx_type"`
	PaymentStatus      string               `json:"payment_status"`
	PaymentAmount      float64              `json:"payment_amount"`
	Total              float64              `json:"total"`
	IsActive           bool                 `json:"is_active"`
	CreatedAt          time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time            `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy          string               `json:"created_by"`
	UpdatedBy          string               `json:"updated_by"`
	Debt               float64              `json:"debt" gorm:"column:debt"`
	TransactionDetails []*TransactionDetail `json:"transaction_details" gorm:"foreignKey:TransactionID"`
}

// TransactionDetail adalah representasi dari detail transaksi.
type TransactionDetail struct {
	ID            string    `json:"id" gorm:"primaryKey" gorm:"tableName=transaction_details"`
	TransactionID string    `json:"transaction_id"`
	MeatID        string    `json:"meat_id"`
	MeatName      string    `json:"meat_name"`
	Qty           float64   `json:"qty"`
	Price         float64   `json:"price"`
	Total         float64   `json:"total"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy     string    `json:"created_by"`
	UpdatedBy     string    `json:"updated_by"`
}

// CalulatedTotal menghitung total transaksi berdasarkan detail transaksi.
func (h *TransactionHeader) CalulatedTotal() {
	total := 0.0
	for _, detail := range h.TransactionDetails {
		detail.Total = detail.Price * detail.Qty
		total += detail.Total
	}
	h.Total = total
}

func (TransactionHeader) TableName() string {
	return "transaction_headers"
}

func (TransactionDetail) TableName() string {
	return "transaction_details"
}
