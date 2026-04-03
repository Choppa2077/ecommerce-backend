package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/PrimeraAizen/e-comm/config"
	"github.com/PrimeraAizen/e-comm/internal/domain"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// testConfig returns a config suitable for unit tests
func testConfig() *config.Config {
	return &config.Config{
		JWT: config.JWT{
			Secret:               "test-secret-key-for-unit-tests",
			AccessTokenDuration:  "15m",
			RefreshTokenDuration: "168h",
		},
	}
}

func newTestAuthService(userRepo *MockUserRepository) AuthService {
	svc, err := NewAuthService(userRepo, testConfig())
	if err != nil {
		panic("failed to create auth service: " + err.Error())
	}
	return svc
}

// --- Register tests ---

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	user := &domain.User{Email: "newuser@example.com", PasswordHash: "password123"}

	mockRepo.On("GetByEmail", ctx, "newuser@example.com").Return(nil, domain.ErrNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	token, err := svc.Register(ctx, user)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.AccessToken)
	assert.NotEmpty(t, token.RefreshToken)
	assert.Equal(t, "Bearer", token.TokenType)
	mockRepo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	existingUser := &domain.User{ID: 1, Email: "existing@example.com", Status: "active"}
	mockRepo.On("GetByEmail", ctx, "existing@example.com").Return(existingUser, nil)

	user := &domain.User{Email: "existing@example.com", PasswordHash: "password123"}
	token, err := svc.Register(ctx, user)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrAlreadyExists, err)
	assert.Nil(t, token)
	mockRepo.AssertExpectations(t)
}

// --- Login tests ---

func TestLogin_ValidCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	existingUser := &domain.User{
		ID:           1,
		Email:        "user@example.com",
		PasswordHash: string(hash),
		Status:       "active",
		CreatedAt:    time.Now(),
	}

	mockRepo.On("GetByEmail", ctx, "user@example.com").Return(existingUser, nil)
	mockRepo.On("UpdateLastLogin", ctx, 1).Return(nil)

	req := &domain.LoginRequest{Email: "user@example.com", Password: "password123"}
	token, err := svc.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.AccessToken)
	assert.NotEmpty(t, token.RefreshToken)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.MinCost)
	existingUser := &domain.User{
		ID:           1,
		Email:        "user@example.com",
		PasswordHash: string(hash),
		Status:       "active",
	}

	mockRepo.On("GetByEmail", ctx, "user@example.com").Return(existingUser, nil)

	req := &domain.LoginRequest{Email: "user@example.com", Password: "wrongpassword"}
	token, err := svc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidCredentials, err)
	assert.Nil(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByEmail", ctx, "notexist@example.com").Return(nil, domain.ErrNotFound)

	req := &domain.LoginRequest{Email: "notexist@example.com", Password: "password123"}
	token, err := svc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidCredentials, err)
	assert.Nil(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InactiveUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	inactiveUser := &domain.User{
		ID:           2,
		Email:        "inactive@example.com",
		PasswordHash: string(hash),
		Status:       "suspended",
	}

	mockRepo.On("GetByEmail", ctx, "inactive@example.com").Return(inactiveUser, nil)

	req := &domain.LoginRequest{Email: "inactive@example.com", Password: "password123"}
	token, err := svc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserInactive, err)
	assert.Nil(t, token)
	mockRepo.AssertExpectations(t)
}

// --- ValidateToken tests ---

func TestValidateToken_Valid(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)
	ctx := context.Background()

	// First register to get a real token
	mockRepo.On("GetByEmail", ctx, "token@example.com").Return(nil, domain.ErrNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	user := &domain.User{ID: 1, Email: "token@example.com"}
	token, err := svc.Register(ctx, user)
	assert.NoError(t, err)

	// Now validate the token
	claims, err := svc.ValidateToken(token.AccessToken)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "token@example.com", claims.Email)
}

func TestValidateToken_Invalid(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)

	claims, err := svc.ValidateToken("this.is.not.a.valid.token")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := newTestAuthService(mockRepo)

	// Token signed with a different secret
	wrongSecretToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjk5OTk5OTk5OTksInVzZXJfaWQiOiIxIn0.wrongsignature"

	claims, err := svc.ValidateToken(wrongSecretToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}
