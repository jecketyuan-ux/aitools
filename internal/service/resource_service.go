package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/eduflow/eduflow/internal/domain"
	"github.com/eduflow/eduflow/internal/pkg/storage"
	"github.com/eduflow/eduflow/internal/repository"
	"github.com/eduflow/eduflow/pkg/constants"
	"github.com/redis/go-redis/v9"
)

type ResourceService interface {
	Create(resource *domain.Resource) error
	GetByID(id int) (*domain.Resource, error)
	Delete(id int) error
	List(page, size int, filters map[string]interface{}) ([]domain.Resource, int64, error)
	UploadVideo(file *multipart.FileHeader, adminID int) (*domain.Resource, error)
	UploadImage(file *multipart.FileHeader, adminID int) (*domain.Resource, error)
}

type resourceService struct {
	resourceRepo repository.ResourceRepository
	storage      *storage.MinIOStorage
	rdb          *redis.Client
}

func NewResourceService(resourceRepo repository.ResourceRepository, storage *storage.MinIOStorage, rdb *redis.Client) ResourceService {
	return &resourceService{
		resourceRepo: resourceRepo,
		storage:      storage,
		rdb:          rdb,
	}
}

func (s *resourceService) Create(resource *domain.Resource) error {
	return s.resourceRepo.Create(resource)
}

func (s *resourceService) GetByID(id int) (*domain.Resource, error) {
	return s.resourceRepo.GetByID(id)
}

func (s *resourceService) Delete(id int) error {
	return s.resourceRepo.Delete(id)
}

func (s *resourceService) List(page, size int, filters map[string]interface{}) ([]domain.Resource, int64, error) {
	return s.resourceRepo.List(page, size, filters)
}

func (s *resourceService) UploadVideo(file *multipart.FileHeader, adminID int) (*domain.Resource, error) {
	if file.Size > constants.MaxVideoSize {
		return nil, fmt.Errorf("file size exceeds limit")
	}

	ext := filepath.Ext(file.Filename)
	objectName := s.storage.GenerateObjectName("videos", ext)

	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ctx := context.Background()
	url, err := s.storage.Upload(ctx, objectName, f, file.Size, "video/mp4")
	if err != nil {
		return nil, err
	}

	resource := &domain.Resource{
		AdminID:   adminID,
		Type:      constants.ResourceTypeVideo,
		URL:       url,
		Name:      file.Filename,
		Extension: ext,
		Size:      file.Size,
		Disk:      "minio",
		Path:      objectName,
	}

	if err := s.resourceRepo.Create(resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *resourceService) UploadImage(file *multipart.FileHeader, adminID int) (*domain.Resource, error) {
	if file.Size > constants.MaxImageSize {
		return nil, fmt.Errorf("file size exceeds limit")
	}

	ext := filepath.Ext(file.Filename)
	objectName := s.storage.GenerateObjectName("images", ext)

	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ctx := context.Background()
	url, err := s.storage.Upload(ctx, objectName, f, file.Size, "image/jpeg")
	if err != nil {
		return nil, err
	}

	resource := &domain.Resource{
		AdminID:   adminID,
		Type:      constants.ResourceTypeImage,
		URL:       url,
		Name:      file.Filename,
		Extension: ext,
		Size:      file.Size,
		Disk:      "minio",
		Path:      objectName,
	}

	if err := s.resourceRepo.Create(resource); err != nil {
		return nil, err
	}

	return resource, nil
}
