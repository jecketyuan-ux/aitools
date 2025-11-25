package domain

import (
	"time"
)

type Resource struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	AdminID    int       `json:"admin_id" gorm:"not null"`
	Type       string    `json:"type" gorm:"size:50;not null;index"`
	CategoryID int       `json:"category_id" gorm:"default:0;index"`
	URL        string    `json:"url" gorm:"size:2000;not null"`
	Name       string    `json:"name" gorm:"size:500;not null;index:idx_name"`
	Extension  string    `json:"extension" gorm:"size:50;not null"`
	Size       int64     `json:"size" gorm:"default:0"`
	Disk       string    `json:"disk" gorm:"size:50;not null"`
	FileID     string    `json:"file_id" gorm:"size:500"`
	Path       string    `json:"path" gorm:"size:2000"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Resource) TableName() string {
	return "resources"
}

type ResourceCategory struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ParentID    int       `json:"parent_id" gorm:"default:0;index"`
	ParentChain string    `json:"parent_chain" gorm:"size:1000"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Sort        int       `json:"sort" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (ResourceCategory) TableName() string {
	return "resource_categories"
}

type ResourceVideo struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	RID       int       `json:"rid" gorm:"not null;uniqueIndex"`
	Duration  int       `json:"duration" gorm:"default:0"`
	Poster    string    `json:"poster" gorm:"size:2000"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (ResourceVideo) TableName() string {
	return "resource_videos"
}
