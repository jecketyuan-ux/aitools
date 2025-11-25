package repository

import (
	"github.com/playedu/playedu-go/internal/domain"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	Create(department *domain.Department) error
	GetByID(id int) (*domain.Department, error)
	Update(department *domain.Department) error
	Delete(id int) error
	List() ([]domain.Department, error)
	GetChildren(parentID int) ([]domain.Department, error)
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) Create(department *domain.Department) error {
	return r.db.Create(department).Error
}

func (r *departmentRepository) GetByID(id int) (*domain.Department, error) {
	var department domain.Department
	if err := r.db.First(&department, id).Error; err != nil {
		return nil, err
	}
	return &department, nil
}

func (r *departmentRepository) Update(department *domain.Department) error {
	return r.db.Save(department).Error
}

func (r *departmentRepository) Delete(id int) error {
	return r.db.Delete(&domain.Department{}, id).Error
}

func (r *departmentRepository) List() ([]domain.Department, error) {
	var departments []domain.Department
	if err := r.db.Order("sort ASC, id ASC").Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}

func (r *departmentRepository) GetChildren(parentID int) ([]domain.Department, error) {
	var departments []domain.Department
	if err := r.db.Where("parent_id = ?", parentID).Order("sort ASC, id ASC").Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}
