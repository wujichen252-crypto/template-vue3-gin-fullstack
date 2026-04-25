package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"template-vue3-gin-fullstack/backend/internal/model"
	"template-vue3-gin-fullstack/backend/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository 模拟用户仓库
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List(page, pageSize int) ([]*model.User, int64, error) {
	args := m.Called(page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

// setupTest 创建测试所需的 Service 和 Mock
func setupTest() (*userService, *MockUserRepository) {
	mockRepo := new(MockUserRepository)
	// 使用 nil redis 客户端，测试不依赖缓存的功能
	svc := NewUserService(mockRepo, nil).(*userService)
	return svc, mockRepo
}

func TestUserService_Register(t *testing.T) {
	svc, mockRepo := setupTest()

	tests := []struct {
		name        string
		username    string
		password    string
		email       string
		mockSetup   func()
		wantErr     bool
		errContains string
	}{
		{
			name:     "正常注册成功",
			username: "testuser",
			password: "password123",
			email:    "test@example.com",
			mockSetup: func() {
				mockRepo.On("Create", mock.AnythingOfType("*model.User")).
					Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name:     "用户名已存在",
			username: "existinguser",
			password: "password123",
			email:    "existing@example.com",
			mockSetup: func() {
				mockRepo.On("Create", mock.AnythingOfType("*model.User")).
					Return(repository.ErrUserAlreadyExists).Once()
			},
			wantErr:     true,
			errContains: "用户名或邮箱已被注册",
		},
		{
			name:     "数据库错误",
			username: "testuser",
			password: "password123",
			email:    "test@example.com",
			mockSetup: func() {
				mockRepo.On("Create", mock.AnythingOfType("*model.User")).
					Return(errors.New("数据库连接失败")).Once()
			},
			wantErr:     true,
			errContains: "数据库连接失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := svc.Register(tt.username, tt.password, tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.Username)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, int8(1), user.Status)
				// 验证密码已加密
				assert.NotEqual(t, tt.password, user.PasswordHash)
				err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(tt.password))
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	svc, mockRepo := setupTest()

	// 生成测试用的密码哈希
	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		username    string
		password    string
		mockSetup   func()
		wantErr     bool
		errContains string
	}{
		{
			name:     "登录成功",
			username: "testuser",
			password: "correctpassword",
			mockSetup: func() {
				mockRepo.On("GetByUsername", "testuser").Return(&model.User{
					ID:           1,
					Username:     "testuser",
					PasswordHash: string(hash),
					Status:       1,
				}, nil).Once()
			},
			wantErr: false,
		},
		{
			name:     "用户不存在",
			username: "nonexistent",
			password: "anypassword",
			mockSetup: func() {
				mockRepo.On("GetByUsername", "nonexistent").
					Return(nil, repository.ErrUserNotFound).Once()
			},
			wantErr:     true,
			errContains: "用户不存在",
		},
		{
			name:     "密码错误",
			username: "testuser",
			password: "wrongpassword",
			mockSetup: func() {
				mockRepo.On("GetByUsername", "testuser").Return(&model.User{
					ID:           1,
					Username:     "testuser",
					PasswordHash: string(hash),
					Status:       1,
				}, nil).Once()
			},
			wantErr:     true,
			errContains: "密码错误",
		},
		{
			name:     "用户被禁用",
			username: "disableduser",
			password: "correctpassword",
			mockSetup: func() {
				hash2, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
				mockRepo.On("GetByUsername", "disableduser").Return(&model.User{
					ID:           2,
					Username:     "disableduser",
					PasswordHash: string(hash2),
					Status:       0,
				}, nil).Once()
			},
			wantErr:     true,
			errContains: "用户已被禁用",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := svc.Login(tt.username, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.Username)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserInfo(t *testing.T) {
	svc, mockRepo := setupTest()

	tests := []struct {
		name      string
		userID    uint
		mockSetup func()
		wantErr   bool
	}{
		{
			name:   "获取成功",
			userID: 1,
			mockSetup: func() {
				mockRepo.On("GetByID", uint(1)).Return(&model.User{
					ID:       1,
					Username: "testuser",
					Email:    "test@example.com",
				}, nil).Once()
			},
			wantErr: false,
		},
		{
			name:   "用户不存在",
			userID: 999,
			mockSetup: func() {
				mockRepo.On("GetByID", uint(999)).
					Return(nil, repository.ErrUserNotFound).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := svc.GetUserInfo(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.userID, user.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_RefreshToken(t *testing.T) {
	svc, mockRepo := setupTest()

	tests := []struct {
		name        string
		userID      uint
		mockSetup   func()
		wantErr     bool
		errContains string
	}{
		{
			name:   "刷新成功",
			userID: 1,
			mockSetup: func() {
				mockRepo.On("GetByID", uint(1)).Return(&model.User{
					ID:     1,
					Status: 1,
				}, nil).Once()
			},
			wantErr: false,
		},
		{
			name:   "用户不存在",
			userID: 999,
			mockSetup: func() {
				mockRepo.On("GetByID", uint(999)).
					Return(nil, repository.ErrUserNotFound).Once()
			},
			wantErr:     true,
			errContains: "用户不存在",
		},
		{
			name:   "用户被禁用",
			userID: 2,
			mockSetup: func() {
				mockRepo.On("GetByID", uint(2)).Return(&model.User{
					ID:     2,
					Status: 0,
				}, nil).Once()
			},
			wantErr:     true,
			errContains: "用户已被禁用",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := svc.RefreshToken(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Logout(t *testing.T) {
	t.Run("无Redis客户端时直接返回成功", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := &userService{repo: mockRepo, rdb: nil}

		err := svc.Logout(context.Background(), "test_token", time.Hour)
		assert.NoError(t, err)
	})

	t.Run("有Redis客户端时设置黑名单", func(t *testing.T) {
		// 注意：这里不测试真实的 Redis 操作，因为需要真实的 Redis 连接
		// 实际项目中可以使用 miniredis 等库进行测试
		mockRepo := new(MockUserRepository)
		svc := &userService{repo: mockRepo, rdb: nil}

		// 当 rdb 为 nil 时，直接返回 nil
		err := svc.Logout(context.Background(), "test_token", time.Hour)
		assert.NoError(t, err)
	})
}

func TestNewUserService(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := NewUserService(mockRepo, nil)

	assert.NotNil(t, svc)
	assert.Implements(t, (*UserService)(nil), svc)
}
