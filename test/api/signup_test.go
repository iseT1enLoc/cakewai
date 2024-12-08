package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"cakewai/cakewai.com/api/handlers"
	appconfig "cakewai/cakewai.com/component/appcfg"
	apperror "cakewai/cakewai.com/component/apperr"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockSignupUseCase is a mock implementation of the SignupUseCase interface
type MockSignupUseCase struct {
	mock.Mock
}

func (m *MockSignupUseCase) SignUp(ctx context.Context, request domain.SignupRequest, env *appconfig.Env) (string, string, string, error) {
	args := m.Called(ctx, request, env)
	return args.String(0), args.String(1), args.String(2), args.Error(3)
}

// MockCartUseCase is a mock implementation of the CartUsecase interface
type MockCartUseCase struct {
	mock.Mock
}

// CreateCartByUserId mocks the CreateCartByUserId method
func (m *MockCartUseCase) CreateCartByUserId(ctx context.Context, userID primitive.ObjectID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// GetAllItemsInCartByUserID mocks the GetAllItemsInCartByUserID method
func (m *MockCartUseCase) GetAllItemsInCartByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.CartItem, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.CartItem), args.Error(1)
}

// GetCartByUserID mocks the GetCartByUserID method
func (m *MockCartUseCase) GetCartByUserID(ctx context.Context, userID primitive.ObjectID) (*domain.Cart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Cart), args.Error(1)
	}
	return nil, args.Error(1)
}

// RemoveItemFromCart mocks the RemoveItemFromCart method
func (m *MockCartUseCase) RemoveItemFromCart(ctx context.Context, cartID primitive.ObjectID, productID primitive.ObjectID, variant string) error {
	args := m.Called(ctx, cartID, productID, variant)
	return args.Error(0)
}

// AddCartItemIntoCart mocks the AddCartItemIntoCart method
func (m *MockCartUseCase) AddCartItemIntoCart(ctx context.Context, cartID primitive.ObjectID, item domain.CartItem) (*primitive.ObjectID, error) {
	args := m.Called(ctx, cartID, item)
	if args.Get(0) != nil {
		return args.Get(0).(*primitive.ObjectID), args.Error(1)
	}
	return nil, args.Error(1)
}

