package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/playedu/playedu-go/internal/config"
	"github.com/playedu/playedu-go/internal/domain"
	"github.com/playedu/playedu-go/internal/handler/backend"
	"github.com/playedu/playedu-go/internal/handler/frontend"
	"github.com/playedu/playedu-go/internal/middleware"
	"github.com/playedu/playedu-go/internal/pkg/jwt"
	"github.com/playedu/playedu-go/internal/pkg/storage"
	"github.com/playedu/playedu-go/internal/repository"
	"github.com/playedu/playedu-go/internal/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	rdb := initRedis(cfg)

	minioStorage, err := initMinIO(cfg)
	if err != nil {
		log.Fatalf("Failed to init MinIO: %v", err)
	}

	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireTime)

	repos := initRepositories(db)
	services := initServices(repos, rdb, minioStorage, jwtManager, cfg)
	handlers := initHandlers(services)

	router := setupRouter(cfg, handlers, jwtManager, rdb, services)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("Server starting on port %d...", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	if err := autoMigrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Department{},
		&domain.UserDepartment{},
		&domain.UserLoginRecord{},
		&domain.UserUploadImageLog{},
		&domain.AdminUser{},
		&domain.AdminRole{},
		&domain.AdminPermission{},
		&domain.AdminRolePermission{},
		&domain.AdminUserRole{},
		&domain.AdminLog{},
		&domain.Course{},
		&domain.CourseChapter{},
		&domain.CourseHour{},
		&domain.CourseCategory{},
		&domain.CourseAttachment{},
		&domain.CourseAttachmentDownloadLog{},
		&domain.CourseDepartmentUser{},
		&domain.UserCourseRecord{},
		&domain.UserCourseHourRecord{},
		&domain.UserLearnDurationRecord{},
		&domain.UserLearnDurationStats{},
		&domain.UserLatestLearn{},
		&domain.Resource{},
		&domain.ResourceCategory{},
		&domain.ResourceVideo{},
		&domain.Category{},
		&domain.AppConfig{},
		&domain.LdapUser{},
		&domain.LdapDepartment{},
		&domain.LdapSyncRecord{},
	)
}

func initRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
}

func initMinIO(cfg *config.Config) (*storage.MinIOStorage, error) {
	return storage.NewMinIOStorage(
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKeyID,
		cfg.MinIO.SecretAccessKey,
		cfg.MinIO.BucketName,
		cfg.MinIO.UseSSL,
	)
}

type Repositories struct {
	User            repository.UserRepository
	Department      repository.DepartmentRepository
	AdminUser       repository.AdminUserRepository
	Course          repository.CourseRepository
	Resource        repository.ResourceRepository
	Category        repository.CategoryRepository
}

func initRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:            repository.NewUserRepository(db),
		Department:      repository.NewDepartmentRepository(db),
		AdminUser:       repository.NewAdminUserRepository(db),
		Course:          repository.NewCourseRepository(db),
		Resource:        repository.NewResourceRepository(db),
		Category:        repository.NewCategoryRepository(db),
	}
}

type Services struct {
	Auth     service.AuthService
	User     service.UserService
	Course   service.CourseService
	Resource service.ResourceService
}

func initServices(repos *Repositories, rdb *redis.Client, storage *storage.MinIOStorage, jwtManager *jwt.JWTManager, cfg *config.Config) *Services {
	return &Services{
		Auth:     service.NewAuthService(repos.User, repos.AdminUser, rdb, jwtManager),
		User:     service.NewUserService(repos.User, repos.Department, rdb),
		Course:   service.NewCourseService(repos.Course, rdb),
		Resource: service.NewResourceService(repos.Resource, storage, rdb),
	}
}

type Handlers struct {
	BackendAuth     *backend.AuthHandler
	BackendUser     *backend.UserHandler
	BackendCourse   *backend.CourseHandler
	BackendResource *backend.ResourceHandler
	FrontendAuth    *frontend.AuthHandler
	FrontendCourse  *frontend.CourseHandler
}

func initHandlers(services *Services) *Handlers {
	return &Handlers{
		BackendAuth:     backend.NewAuthHandler(services.Auth),
		BackendUser:     backend.NewUserHandler(services.User),
		BackendCourse:   backend.NewCourseHandler(services.Course),
		BackendResource: backend.NewResourceHandler(services.Resource),
		FrontendAuth:    frontend.NewAuthHandler(services.Auth),
		FrontendCourse:  frontend.NewCourseHandler(services.Course),
	}
}

func setupRouter(cfg *config.Config, handlers *Handlers, jwtManager *jwt.JWTManager, rdb *redis.Client, services *Services) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	backendAuth := middleware.NewAuthMiddleware(jwtManager, rdb, "admin")
	frontendAuth := middleware.NewAuthMiddleware(jwtManager, rdb, "user")
	rateLimiter := middleware.NewRateLimiter(rdb, cfg.RateLimit.Duration, cfg.RateLimit.Limit)

	backendV1 := router.Group("/backend/v1")
	backendV1.Use(rateLimiter.Limit())
	{
		backendV1.POST("/auth/login", handlers.BackendAuth.Login)
		
		authGroup := backendV1.Group("")
		authGroup.Use(backendAuth.Authenticate())
		{
			authGroup.POST("/auth/logout", handlers.BackendAuth.Logout)
			authGroup.GET("/auth/detail", handlers.BackendAuth.GetDetail)

			authGroup.GET("/user", handlers.BackendUser.List)
			authGroup.POST("/user", handlers.BackendUser.Create)
			authGroup.GET("/user/:id", handlers.BackendUser.GetByID)
			authGroup.PUT("/user/:id", handlers.BackendUser.Update)
			authGroup.DELETE("/user/:id", handlers.BackendUser.Delete)

			authGroup.GET("/course", handlers.BackendCourse.List)
			authGroup.POST("/course", handlers.BackendCourse.Create)
			authGroup.GET("/course/:id", handlers.BackendCourse.GetByID)
			authGroup.PUT("/course/:id", handlers.BackendCourse.Update)
			authGroup.DELETE("/course/:id", handlers.BackendCourse.Delete)

			authGroup.GET("/resource", handlers.BackendResource.List)
			authGroup.POST("/resource/video/upload", handlers.BackendResource.UploadVideo)
			authGroup.POST("/resource/image/upload", handlers.BackendResource.UploadImage)
			authGroup.DELETE("/resource/:id", handlers.BackendResource.Delete)
		}
	}

	apiV1 := router.Group("/api/v1")
	apiV1.Use(rateLimiter.Limit())
	{
		apiV1.POST("/auth/login", handlers.FrontendAuth.Login)
		apiV1.POST("/auth/register", handlers.FrontendAuth.Register)

		authGroup := apiV1.Group("")
		authGroup.Use(frontendAuth.Authenticate())
		{
			authGroup.POST("/auth/logout", handlers.FrontendAuth.Logout)
			authGroup.GET("/auth/detail", handlers.FrontendAuth.GetDetail)

			authGroup.GET("/courses", handlers.FrontendCourse.List)
			authGroup.GET("/course/:id", handlers.FrontendCourse.GetByID)
		}
	}

	return router
}
