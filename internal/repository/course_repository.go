package repository

import (
	"github.com/eduflow/eduflow/internal/domain"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *domain.Course) error
	GetByID(id int) (*domain.Course, error)
	Update(course *domain.Course) error
	Delete(id int) error
	List(page, size int, filters map[string]interface{}) ([]domain.Course, int64, error)
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(course *domain.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) GetByID(id int) (*domain.Course, error) {
	var course domain.Course
	if err := r.db.Where("deleted_at IS NULL").First(&course, id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) Update(course *domain.Course) error {
	return r.db.Save(course).Error
}

func (r *courseRepository) Delete(id int) error {
	return r.db.Model(&domain.Course{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *courseRepository) List(page, size int, filters map[string]interface{}) ([]domain.Course, int64, error) {
	var courses []domain.Course
	var total int64

	query := r.db.Model(&domain.Course{}).Where("deleted_at IS NULL")

	if title, ok := filters["title"].(string); ok && title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	if isShow, ok := filters["is_show"].(int); ok {
		query = query.Where("is_show = ?", isShow)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Order("sort_at DESC, id DESC").Offset(offset).Limit(size).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}
