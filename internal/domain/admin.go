package domain

import (
	"time"
)

type AdminUser struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"size:255;not null"`
	Email       string     `json:"email" gorm:"uniqueIndex;size:255;not null"`
	Password    string     `json:"-" gorm:"size:255;not null"`
	Salt        string     `json:"-" gorm:"size:255;not null"`
	LoginIP     string     `json:"login_ip" gorm:"size:255"`
	LoginAt     *time.Time `json:"login_at"`
	IsBanLogin  int        `json:"is_ban_login" gorm:"default:0"`
	LoginTimes  int        `json:"login_times" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (AdminUser) TableName() string {
	return "admin_users"
}

type AdminRole struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	Slug      string    `json:"slug" gorm:"uniqueIndex;size:255;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (AdminRole) TableName() string {
	return "admin_roles"
}

type AdminPermission struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Type      string    `json:"type" gorm:"size:50;not null"`
	GroupName string    `json:"group_name" gorm:"size:255;not null"`
	Sort      int       `json:"sort" gorm:"default:0"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	Slug      string    `json:"slug" gorm:"uniqueIndex;size:255;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (AdminPermission) TableName() string {
	return "admin_permissions"
}

type AdminRolePermission struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleID       int       `json:"role_id" gorm:"not null;uniqueIndex:idx_role_perm"`
	PermissionID int       `json:"permission_id" gorm:"not null;uniqueIndex:idx_role_perm"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (AdminRolePermission) TableName() string {
	return "admin_role_permission"
}

type AdminUserRole struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	AdminID   int       `json:"admin_id" gorm:"not null;uniqueIndex:idx_admin_role"`
	RoleID    int       `json:"role_id" gorm:"not null;uniqueIndex:idx_admin_role"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (AdminUserRole) TableName() string {
	return "admin_user_role"
}

type AdminLog struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	AdminID   int       `json:"admin_id" gorm:"not null;index"`
	AdminName string    `json:"admin_name" gorm:"size:255;not null"`
	Module    string    `json:"module" gorm:"size:255;not null"`
	Title     string    `json:"title" gorm:"size:500;not null"`
	Opt       string    `json:"opt" gorm:"size:50;not null"`
	Method    string    `json:"method" gorm:"size:50;not null"`
	URL       string    `json:"url" gorm:"size:2000;not null"`
	IP        string    `json:"ip" gorm:"size:255;not null"`
	IPArea    string    `json:"ip_area" gorm:"size:500"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;index"`
}

func (AdminLog) TableName() string {
	return "admin_logs"
}
