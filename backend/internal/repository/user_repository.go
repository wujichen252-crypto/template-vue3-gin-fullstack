package repository

import (
	"errors"
	"template-vue3-gin-fullstack/backend/internal/model"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrUserAlreadyExists = errors.New("用户名或邮箱已存在")
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	List(page, pageSize int) ([]*model.User, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return containsAny(errStr, []string{
		"duplicate key",
		"UNIQUE constraint",
		"unique_violation",
		"23505",
	})
}

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.Select("id, username, email, avatar_url, status, created_at, updated_at").
		First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Select("id", "username", "email", "password_hash", "avatar_url", "status", "created_at", "updated_at", "deleted_at").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Select("id", "username", "email", "avatar_url", "status", "created_at", "updated_at", "deleted_at").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) List(page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Select("id", "username", "email", "avatar_url", "status", "created_at", "updated_at").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}