package service

import (
	"testing"
	"template-vue3-gin-fullstack/backend/internal/model"
	"template-vue3-gin-fullstack/backend/internal/repository"

	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

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
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func TestUserService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

	svc := NewUserService(mockRepo, nil)

	user, _, err := svc.Register("testuser", "password123", "test@example.com")
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Username mismatch: got %s, want testuser", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Email mismatch: got %s, want test@example.com", user.Email)
	}

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := &model.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Email:        "test@example.com",
		Status:       1,
	}

	mockRepo.On("GetByUsername", "testuser").Return(mockUser, nil)

	svc := NewUserService(mockRepo, nil)

	user, _, err := svc.Login("testuser", "password123")
	if err != nil {
		t.Errorf("Login failed: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Username mismatch: got %s, want testuser", user.Username)
	}

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("GetByUsername", "nonexistent").Return(nil, repository.ErrUserNotFound)

	svc := NewUserService(mockRepo, nil)

	_, _, err := svc.Login("nonexistent", "password123")
	if err == nil {
		t.Error("Login should fail for nonexistent user")
	}

	if err.Error() != "用户不存在" {
		t.Errorf("Error message mismatch: got %s, want 用户不存在", err.Error())
	}

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	mockUser := &model.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Email:        "test@example.com",
		Status:       1,
	}

	mockRepo.On("GetByUsername", "testuser").Return(mockUser, nil)

	svc := NewUserService(mockRepo, nil)

	_, _, err := svc.Login("testuser", "wrongpassword")
	if err == nil {
		t.Error("Login should fail for wrong password")
	}

	if err.Error() != "密码错误" {
		t.Errorf("Error message mismatch: got %s, want 密码错误", err.Error())
	}

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_DisabledUser(t *testing.T) {
	mockRepo := new(MockUserRepository)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := &model.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Email:        "test@example.com",
		Status:       0,
	}

	mockRepo.On("GetByUsername", "testuser").Return(mockUser, nil)

	svc := NewUserService(mockRepo, nil)

	_, _, err := svc.Login("testuser", "password123")
	if err == nil {
		t.Error("Login should fail for disabled user")
	}

	if err.Error() != "用户已被禁用" {
		t.Errorf("Error message mismatch: got %s, want 用户已被禁用", err.Error())
	}

	mockRepo.AssertExpectations(t)
}