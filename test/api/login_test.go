package api_test

import (
	"bytes"
	"cakewai/cakewai.com/api/handlers"
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/domain"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
)

type MockLoginUseCase struct {
	mock.Mock
}

func (m *MockLoginUseCase) Login(c context.Context, request domain.LoginRequest, env *appconfig.Env) (string, string, error) {
	args := m.Called(c, request, env)
	return args.String(0), args.String(1), args.Error(2)
}
func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test cases
	tests := []struct {
		name         string
		requestBody  interface{}
		setupMocks   func(mockLoginUseCase *MockLoginUseCase, env *appconfig.Env)
		expectedCode int
		expectedBody interface{}
	}{
		{
			name: "Successful Login",
			requestBody: domain.LoginRequest{
				Email:    "test@gm.uit.edu.vn",
				Password: "securepassword",
			},
			setupMocks: func(mockLoginUseCase *MockLoginUseCase, env *appconfig.Env) {
				mockLoginUseCase.On("Login", mock.Anything, mock.Anything, env).
					Return("access-token", "refresh-token", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code":    float64(http.StatusOK),
				"message": "Login successful",
				"data": map[string]interface{}{
					"access_token":  "access-token",
					"refresh_token": "refresh-token",
				},
			},
		},
		{
			name: "Invalid Request Body",
			requestBody: domain.LoginRequest{
				Email:    "invalid-email",
				Password: "securepassword",
			},
			setupMocks:   func(mockLoginUseCase *MockLoginUseCase, env *appconfig.Env) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code":    float64(http.StatusBadRequest),
				"error":   "Key: 'LoginRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag",
				"message": "Invalid request payload",
			},
		},
		{
			name: "Invalid Credentials",
			requestBody: domain.LoginRequest{
				Email:    "test@gm.uit.edu.vn",
				Password: "wrongpassword",
			},
			setupMocks: func(mockLoginUseCase *MockLoginUseCase, env *appconfig.Env) {
				mockLoginUseCase.On("Login", mock.Anything, mock.Anything, env).
					Return("", "", errors.New("http.StatusUnauthorized"))
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"code":    float64(http.StatusUnauthorized),
				"message": "Invalid credentials",
				"error":   "http.StatusUnauthorized",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks and controller
			mockLoginUseCase := new(MockLoginUseCase)
			env := &appconfig.Env{}
			handler := handlers.LoginHandler{
				LoginUsecase: mockLoginUseCase,
				Env:          env,
			}

			// Setup mocks for each test case
			tt.setupMocks(mockLoginUseCase, env)

			// Register cleanup to assert mock expectations after test
			t.Cleanup(func() {
				mockLoginUseCase.AssertExpectations(t)
			})

			// Create a request
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			recorder := httptest.NewRecorder()
			router := gin.Default()
			router.POST("/login", handler.LoginHandler())

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
