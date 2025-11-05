package missions

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MissionService interface {
	CreateMission(req CreateMissionRequest) (*Mission, error)
	DeleteMission(id int64) error
	MarkMissionComplete(id int64) error
	AddTarget(missionID int64, req CreateTarget) error
	UpdateTarget(id int64, req UpdateTargetRequest) error
	DeleteTarget(id int64) error
	GetAllMissions() ([]Mission, error)
	GetMissionByID(id int64) (*Mission, error)
	AssignCat(missionID, catID int64) error
}

type Handler struct {
	service MissionService
}

func NewHandler(service MissionService) *Handler {
	return &Handler{service: service}
}

// CreateMission creates a new mission
// @Summary      Create a mission
// @Description  Create a new mission with targets
// @Tags         missions
// @Accept       json
// @Produce      json
// @Param        mission  body      CreateMissionRequest  true  "Mission information"
// @Success      201      {object}  Mission               "Successfully created mission"
// @Failure      400      {object}  map[string]string     "Invalid input"
// @Failure      500      {object}  map[string]string     "Internal server error"
// @Router       /missions [post]
func (h *Handler) CreateMission(c *gin.Context) {
	var req CreateMissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mission, err := h.service.CreateMission(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, mission)
}

// DeleteMission deletes a mission
// @Summary      Delete a mission
// @Description  Delete a mission by its ID
// @Tags         missions
// @Param        id  path      int  true  "Mission ID"
// @Success      200 {object}  map[string]string "Mission deleted successfully"
// @Failure      400 {object}  map[string]string "Bad request"
// @Router       /missions/{id} [delete]
func (h *Handler) DeleteMission(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.DeleteMission(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "mission deleted"})
}

// MarkMissionComplete marks a mission as complete
// @Summary      Mark mission complete
// @Description  Mark a mission as complete by its ID
// @Tags         missions
// @Param        id  path      int  true  "Mission ID"
// @Success      200 {object}  map[string]string "Mission marked complete"
// @Failure      404 {object}  map[string]string "Mission not found"
// @Failure      500 {object}  map[string]string "Internal server error"
// @Router       /missions/{id}/complete [patch]
func (h *Handler) MarkMissionComplete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.MarkMissionComplete(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "mission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "mission marked complete"})
}

// AddTarget adds a target to a mission
// @Summary      Add target to mission
// @Description  Add a new target to an existing mission
// @Tags         missions
// @Accept       json
// @Produce      json
// @Param        id      path      int           true  "Mission ID"
// @Param        target  body      CreateTarget  true  "Target information"
// @Success      201     {object}  map[string]string "Target added successfully"
// @Failure      400     {object}  map[string]string "Bad request"
// @Router       /missions/{id}/targets [post]
func (h *Handler) AddTarget(c *gin.Context) {
	missionID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req CreateTarget
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AddTarget(missionID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "target added"})
}

// UpdateTarget updates a target
// @Summary      Update target
// @Description  Update target completion status and notes
// @Tags         missions
// @Accept       json
// @Produce      json
// @Param        targetId  path      int                   true  "Target ID"
// @Param        target    body      UpdateTargetRequest   true  "Target update information"
// @Success      200       {object}  map[string]string     "Target updated successfully"
// @Failure      400       {object}  map[string]string     "Bad request"
// @Router       /missions/targets/{targetId} [patch]
func (h *Handler) UpdateTarget(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("targetId"), 10, 64)
	var req UpdateTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateTarget(id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "target updated"})
}

// DeleteTarget deletes a target
// @Summary      Delete target
// @Description  Delete a target by its ID
// @Tags         missions
// @Param        targetId  path      int  true  "Target ID"
// @Success      200       {object}  map[string]string "Target deleted successfully"
// @Failure      400       {object}  map[string]string "Bad request"
// @Router       /missions/targets/{targetId} [delete]
func (h *Handler) DeleteTarget(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("targetId"), 10, 64)
	if err := h.service.DeleteTarget(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "target deleted"})
}

// GetAllMissions retrieves all missions
// @Summary      List all missions
// @Description  Get a list of all missions in the system
// @Tags         missions
// @Produce      json
// @Success      200  {array}   Mission "List of missions"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /missions [get]
func (h *Handler) GetAllMissions(c *gin.Context) {
	missions, err := h.service.GetAllMissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch missions"})
		return
	}
	c.JSON(http.StatusOK, missions)
}

// GetMissionByID retrieves a specific mission by ID
// @Summary      Get a mission
// @Description  Get a mission by its ID
// @Tags         missions
// @Produce      json
// @Param        id   path      int  true  "Mission ID"
// @Success      200  {object}  Mission "Mission information"
// @Failure      404  {object}  map[string]string "Mission not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /missions/{id} [get]
func (h *Handler) GetMissionByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	mission, err := h.service.GetMissionByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "mission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch mission"})
		return
	}
	c.JSON(http.StatusOK, mission)
}

// AssignCat assigns a cat to a mission
// @Summary      Assign cat to mission
// @Description  Assign a spy cat to a mission
// @Tags         missions
// @Accept       json
// @Produce      json
// @Param        id       path      int               true  "Mission ID"
// @Param        request  body      AssignCatRequest  true  "Cat assignment information"
// @Success      200      {object}  map[string]string "Cat assigned successfully"
// @Failure      400      {object}  map[string]string "Bad request"
// @Failure      404      {object}  map[string]string "Mission not found"
// @Failure      500      {object}  map[string]string "Internal server error"
// @Router       /missions/{id}/assign [put]
func (h *Handler) AssignCat(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req AssignCatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AssignCat(id, req.CatID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "mission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign cat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cat assigned successfully"})
}
