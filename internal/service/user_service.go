package service

import (
	"github.com/playedu/playedu-go/internal/domain"
	"github.com/playedu/playedu-go/internal/repository"
	"github.com/redis/go-redis/v9"
)

type UserService interface {
	Create(user *domain.User) error
	GetByID(id int) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id int) error
	List(page, size int, filters map[string]interface{}) ([]domain.User, int64, error)
}

type userService struct {
	userRepo repository.UserRepository
	deptRepo repository.DepartmentRepository
	rdb      *redis.Client
}

func NewUserService(userRepo repository.UserRepository, deptRepo repository.DepartmentRepository, rdb *redis.Client) UserService {
	return &userService{
		userRepo: userRepo,
		deptRepo: deptRepo,
		rdb:      rdb,
	}
}

func (s *userService) Create(user *domain.User) error {
	return s.userRepo.Create(user)
}

func (s *userService) GetByID(id int) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) Update(user *domain.User) error {
	return s.userRepo.Update(user)
}

func (s *userService) Delete(id int) error {
	return s.userRepo.Delete(id)
}

func (s *userService) List(page, size int, filters map[string]interface{}) ([]domain.User, int64, error) {
	return s.userRepo.List(page, size, filters)
}
