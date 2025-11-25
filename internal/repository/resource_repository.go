package repository

import (
	"github.com/eduflow/eduflow/internal/domain"
	"gorm.io/gorm"
)

type ResourceRepository interface {
	Create(resource *domain.Resource) error
	GetByID(id int) (*domain.Resource, error)
	Update(resource *domain.Resource) error
	Delete(id int) error
	List(page, size int, filters map[string]interface{}) ([]domain.Resource, int64, error)
}

type resourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) ResourceRepository {
	return &resourceRepository{db: db}
}

func (r *resourceRepository) Create(resource *domain.Resource) error {
	return r.db.Create(resource).Error
}

func (r *resourceRepository) GetByID(id int) (*domain.Resource, error) {
	var resource domain.Resource
	if err := r.db.First(&resource, id).Error; err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *resourceRepository) Update(resource *domain.Resource) error {
	return r.db.Save(resource).Error
}

func (r *resourceRepository) Delete(id int) error {
	return r.db.Delete(&domain.Resource{}, id).Error
}

func (r *resourceRepository) List(page, size int, filters map[string]interface{}) ([]domain.Resource, int64, error) {
	var resources []domain.Resource
	var total int64

	query := r.db.Model(&domain.Resource{})

	if resourceType, ok := filters["type"].(string); ok && resourceType != "" {
		query = query.Where("type = ?", resourceType)
	}

	if categoryID, ok := filters["category_id"].(int); ok && categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Order("id DESC").Offset(offset).Limit(size).Find(&resources).Error; err != nil {
		return nil, 0, err
	}

	return resources, total, nil
}
