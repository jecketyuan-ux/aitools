package repository

import (
	"github.com/playedu/playedu-go/internal/domain"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *domain.Category) error
	GetByID(id int) (*domain.Category, error)
	Update(category *domain.Category) error
	Delete(id int) error
	List() ([]domain.Category, error)
	GetChildren(parentID int) ([]domain.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id int) (*domain.Category, error) {
	var category domain.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id int) error {
	return r.db.Delete(&domain.Category{}, id).Error
}

func (r *categoryRepository) List() ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.db.Order("sort ASC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) GetChildren(parentID int) ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.db.Where("parent_id = ?", parentID).Order("sort ASC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
