package domain

import (
	"time"
)

type AppConfig struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	KeyName   string    `json:"key_name" gorm:"column:key_name;uniqueIndex;size:255;not null"`
	KeyValue  string    `json:"key_value" gorm:"column:key_value;type:text"`
	IsPrivate int       `json:"is_private" gorm:"default:0"`
	IsHidden  int       `json:"is_hidden" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (AppConfig) TableName() string {
	return "app_config"
}
