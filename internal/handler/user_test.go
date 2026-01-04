package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/olmits/budget-tracker-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// 1. Define the Mock User Repository
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Case", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		// Expect CreateUser to be called.
		// Note: We use mock.AnythingOfType for the user struct because the
		// handler generates the password hash dynamically, so we can't predict the exact struct.
		mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

		h := &UserHandler{Repo: mockRepo}
		r := gin.Default()
		r.POST("/register", h.Register)

		w := httptest.NewRecorder()
		jsonBody := []byte(`{"email": "test@test.com", "password": "123321"}`)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "User registered successfully")

		mockRepo.AssertExpectations(t)
	})
	t.Run("Invalid Input (Short Password)", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		// Repo should NOT be called

		h := &UserHandler{Repo: mockRepo}
		r := gin.Default()
		r.POST("/register", h.Register)

		w := httptest.NewRecorder()
		jsonBody := []byte(`{"email": "test@test.com", "password": "123"}`)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "'Password' failed on the 'min' tag")

		mockRepo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// IMPORTANT: Set JWT_SECRET for the test environment
	os.Setenv("JWT_SECRET", "super_secret_test_key")

	t.Run("Success Case", func(t *testing.T) {
		mockRepo := new(MockUserRepo)

		// 1. Prepare a user with a real hashed password
		password := "123321"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		dummyUser := &models.User{
			Email:        "test@test.com",
			PasswordHash: string(hashedPassword),
		}

		// 2. Mock finding the user
		mockRepo.On("GetUserByEmail", mock.Anything, "test@test.com").Return(dummyUser, nil)

		h := &UserHandler{Repo: mockRepo}
		r := gin.Default()
		r.POST("/login", h.Login)

		w := httptest.NewRecorder()
		jsonBody := []byte(`{"email": "test@test.com", "password": "123321"}`)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")

		mockRepo.AssertExpectations(t)
	})
	t.Run("Wrong Password", func(t *testing.T) {
		mockRepo := new(MockUserRepo)

		password := "123321"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		dummyUser := &models.User{
			Email:        "test@test.com",
			PasswordHash: string(hashedPassword),
		}

		mockRepo.On("GetUserByEmail", mock.Anything, "test@test.com").Return(dummyUser, nil)

		h := &UserHandler{Repo: mockRepo}
		r := gin.Default()
		r.POST("/login", h.Login)

		w := httptest.NewRecorder()
		jsonBody := []byte(`{"email": "test@test.com", "password": "321123"}`)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email or password")

		mockRepo.AssertExpectations(t)
	})
	t.Run("User Not FOund", func(t *testing.T) {
		mockRepo := new(MockUserRepo)

		mockRepo.On("GetUserByEmail", mock.Anything, "test@test.com").Return(nil, errors.New("user not found"))

		h := &UserHandler{Repo: mockRepo}
		r := gin.Default()
		r.POST("/login", h.Login)

		w := httptest.NewRecorder()
		jsonBody := []byte(`{"email": "test@test.com", "password": "123321"}`)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockRepo.AssertExpectations(t)
	})
}
