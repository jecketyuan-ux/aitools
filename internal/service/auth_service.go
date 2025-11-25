package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eduflow/eduflow/internal/domain"
	"github.com/eduflow/eduflow/internal/pkg/crypto"
	"github.com/eduflow/eduflow/internal/pkg/jwt"
	"github.com/eduflow/eduflow/internal/repository"
	"github.com/eduflow/eduflow/pkg/constants"
	"github.com/redis/go-redis/v9"
)

type AuthService interface {
	LoginUser(email, password string) (string, *domain.User, error)
	LoginAdmin(email, password string) (string, *domain.AdminUser, error)
	RegisterUser(email, password, name string) (*domain.User, error)
	Logout(jti string, expireTime time.Duration) error
	IsTokenBlacklisted(jti string) (bool, error)
}

type authService struct {
	userRepo      repository.UserRepository
	adminUserRepo repository.AdminUserRepository
	rdb           *redis.Client
	jwtManager    *jwt.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, adminUserRepo repository.AdminUserRepository, rdb *redis.Client, jwtManager *jwt.JWTManager) AuthService {
	return &authService{
		userRepo:      userRepo,
		adminUserRepo: adminUserRepo,
		rdb:           rdb,
		jwtManager:    jwtManager,
	}
}

func (s *authService) LoginUser(email, password string) (string, *domain.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, errors.New("user not found")
	}

	if user.IsLock == 1 {
		return "", nil, errors.New("user is locked")
	}

	if !crypto.VerifyPassword(password, user.Salt, user.Password) {
		return "", nil, errors.New("incorrect password")
	}

	token, jti, expiresAt, err := s.jwtManager.Generate(user.ID, user.Email, constants.RoleUser)
	if err != nil {
		return "", nil, err
	}

	now := time.Now()
	user.LoginAt = &now
	s.userRepo.Update(user)

	loginRecord := &domain.UserLoginRecord{
		UserID:    user.ID,
		JTI:       jti,
		ExpiredAt: expiresAt,
	}
	
	ctx := context.Background()
	s.rdb.Set(ctx, fmt.Sprintf("user_login:%d:%s", user.ID, jti), "1", time.Until(expiresAt))

	return token, user, nil
}

func (s *authService) LoginAdmin(email, password string) (string, *domain.AdminUser, error) {
	adminUser, err := s.adminUserRepo.GetByEmail(email)
	if err != nil {
		return "", nil, errors.New("admin user not found")
	}

	if adminUser.IsBanLogin == 1 {
		return "", nil, errors.New("admin user is banned")
	}

	if !crypto.VerifyPassword(password, adminUser.Salt, adminUser.Password) {
		return "", nil, errors.New("incorrect password")
	}

	token, _, _, err := s.jwtManager.Generate(adminUser.ID, adminUser.Email, constants.RoleAdmin)
	if err != nil {
		return "", nil, err
	}

	now := time.Now()
	adminUser.LoginAt = &now
	adminUser.LoginTimes++
	s.adminUserRepo.Update(adminUser)

	return token, adminUser, nil
}

func (s *authService) RegisterUser(email, password, name string) (*domain.User, error) {
	existingUser, _ := s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	salt, err := crypto.GenerateSalt(16)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := crypto.HashPassword(password, salt)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:         email,
		Name:          name,
		Password:      hashedPassword,
		Salt:          salt,
		IsActive:      1,
		IsLock:        0,
		IsVerify:      0,
		IsSetPassword: 1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Logout(jti string, expireTime time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf(constants.CacheKeyTokenBlacklist, jti)
	return s.rdb.Set(ctx, key, "1", expireTime).Err()
}

func (s *authService) IsTokenBlacklisted(jti string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf(constants.CacheKeyTokenBlacklist, jti)
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}
