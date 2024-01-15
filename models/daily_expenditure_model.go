package model

import (
	"time"
)

type DailyExpenditure struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	DeNote      string    `json:"de_note"`
	Amount      float64   `json:"amount" binding:"required"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
	Date        string    `json:"date"`
}