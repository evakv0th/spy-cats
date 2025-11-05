package cats_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"spy-cats/internal/cats"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) CreateCat(req cats.CreateCatRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockService) GetAllCats() ([]cats.Cat, error) {
	args := m.Called()
	return args.Get(0).([]cats.Cat), args.Error(1)
}
func (m *mockService) GetCat(id int64) (*cats.Cat, error) {
	args := m.Called(id)
	return args.Get(0).(*cats.Cat), args.Error(1)
}
func (m *mockService) UpdateSalary(id int64, salary float64) error {
	args := m.Called(id, salary)
	return args.Error(0)
}
func (m *mockService) DeleteCat(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateCat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           any
		mockReturnID   int64
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			body: cats.CreateCatRequest{
				Name:              "Whiskers",
				YearsOfExperience: 5,
				Breed:             "Siamese",
				Salary:            1000.0,
			},
			mockReturnID:   42,
			mockReturnErr:  nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `"id":42`,
		},
		{
			name:           "invalid JSON body",
			body:           `{"name": 123}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error"`,
		},
		{
			name: "invalid breed",
			body: cats.CreateCatRequest{
				Name:              "Tom",
				YearsOfExperience: 3,
				Breed:             "UnknownBreed",
				Salary:            500.0,
			},
			mockReturnID:   0,
			mockReturnErr:  errors.New("invalid cat breed: UnknownBreed"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"invalid cat breed: UnknownBreed"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := cats.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.POST("/cats", h.CreateCat)

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				b, _ := json.Marshal(v)
				bodyBytes = b
			}

			req, _ := http.NewRequest(http.MethodPost, "/cats", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			if req.Header.Get("Content-Type") == "application/json" && tt.mockReturnErr != nil || tt.mockReturnID != 0 {
				mockSvc.On("CreateCat", mock.Anything).Return(tt.mockReturnID, tt.mockReturnErr)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestListCats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockReturnCats []cats.Cat
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success with cats",
			mockReturnCats: []cats.Cat{
				{ID: 1, Name: "Whiskers", YearsOfExperience: 5, Breed: "Siamese", Salary: 1000.0},
				{ID: 2, Name: "Mittens", YearsOfExperience: 3, Breed: "Persian", Salary: 800.0},
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"name":"Whiskers"`,
		},
		{
			name:           "success with empty list",
			mockReturnCats: []cats.Cat{},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `[]`,
		},
		{
			name:           "service error",
			mockReturnCats: []cats.Cat{},
			mockReturnErr:  errors.New("database connection error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error":"failed to get cats"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := cats.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.GET("/cats", h.ListCats)

			mockSvc.On("GetAllCats").Return(tt.mockReturnCats, tt.mockReturnErr)

			req, _ := http.NewRequest(http.MethodGet, "/cats", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestGetCat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		catID          string
		mockReturnCat  *cats.Cat
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:  "success",
			catID: "1",
			mockReturnCat: &cats.Cat{
				ID:                1,
				Name:              "Whiskers",
				YearsOfExperience: 5,
				Breed:             "Siamese",
				Salary:            1000.0,
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"name":"Whiskers"`,
		},
		{
			name:           "cat not found",
			catID:          "999",
			mockReturnCat:  nil,
			mockReturnErr:  nil,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"error":"cat not found with id 999"`,
		},
		{
			name:           "service error",
			catID:          "1",
			mockReturnCat:  nil,
			mockReturnErr:  errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error":"failed to get cat"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := cats.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.GET("/cats/:id", h.GetCat)

			mockSvc.On("GetCat", mock.AnythingOfType("int64")).Return(tt.mockReturnCat, tt.mockReturnErr)

			req, _ := http.NewRequest(http.MethodGet, "/cats/"+tt.catID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestUpdateSalary(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		catID          string
		body           any
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:  "success",
			catID: "1",
			body: cats.UpdateSalaryRequest{
				Salary: 1500.0,
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"message":"salary updated successfully"`,
		},
		{
			name:           "invalid cat ID",
			catID:          "invalid",
			body:           cats.UpdateSalaryRequest{Salary: 1500.0},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"invalid cat id"`,
		},
		{
			name:           "invalid JSON body",
			catID:          "1",
			body:           `{"salary": "invalid"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error"`,
		},
		{
			name:  "cat not found",
			catID: "999",
			body: cats.UpdateSalaryRequest{
				Salary: 1500.0,
			},
			mockReturnErr:  errors.New("cat not found with id 999"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"error":"cat not found with id 999"`,
		},
		{
			name:  "service error",
			catID: "1",
			body: cats.UpdateSalaryRequest{
				Salary: 1500.0,
			},
			mockReturnErr:  errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error":"failed to update salary"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := cats.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.PUT("/cats/:id/salary", h.UpdateSalary)

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				b, _ := json.Marshal(v)
				bodyBytes = b
			}

			if tt.catID != "invalid" && tt.body != `{"salary": "invalid"}` {
				mockSvc.On("UpdateSalary", mock.AnythingOfType("int64"), mock.AnythingOfType("float64")).Return(tt.mockReturnErr)
			}

			req, _ := http.NewRequest(http.MethodPut, "/cats/"+tt.catID+"/salary", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestDeleteCat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		catID          string
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success",
			catID:          "1",
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"message":"cat deleted"`,
		},
		{
			name:           "service error",
			catID:          "1",
			mockReturnErr:  errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error":"failed to delete cat"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := cats.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.DELETE("/cats/:id", h.DeleteCat)

			mockSvc.On("DeleteCat", mock.AnythingOfType("int64")).Return(tt.mockReturnErr)

			req, _ := http.NewRequest(http.MethodDelete, "/cats/"+tt.catID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}
