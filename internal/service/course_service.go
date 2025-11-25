package service

import (
	"github.com/eduflow/eduflow/internal/domain"
	"github.com/eduflow/eduflow/internal/repository"
	"github.com/redis/go-redis/v9"
)

type CourseService interface {
	Create(course *domain.Course) error
	GetByID(id int) (*domain.Course, error)
	Update(course *domain.Course) error
	Delete(id int) error
	List(page, size int, filters map[string]interface{}) ([]domain.Course, int64, error)
}

type courseService struct {
	courseRepo repository.CourseRepository
	rdb        *redis.Client
}

func NewCourseService(courseRepo repository.CourseRepository, rdb *redis.Client) CourseService {
	return &courseService{
		courseRepo: courseRepo,
		rdb:        rdb,
	}
}

func (s *courseService) Create(course *domain.Course) error {
	return s.courseRepo.Create(course)
}

func (s *courseService) GetByID(id int) (*domain.Course, error) {
	return s.courseRepo.GetByID(id)
}

func (s *courseService) Update(course *domain.Course) error {
	return s.courseRepo.Update(course)
}

func (s *courseService) Delete(id int) error {
	return s.courseRepo.Delete(id)
}

func (s *courseService) List(page, size int, filters map[string]interface{}) ([]domain.Course, int64, error) {
	return s.courseRepo.List(page, size, filters)
}