// UpdateCartItemByID mocks the UpdateCartItemByID method
func (m *MockCartUseCase) UpdateCartItemByID(ctx context.Context, cartID primitive.ObjectID, updatedItem domain.CartItem) (*domain.CartItem, error) {
	args := m.Called(ctx, cartID, updatedItem)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.CartItem), args.Error(1)
	}
	return nil, args.Error(1)
}
func TestSignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test cases
	tests := []struct {
		name         string
		requestBody  interface{}
		setupMocks   func(mockSignupUseCase *MockSignupUseCase, mockCartUseCase *MockCartUseCase, env *appconfig.Env)
		expectedCode int
		expectedBody interface{}
	}{
		{
			name: "Successful SignUp",
			requestBody: domain.SignupRequest{
				Name:     "Nguyen Vo Tien Loc",
				Email:    "test@gm.uit.edu.vn",
				Password: "securepassword",
			},
			setupMocks: func(mockSignupUseCase *MockSignupUseCase, mockCartUseCase *MockCartUseCase, env *appconfig.Env) {
				mockSignupUseCase.On("SignUp", mock.Anything, mock.Anything, env).
					Return("access-token", "refresh-token", "675097b7c62c4aef0cfa828b", nil)
				mockCartUseCase.On("CreateCartByUserId", mock.Anything, mock.Anything).
					Return(nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"code":    float64(201),
				"message": "Registered successfully",
				"data": map[string]interface{}{
					"access_token":  "access-token",
					"refresh_token": "refresh-token",
				},
			},
		},
		{
			name: "Invalid Request Body",
			requestBody: domain.SignupRequest{
				Name:     "Nguyen Vo Tien Loc",
				Email:    "test@hii",
				Password: "securepassword",
			},
			setupMocks:   func(mockSignupUseCase *MockSignupUseCase, mockCartUseCase *MockCartUseCase, env *appconfig.Env) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(http.StatusBadRequest),
				"error":   "Key: 'SignupRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag",
				"message": "Cannot parse the request body",
			},
		},
		{
			name: "Email Already Exists",
			requestBody: domain.SignupRequest{
				Name:     "Nguyen Vo Tien Loc",
				Email:    "test@gm.uit.edu.vn",
				Password: "securepassword",
			},
			setupMocks: func(mockSignupUseCase *MockSignupUseCase, mockCartUseCase *MockCartUseCase, env *appconfig.Env) {
				mockSignupUseCase.On("SignUp", mock.Anything, mock.Anything, env).
					Return("", "", "", apperror.ErrEmailAlreadyExist)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(http.StatusBadRequest),
				"error":   "email already exist",
				"message": "Email already exists",
			},
		},
		{
			name: "Invalid User ID",
			requestBody: domain.SignupRequest{
				Name:     "Nguyen Vo Tien Loc",
				Email:    "test@example.com",
				Password: "securepassword",
			},
			setupMocks: func(mockSignupUseCase *MockSignupUseCase, mockCartUseCase *MockCartUseCase, env *appconfig.Env) {
				mockSignupUseCase.On("SignUp", mock.Anything, mock.Anything, env).
					Return("access-token", "refresh-token", "invalid-object-id", nil)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(http.StatusBadRequest),
				"error":   "the provided hex string is not a valid ObjectID",
				"message": "Invalid user ID",
			},
		},
		{
			name: "Error Creating Cart",
			requestBody: domain.SignupRequest{
				Name:     "Nguyen Vo Tien Loc",
				Email:    "test@example.com",
				Password: "securepassword",
			},
			setupMocks: func(mockSignupUseCase *MockSignupUseCase, mockCartUseCase *MockCartUseCase, env *appconfig.Env) {
				mockSignupUseCase.On("SignUp", mock.Anything, mock.Anything, env).
					Return("access-token", "refresh-token", "645a0c2f4a9a2d3a58c96b84", nil)
				mockCartUseCase.On("CreateCartByUserId", mock.Anything, mock.Anything).
					Return(errors.New("cart creation error"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(http.StatusBadRequest),
				"error":   "cart creation error",
				"message": "Error creating empty cart",
			},
		},
		{
			name: "Missing Password",
			requestBody: domain.SignupRequest{
				Name:  "Nguyen Vo Tien Loc",
				Email: "test@gm.uit.edu.vn",
			},
			setupMocks:   func(mockSignupUseCase *MockSignupUseCase, mockCartUseCase *MockCartUseCase, env *appconfig.Env) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(http.StatusBadRequest),
				"error":   "Key: 'SignupRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag",
				"message": "Cannot parse the request body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks and controller
			mockSignupUseCase := new(MockSignupUseCase)
			mockCartUseCase := new(MockCartUseCase)
			env := &appconfig.Env{}
			controller := handlers.SignupController{
				SignupUseCase: mockSignupUseCase,
				CartUseCase:   mockCartUseCase,
				Env:           env,
			}

			// Setup mocks for each test case
			tt.setupMocks(mockSignupUseCase, mockCartUseCase, env)

			// Register cleanup to assert mock expectations after test
			t.Cleanup(func() {
				mockSignupUseCase.AssertExpectations(t)
				mockCartUseCase.AssertExpectations(t)
			})

			// Create a request
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			recorder := httptest.NewRecorder()
			router := gin.Default()
			router.POST("/signup", controller.SignUp())

			// Perform the request
			router.ServeHTTP(recorder, req)

			// Assert response code
			assert.Equal(t, tt.expectedCode, recorder.Code)

			// Assert response body
			var actualBody interface{}
			if recorder.Code == 200 || recorder.Code == 201 {
				actualBody = response.Success{}
			} else {
				actualBody = response.FailedResponse{}
			}
			_ = json.Unmarshal(recorder.Body.Bytes(), &actualBody)
			assert.Equal(t, tt.expectedBody, actualBody)
		})
	}
}
