package repository

import (
	"github.com/playedu/playedu-go/internal/domain"
	"gorm.io/gorm"
)

type AdminUserRepository interface {
	Create(adminUser *domain.AdminUser) error
	GetByID(id int) (*domain.AdminUser, error)
	GetByEmail(email string) (*domain.AdminUser, error)
	Update(adminUser *domain.AdminUser) error
	Delete(id int) error
	List(page, size int) ([]domain.AdminUser, int64, error)
}

type adminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) AdminUserRepository {
	return &adminUserRepository{db: db}
}

func (r *adminUserRepository) Create(adminUser *domain.AdminUser) error {
	return r.db.Create(adminUser).Error
}

func (r *adminUserRepository) GetByID(id int) (*domain.AdminUser, error) {
	var adminUser domain.AdminUser
	if err := r.db.First(&adminUser, id).Error; err != nil {
		return nil, err
	}
	return &adminUser, nil
}

func (r *adminUserRepository) GetByEmail(email string) (*domain.AdminUser, error) {
	var adminUser domain.AdminUser
	if err := r.db.Where("email = ?", email).First(&adminUser).Error; err != nil {
		return nil, err
	}
	return &adminUser, nil
}

func (r *adminUserRepository) Update(adminUser *domain.AdminUser) error {
	return r.db.Save(adminUser).Error
}

func (r *adminUserRepository) Delete(id int) error {
	return r.db.Delete(&domain.AdminUser{}, id).Error
}

func (r *adminUserRepository) List(page, size int) ([]domain.AdminUser, int64, error) {
	var adminUsers []domain.AdminUser
	var total int64

	query := r.db.Model(&domain.AdminUser{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Find(&adminUsers).Error; err != nil {
		return nil, 0, err
	}

	return adminUsers, total, nil
}
