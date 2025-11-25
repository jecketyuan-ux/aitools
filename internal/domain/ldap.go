package domain

import (
	"time"
)

type LdapUser struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID        string     `json:"uuid" gorm:"uniqueIndex;size:255;not null"`
	OU          string     `json:"ou" gorm:"size:500"`
	CN          string     `json:"cn" gorm:"size:500"`
	DisplayName string     `json:"display_name" gorm:"size:500"`
	UserID      *int       `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (LdapUser) TableName() string {
	return "ldap_users"
}

type LdapDepartment struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID      string    `json:"uuid" gorm:"uniqueIndex;size:255;not null"`
	OU        string    `json:"ou" gorm:"size:500"`
	Name      string    `json:"name" gorm:"size:500"`
	DepID     *int      `json:"dep_id"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (LdapDepartment) TableName() string {
	return "ldap_departments"
}

type LdapSyncRecord struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Action    string     `json:"action" gorm:"size:50;not null"`
	StartAt   time.Time  `json:"start_at" gorm:"not null"`
	EndAt     *time.Time `json:"end_at"`
	Status    string     `json:"status" gorm:"size:50;not null"`
	ErrorMsg  string     `json:"error_msg" gorm:"type:text"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (LdapSyncRecord) TableName() string {
	return "ldap_sync_records"
}
