package missions_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"spy-cats/internal/missions"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) CreateMission(req missions.CreateMissionRequest) (*missions.Mission, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*missions.Mission), args.Error(1)
}

func (m *mockService) DeleteMission(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockService) MarkMissionComplete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockService) AddTarget(missionID int64, req missions.CreateTarget) error {
	args := m.Called(missionID, req)
	return args.Error(0)
}

func (m *mockService) UpdateTarget(id int64, req missions.UpdateTargetRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *mockService) DeleteTarget(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockService) GetAllMissions() ([]missions.Mission, error) {
	args := m.Called()
	return args.Get(0).([]missions.Mission), args.Error(1)
}

func (m *mockService) GetMissionByID(id int64) (*missions.Mission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*missions.Mission), args.Error(1)
}

func (m *mockService) AssignCat(missionID, catID int64) error {
	args := m.Called(missionID, catID)
	return args.Error(0)
}

func TestCreateMission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name              string
		body              any
		mockReturnMission *missions.Mission
		mockReturnErr     error
		expectedStatus    int
		expectedBody      string
	}{
		{
			name: "success",
			body: missions.CreateMissionRequest{
				CatID: nil,
				Name:  "Operation Stealth",
				Targets: []missions.CreateTarget{
					{
						Name:    "Target Alpha",
						Country: "Russia",
						Notes:   "High priority",
					},
				},
				IsComplete: false,
			},
			mockReturnMission: &missions.Mission{
				ID:         1,
				CatID:      nil,
				Name:       "Operation Stealth",
				IsComplete: false,
				Targets: []missions.Target{
					{
						ID:         1,
						MissionID:  1,
						Name:       "Target Alpha",
						Country:    "Russia",
						Notes:      "High priority",
						IsComplete: false,
					},
				},
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `"name":"Operation Stealth"`,
		},
		{
			name: "success with cat assigned",
			body: missions.CreateMissionRequest{
				CatID: func() *int64 { id := int64(5); return &id }(),
				Name:  "Operation Silent Paws",
				Targets: []missions.CreateTarget{
					{
						Name:    "Target Beta",
						Country: "China",
						Notes:   "Surveillance only",
					},
				},
				IsComplete: false,
			},
			mockReturnMission: &missions.Mission{
				ID:         2,
				CatID:      func() *int64 { id := int64(5); return &id }(),
				Name:       "Operation Silent Paws",
				IsComplete: false,
				Targets: []missions.Target{
					{
						ID:         2,
						MissionID:  2,
						Name:       "Target Beta",
						Country:    "China",
						Notes:      "Surveillance only",
						IsComplete: false,
					},
				},
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `"name":"Operation Silent Paws"`,
		},
		{
			name:           "invalid JSON body",
			body:           `{"name": 123}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error"`,
		},
		{
			name: "missing required fields",
			body: missions.CreateMissionRequest{
				CatID:      nil,
				Name:       "",                        // Missing required name
				Targets:    []missions.CreateTarget{}, // Missing required targets
				IsComplete: false,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error"`,
		},
		{
			name: "service error",
			body: missions.CreateMissionRequest{
				CatID: nil,
				Name:  "Operation Fail",
				Targets: []missions.CreateTarget{
					{
						Name:    "Target Gamma",
						Country: "Iran",
						Notes:   "Will fail",
					},
				},
				IsComplete: false,
			},
			mockReturnMission: nil,
			mockReturnErr:     errors.New("database connection error"),
			expectedStatus:    http.StatusInternalServerError,
			expectedBody:      `"error":"database connection error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := missions.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.POST("/missions", h.CreateMission)

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				b, _ := json.Marshal(v)
				bodyBytes = b
			}

			if tt.mockReturnMission != nil || tt.mockReturnErr != nil {
				mockSvc.On("CreateMission", mock.Anything).Return(tt.mockReturnMission, tt.mockReturnErr)
			}

			req, _ := http.NewRequest(http.MethodPost, "/missions", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestDeleteMission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		missionID      string
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success",
			missionID:      "1",
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"message":"mission deleted"`,
		},
		{
			name:           "mission not found",
			missionID:      "999",
			mockReturnErr:  errors.New("mission not found"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"mission not found"`,
		},
		{
			name:           "service error",
			missionID:      "1",
			mockReturnErr:  errors.New("database error"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"database error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := missions.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.DELETE("/missions/:id", h.DeleteMission)

			mockSvc.On("DeleteMission", mock.AnythingOfType("int64")).Return(tt.mockReturnErr)

			req, _ := http.NewRequest(http.MethodDelete, "/missions/"+tt.missionID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestMarkMissionComplete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		missionID      string
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success",
			missionID:      "1",
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"message":"mission marked complete"`,
		},
		{
			name:           "mission not found",
			missionID:      "999",
			mockReturnErr:  sql.ErrNoRows,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"error":"mission not found"`,
		},
		{
			name:           "service error",
			missionID:      "1",
			mockReturnErr:  errors.New("database connection error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error":"database connection error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := missions.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.PUT("/missions/:id/complete", h.MarkMissionComplete)

			mockSvc.On("MarkMissionComplete", mock.AnythingOfType("int64")).Return(tt.mockReturnErr)

			req, _ := http.NewRequest(http.MethodPut, "/missions/"+tt.missionID+"/complete", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestAddTarget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		missionID      string
		body           any
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "success",
			missionID: "1",
			body: missions.CreateTarget{
				Name:       "New Target",
				Country:    "France",
				Notes:      "High priority target",
				IsComplete: false,
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `"message":"target added"`,
		},
		{
			name:           "invalid JSON body",
			missionID:      "1",
			body:           `{"name": 123}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error"`,
		},
		{
			name:      "missing required fields",
			missionID: "1",
			body: missions.CreateTarget{
				Name:    "", // Missing required name
				Country: "", // Missing required country
				Notes:   "Some notes",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error"`,
		},
		{
			name:      "service error",
			missionID: "1",
			body: missions.CreateTarget{
				Name:    "Target Delta",
				Country: "Germany",
				Notes:   "Service will fail",
			},
			mockReturnErr:  errors.New("cannot add target to completed mission"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"cannot add target to completed mission"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := missions.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.POST("/missions/:id/targets", h.AddTarget)

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				b, _ := json.Marshal(v)
				bodyBytes = b
			}

			if tt.mockReturnErr != nil || tt.expectedStatus == http.StatusCreated {
				mockSvc.On("AddTarget", mock.AnythingOfType("int64"), mock.Anything).Return(tt.mockReturnErr)
			}

			req, _ := http.NewRequest(http.MethodPost, "/missions/"+tt.missionID+"/targets", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestUpdateTarget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		targetID       string
		body           any
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "success",
			targetID: "1",
			body: missions.UpdateTargetRequest{
				IsComplete: func() *bool { b := true; return &b }(),
				Notes:      func() *string { s := "Mission accomplished"; return &s }(),
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"message":"target updated"`,
		},
		{
			name:     "success with partial update",
			targetID: "2",
			body: missions.UpdateTargetRequest{
				IsComplete: func() *bool { b := false; return &b }(),
				Notes:      nil, // Only updating completion status
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"message":"target updated"`,
		},
		{
			name:           "invalid JSON body",
			targetID:       "1",
			body:           `{"is_complete": "invalid"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error"`,
		},
		{
			name:     "service error",
			targetID: "1",
			body: missions.UpdateTargetRequest{
				IsComplete: func() *bool { b := true; return &b }(),
			},
			mockReturnErr:  errors.New("target not found"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"target not found"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := missions.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.PUT("/missions/:id/targets/:targetId", h.UpdateTarget)

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				b, _ := json.Marshal(v)
				bodyBytes = b
			}

			if tt.mockReturnErr != nil || tt.expectedStatus == http.StatusOK {
				mockSvc.On("UpdateTarget", mock.AnythingOfType("int64"), mock.Anything).Return(tt.mockReturnErr)
			}

			req, _ := http.NewRequest(http.MethodPut, "/missions/1/targets/"+tt.targetID, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestDeleteTarget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		targetID       string
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success",
			targetID:       "1",
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"message":"target deleted"`,
		},
		{
			name:           "target not found",
			targetID:       "999",
			mockReturnErr:  errors.New("target not found"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"target not found"`,
		},
		{
			name:           "service error",
			targetID:       "1",
			mockReturnErr:  errors.New("database error"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"database error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := missions.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.DELETE("/missions/:id/targets/:targetId", h.DeleteTarget)

			mockSvc.On("DeleteTarget", mock.AnythingOfType("int64")).Return(tt.mockReturnErr)

			req, _ := http.NewRequest(http.MethodDelete, "/missions/1/targets/"+tt.targetID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestGetAllMissions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		mockReturnMissions []missions.Mission
		mockReturnErr      error
		expectedStatus     int
		expectedBody       string
	}{
		{
			name: "success with missions",
			mockReturnMissions: []missions.Mission{
				{
					ID:         1,
					CatID:      func() *int64 { id := int64(5); return &id }(),
					Name:       "Operation Alpha",
					IsComplete: false,
					Targets: []missions.Target{
						{ID: 1, MissionID: 1, Name: "Target A", Country: "Russia", Notes: "High priority", IsComplete: false},
					},
				},
				{
					ID:         2,
					CatID:      nil,
					Name:       "Operation Beta",
					IsComplete: true,
					Targets: []missions.Target{
						{ID: 2, MissionID: 2, Name: "Target B", Country: "China", Notes: "Completed", IsComplete: true},
					},
				},
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"name":"Operation Alpha"`,
		},
		{
			name:               "success with empty list",
			mockReturnMissions: []missions.Mission{},
			mockReturnErr:      nil,
			expectedStatus:     http.StatusOK,
			expectedBody:       `[]`,
		},
		{
			name:               "service error",
			mockReturnMissions: []missions.Mission{},
			mockReturnErr:      errors.New("database connection error"),
			expectedStatus:     http.StatusInternalServerError,
			expectedBody:       `"error":"failed to fetch missions"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := missions.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.GET("/missions", h.GetAllMissions)

			mockSvc.On("GetAllMissions").Return(tt.mockReturnMissions, tt.mockReturnErr)

			req, _ := http.NewRequest(http.MethodGet, "/missions", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestGetMissionByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name              string
		missionID         string
		mockReturnMission *missions.Mission
		mockReturnErr     error
		expectedStatus    int
		expectedBody      string
	}{
		{
			name:      "success",
			missionID: "1",
			mockReturnMission: &missions.Mission{
				ID:         1,
				CatID:      func() *int64 { id := int64(3); return &id }(),
				Name:       "Operation Gamma",
				IsComplete: false,
				Targets: []missions.Target{
					{ID: 1, MissionID: 1, Name: "Target Gamma", Country: "Iran", Notes: "Surveillance", IsComplete: false},
				},
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"name":"Operation Gamma"`,
		},
		{
			name:              "mission not found",
			missionID:         "999",
			mockReturnMission: nil,
			mockReturnErr:     sql.ErrNoRows,
			expectedStatus:    http.StatusNotFound,
			expectedBody:      `"error":"mission not found"`,
		},
		{
			name:              "service error",
			missionID:         "1",
			mockReturnMission: nil,
			mockReturnErr:     errors.New("database error"),
			expectedStatus:    http.StatusInternalServerError,
			expectedBody:      `"error":"failed to fetch mission"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := missions.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.GET("/missions/:id", h.GetMissionByID)

			mockSvc.On("GetMissionByID", mock.AnythingOfType("int64")).Return(tt.mockReturnMission, tt.mockReturnErr)

			req, _ := http.NewRequest(http.MethodGet, "/missions/"+tt.missionID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestAssignCat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		missionID      string
		body           any
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "success",
			missionID: "1",
			body: missions.AssignCatRequest{
				CatID: 5,
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"message":"cat assigned successfully"`,
		},
		{
			name:           "invalid JSON body",
			missionID:      "1",
			body:           `{"cat_id": "invalid"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error"`,
		},
		{
			name:      "missing required field",
			missionID: "1",
			body:      missions.AssignCatRequest{
				// Missing CatID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error"`,
		},
		{
			name:      "mission not found",
			missionID: "999",
			body: missions.AssignCatRequest{
				CatID: 5,
			},
			mockReturnErr:  sql.ErrNoRows,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"error":"mission not found"`,
		},
		{
			name:      "service error",
			missionID: "1",
			body: missions.AssignCatRequest{
				CatID: 999, // Cat doesn't exist
			},
			mockReturnErr:  errors.New("cat not found"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error":"failed to assign cat"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			h := missions.NewHandler(mockSvc)

			// prepare Gin router
			r := gin.Default()
			r.PUT("/missions/:id/assign", h.AssignCat)

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				b, _ := json.Marshal(v)
				bodyBytes = b
			}

			if tt.mockReturnErr != nil || tt.expectedStatus == http.StatusOK {
				mockSvc.On("AssignCat", mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return(tt.mockReturnErr)
			}

			req, _ := http.NewRequest(http.MethodPut, "/missions/"+tt.missionID+"/assign", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}
