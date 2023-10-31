package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:uuid;primary_key;" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username" binding:"required"`
	Password  string         `json:"password" binding:"required"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	Role      string         `json:"role"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy string         `json:"created_by"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Jika Anda ingin soft delete
}
