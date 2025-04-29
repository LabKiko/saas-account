package model

import (
	"gorm.io/gorm"
)

// Base 是所有模型的基础结构
type Base struct {
	ID        int64          `gorm:"primarykey" json:"id"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除
}
