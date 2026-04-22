package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"template-vue3-gin-fullstack/backend/internal/model"
	"template-vue3-gin-fullstack/backend/internal/repository"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(username, password, email string) (*model.User, string, error)
	Login(username, password string) (*model.User, string, error)
	GetUserInfo(id uint) (*model.User, error)
	RefreshToken(userID uint) error
	Logout(token string, exp time.Duration) error
}

type userService struct {
	repo repository.UserRepository
	rdb  *redis.Client
}

func NewUserService(repo repository.UserRepository, rdb *redis.Client) UserService {
	return &userService{repo: repo, rdb: rdb}
}

func (s *userService) Register(username, password, email string) (*model.User, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(hash),
		Email:        email,
		AvatarURL:    "",
		Status:       1,
	}

	if err := s.repo.Create(user); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, "", errors.New("用户名或邮箱已被注册")
		}
		return nil, "", err
	}

	return user, "", nil
}

func (s *userService) Login(username, password string) (*model.User, string, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, "", errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("密码错误")
	}

	if user.Status != 1 {
		return nil, "", errors.New("用户已被禁用")
	}

	return user, "", nil
}

func (s *userService) GetUserInfo(id uint) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) RefreshToken(userID uint) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	if user.Status != 1 {
		return errors.New("用户已被禁用")
	}

	return nil
}

func (s *userService) Logout(token string, exp time.Duration) error {
	if s.rdb == nil {
		return nil
	}
	ctx := context.Background()
	key := fmt.Sprintf("jwt_blacklist:%s", token)
	return s.rdb.Set(ctx, key, "1", exp).Err()
}