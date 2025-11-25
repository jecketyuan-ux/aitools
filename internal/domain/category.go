package domain

import (
	"time"
)

type Category struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ParentID    int       `json:"parent_id" gorm:"default:0;index"`
	ParentChain string    `json:"parent_chain" gorm:"size:1000"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Sort        int       `json:"sort" gorm:"default:0"`
	IsShow      int       `json:"is_show" gorm:"default:1"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Category) TableName() string {
	return "categories"
}
