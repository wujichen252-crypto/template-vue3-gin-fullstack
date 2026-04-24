// file: backend/internal/repository/user_repository_test.go
package repository

import (
	"errors"
	"testing"

	"template-vue3-gin-fullstack/backend/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// MockUserRepository 用于测试的 Mock 实现
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

func TestMockUserRepository_Create(t *testing.T) {
	mockRepo := new(MockUserRepository)

	tests := []struct {
		name    string
		user    *model.User
		mockErr error
		wantErr bool
	}{
		{
			name: "正常创建用户",
			user: &model.User{
				Username:     "testuser",
				PasswordHash: "hashedpassword",
				Email:        "test@example.com",
				Status:       1,
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "用户名已存在",
			user: &model.User{
				Username:     "existing",
				PasswordHash: "hashedpassword",
				Email:        "existing@example.com",
				Status:       1,
			},
			mockErr: ErrUserAlreadyExists,
			wantErr: true,
		},
		{
			name: "数据库连接错误",
			user: &model.User{
				Username:     "testuser2",
				PasswordHash: "hashedpassword",
				Email:        "test2@example.com",
				Status:       1,
			},
			mockErr: errors.New("connection refused"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("Create", tt.user).Return(tt.mockErr).Once()

			err := mockRepo.Create(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.mockErr, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMockUserRepository_GetByID(t *testing.T) {
	mockRepo := new(MockUserRepository)

	now := time.Now()
	testUser := &model.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Email:        "test@example.com",
		Status:       1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tests := []struct {
		name     string
		id       uint
		mockUser *model.User
		mockErr  error
		wantErr  bool
	}{
		{
			name:     "正常获取用户",
			id:       1,
			mockUser: testUser,
			mockErr:  nil,
			wantErr:  false,
		},
		{
			name:     "用户不存在",
			id:       9999,
			mockUser: nil,
			mockErr:  ErrUserNotFound,
			wantErr:  true,
		},
		{
			name:     "ID为0",
			id:       0,
			mockUser: nil,
			mockErr:  ErrUserNotFound,
			wantErr:  true,
		},
		{
			name:     "数据库错误",
			id:       1,
			mockUser: nil,
			mockErr:  errors.New("database error"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetByID", tt.id).Return(tt.mockUser, tt.mockErr).Once()

			got, err := mockRepo.GetByID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.mockErr, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.mockUser.Username, got.Username)
				assert.Equal(t, tt.mockUser.Email, got.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMockUserRepository_GetByUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)

	testUser := &model.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Email:        "test@example.com",
		Status:       1,
	}

	tests := []struct {
		name     string
		username string
		mockUser *model.User
		mockErr  error
		wantErr  bool
	}{
		{
			name:     "正常获取用户",
			username: "testuser",
			mockUser: testUser,
			mockErr:  nil,
			wantErr:  false,
		},
		{
			name:     "用户不存在",
			username: "nonexistent",
			mockUser: nil,
			mockErr:  ErrUserNotFound,
			wantErr:  true,
		},
		{
			name:     "空用户名",
			username: "",
			mockUser: nil,
			mockErr:  ErrUserNotFound,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetByUsername", tt.username).Return(tt.mockUser, tt.mockErr).Once()

			got, err := mockRepo.GetByUsername(tt.username)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.mockUser.Username, got.Username)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMockUserRepository_GetByEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)

	testUser := &model.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Email:        "test@example.com",
		Status:       1,
	}

	tests := []struct {
		name     string
		email    string
		mockUser *model.User
		mockErr  error
		wantErr  bool
	}{
		{
			name:     "正常获取用户",
			email:    "test@example.com",
			mockUser: testUser,
			mockErr:  nil,
			wantErr:  false,
		},
		{
			name:     "邮箱不存在",
			email:    "nonexistent@example.com",
			mockUser: nil,
			mockErr:  ErrUserNotFound,
			wantErr:  true,
		},
		{
			name:     "空邮箱",
			email:    "",
			mockUser: nil,
			mockErr:  ErrUserNotFound,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetByEmail", tt.email).Return(tt.mockUser, tt.mockErr).Once()

			got, err := mockRepo.GetByEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.mockUser.Email, got.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMockUserRepository_Update(t *testing.T) {
	mockRepo := new(MockUserRepository)

	tests := []struct {
		name    string
		user    *model.User
		mockErr error
		wantErr bool
	}{
		{
			name: "正常更新用户",
			user: &model.User{
				ID:           1,
				Username:     "updated",
				PasswordHash: "newhash",
				Email:        "updated@example.com",
				Status:       1,
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "更新不存在用户",
			user: &model.User{
				ID:           9999,
				Username:     "nonexistent",
				PasswordHash: "hash",
				Email:        "non@existent.com",
			},
			mockErr: ErrUserNotFound,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("Update", tt.user).Return(tt.mockErr).Once()

			err := mockRepo.Update(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMockUserRepository_Delete(t *testing.T) {
	mockRepo := new(MockUserRepository)

	tests := []struct {
		name    string
		id      uint
		mockErr error
		wantErr bool
	}{
		{
			name:    "正常删除用户",
			id:      1,
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "删除不存在用户",
			id:      9999,
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "删除ID为0",
			id:      0,
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "数据库错误",
			id:      1,
			mockErr: errors.New("database error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("Delete", tt.id).Return(tt.mockErr).Once()

			err := mockRepo.Delete(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMockUserRepository_List(t *testing.T) {
	mockRepo := new(MockUserRepository)

	now := time.Now()
	testUsers := []*model.User{
		{ID: 1, Username: "user1", PasswordHash: "hash1", Email: "user1@test.com", Status: 1, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Username: "user2", PasswordHash: "hash2", Email: "user2@test.com", Status: 1, CreatedAt: now, UpdatedAt: now},
		{ID: 3, Username: "user3", PasswordHash: "hash3", Email: "user3@test.com", Status: 1, CreatedAt: now, UpdatedAt: now},
	}

	tests := []struct {
		name         string
		page         int
		pageSize     int
		mockUsers    []*model.User
		mockTotal    int64
		mockErr      error
		wantCount    int
		wantTotal    int64
		wantErr      bool
	}{
		{
			name:      "正常分页-第一页",
			page:      1,
			pageSize:  2,
			mockUsers: testUsers[:2],
			mockTotal: 3,
			mockErr:   nil,
			wantCount: 2,
			wantTotal: 3,
			wantErr:   false,
		},
		{
			name:      "正常分页-第二页",
			page:      2,
			pageSize:  2,
			mockUsers: testUsers[2:3],
			mockTotal: 3,
			mockErr:   nil,
			wantCount: 1,
			wantTotal: 3,
			wantErr:   false,
		},
		{
			name:      "空结果",
			page:      1,
			pageSize:  10,
			mockUsers: []*model.User{},
			mockTotal: 0,
			mockErr:   nil,
			wantCount: 0,
			wantTotal: 0,
			wantErr:   false,
		},
		{
			name:      "数据库错误",
			page:      1,
			pageSize:  10,
			mockUsers: nil,
			mockTotal: 0,
			mockErr:   errors.New("database error"),
			wantCount: 0,
			wantTotal: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("List", tt.page, tt.pageSize).Return(tt.mockUsers, tt.mockTotal, tt.mockErr).Once()

			got, total, err := mockRepo.List(tt.page, tt.pageSize)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.Zero(t, total)
			} else {
				assert.NoError(t, err)
				assert.Len(t, got, tt.wantCount)
				assert.Equal(t, tt.wantTotal, total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestIsDuplicateKeyError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil错误",
			err:  nil,
			want: false,
		},
		{
			name: "duplicate key错误",
			err:  errors.New("ERROR: duplicate key value violates unique constraint \"users_username_key\""),
			want: true,
		},
		{
			name: "UNIQUE constraint错误",
			err:  errors.New("UNIQUE constraint failed: users.username"),
			want: true,
		},
		{
			name: "其他错误",
			err:  errors.New("connection refused"),
			want: false,
		},
		{
			name: "空错误",
			err:  errors.New(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDuplicateKeyError(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
