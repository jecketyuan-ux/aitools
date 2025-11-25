package domain

import (
	"time"
)

type Course struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string     `json:"title" gorm:"size:500;not null;index:idx_title"`
	Thumb       *int       `json:"thumb"`
	Charge      int        `json:"charge" gorm:"default:0"`
	ShortDesc   string     `json:"short_desc" gorm:"type:text"`
	IsRequired  int        `json:"is_required" gorm:"default:0"`
	ClassHour   int        `json:"class_hour" gorm:"default:0"`
	IsShow      int        `json:"is_show" gorm:"default:1;index"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	SortAt      *time.Time `json:"sort_at" gorm:"index"`
	PublishedAt *time.Time `json:"published_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}

func (Course) TableName() string {
	return "courses"
}

type CourseChapter struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CourseID  int       `json:"course_id" gorm:"not null;index"`
	Name      string    `json:"name" gorm:"size:500;not null"`
	Sort      int       `json:"sort" gorm:"default:0;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (CourseChapter) TableName() string {
	return "course_chapters"
}

type CourseHour struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	CourseID    int        `json:"course_id" gorm:"not null;index"`
	ChapterID   int        `json:"chapter_id" gorm:"default:0;index"`
	Sort        int        `json:"sort" gorm:"default:0"`
	Title       string     `json:"title" gorm:"size:500;not null"`
	Type        string     `json:"type" gorm:"size:50;not null"`
	RID         int        `json:"rid" gorm:"not null;index"`
	Duration    int        `json:"duration" gorm:"default:0"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (CourseHour) TableName() string {
	return "course_hours"
}

type CourseCategory struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CourseID   int       `json:"course_id" gorm:"not null;uniqueIndex:idx_course_cat"`
	CategoryID int       `json:"category_id" gorm:"not null;uniqueIndex:idx_course_cat"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (CourseCategory) TableName() string {
	return "course_categories"
}

type CourseAttachment struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CourseID  int       `json:"course_id" gorm:"not null;index"`
	Sort      int       `json:"sort" gorm:"default:0"`
	Title     string    `json:"title" gorm:"size:500;not null"`
	Type      string    `json:"type" gorm:"size:50;not null"`
	RID       int       `json:"rid" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (CourseAttachment) TableName() string {
	return "course_attachments"
}

type CourseAttachmentDownloadLog struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       int       `json:"user_id" gorm:"not null;index"`
	CourseID     int       `json:"course_id" gorm:"not null;index"`
	Title        string    `json:"title" gorm:"size:500;not null"`
	AttachmentID int       `json:"attachment_id" gorm:"not null"`
	RID          int       `json:"rid" gorm:"not null"`
	IP           string    `json:"ip" gorm:"size:255"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (CourseAttachmentDownloadLog) TableName() string {
	return "course_attachment_download_log"
}

type CourseDepartmentUser struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CourseID  int       `json:"course_id" gorm:"not null;index"`
	DepID     int       `json:"dep_id" gorm:"not null"`
	UserID    int       `json:"user_id" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (CourseDepartmentUser) TableName() string {
	return "course_department_user"
}

type UserCourseRecord struct {
	ID            int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        int        `json:"user_id" gorm:"not null;uniqueIndex:idx_user_course"`
	CourseID      int        `json:"course_id" gorm:"not null;uniqueIndex:idx_user_course;index"`
	HourCount     int        `json:"hour_count" gorm:"default:0"`
	FinishedCount int        `json:"finished_count" gorm:"default:0"`
	Progress      int        `json:"progress" gorm:"default:0"`
	IsFinished    int        `json:"is_finished" gorm:"default:0"`
	FinishedAt    *time.Time `json:"finished_at"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (UserCourseRecord) TableName() string {
	return "user_course_records"
}

type UserCourseHourRecord struct {
	ID               int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           int        `json:"user_id" gorm:"not null;uniqueIndex:idx_user_hour"`
	CourseID         int        `json:"course_id" gorm:"not null;index"`
	HourID           int        `json:"hour_id" gorm:"not null;uniqueIndex:idx_user_hour"`
	TotalDuration    int        `json:"total_duration" gorm:"default:0"`
	FinishedDuration int        `json:"finished_duration" gorm:"default:0"`
	RealDuration     int        `json:"real_duration" gorm:"default:0"`
	IsFinished       int        `json:"is_finished" gorm:"default:0"`
	FinishedAt       *time.Time `json:"finished_at"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (UserCourseHourRecord) TableName() string {
	return "user_course_hour_records"
}

type UserLearnDurationRecord struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id" gorm:"not null;index"`
	CourseID  int       `json:"course_id" gorm:"not null"`
	HourID    int       `json:"hour_id" gorm:"not null"`
	Duration  int       `json:"duration" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;index"`
}

func (UserLearnDurationRecord) TableName() string {
	return "user_learn_duration_records"
}

type UserLearnDurationStats struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int       `json:"user_id" gorm:"not null;uniqueIndex:idx_user_date"`
	Duration    int       `json:"duration" gorm:"not null"`
	CreatedDate time.Time `json:"created_date" gorm:"type:date;not null;uniqueIndex:idx_user_date"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (UserLearnDurationStats) TableName() string {
	return "user_learn_duration_stats"
}

type UserLatestLearn struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id" gorm:"not null;uniqueIndex:idx_user_course"`
	CourseID  int       `json:"course_id" gorm:"not null;uniqueIndex:idx_user_course"`
	HourID    int       `json:"hour_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (UserLatestLearn) TableName() string {
	return "user_latest_learn"
}
