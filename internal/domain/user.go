package domain

import (
	"time"
)

type User struct {
	ID            int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Email         string     `json:"email" gorm:"uniqueIndex;size:255;not null"`
	Name          string     `json:"name" gorm:"size:255;not null"`
	Avatar        *int       `json:"avatar" gorm:"index"`
	Password      string     `json:"-" gorm:"size:255"`
	Salt          string     `json:"-" gorm:"size:255"`
	IDCard        string     `json:"id_card" gorm:"column:id_card;size:255"`
	Credit        int        `json:"credit" gorm:"default:0"`
	CreateIP      string     `json:"create_ip" gorm:"size:255"`
	CreateCity    string     `json:"create_city" gorm:"size:255"`
	IsActive      int        `json:"is_active" gorm:"default:1"`
	IsLock        int        `json:"is_lock" gorm:"default:0"`
	IsVerify      int        `json:"is_verify" gorm:"default:0"`
	VerifyAt      *time.Time `json:"verify_at"`
	IsSetPassword int        `json:"is_set_password" gorm:"default:0"`
	LoginAt       *time.Time `json:"login_at"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

type Department struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	ParentID    int       `json:"parent_id" gorm:"default:0;index"`
	ParentChain string    `json:"parent_chain" gorm:"size:1000"`
	Sort        int       `json:"sort" gorm:"default:0;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Department) TableName() string {
	return "departments"
}

type UserDepartment struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       int       `json:"user_id" gorm:"not null;uniqueIndex:idx_user_dep"`
	DepartmentID int       `json:"department_id" gorm:"not null;uniqueIndex:idx_user_dep;index"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (UserDepartment) TableName() string {
	return "user_department"
}

type UserLoginRecord struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         int        `json:"user_id" gorm:"not null;index"`
	JTI            string     `json:"jti" gorm:"size:255;not null;index"`
	IP             string     `json:"ip" gorm:"size:255"`
	IPArea         string     `json:"ip_area" gorm:"size:500"`
	Browser        string     `json:"browser" gorm:"size:500"`
	BrowserVersion string     `json:"browser_version" gorm:"size:255"`
	OS             string     `json:"os" gorm:"size:255"`
	IsLogout       int        `json:"is_logout" gorm:"default:0"`
	ExpiredAt      time.Time  `json:"expired_at" gorm:"not null"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (UserLoginRecord) TableName() string {
	return "user_login_records"
}

type UserUploadImageLog struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id" gorm:"not null;index"`
	Scene     string    `json:"scene" gorm:"size:255;not null"`
	RID       int       `json:"rid" gorm:"not null"`
	IP        string    `json:"ip" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (UserUploadImageLog) TableName() string {
	return "user_upload_image_logs"
}
