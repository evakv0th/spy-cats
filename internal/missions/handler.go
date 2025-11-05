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

func (h *Handler) DeleteMission(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.DeleteMission(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "mission deleted"})
}

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

func (h *Handler) DeleteTarget(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("targetId"), 10, 64)
	if err := h.service.DeleteTarget(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "target deleted"})
}

func (h *Handler) GetAllMissions(c *gin.Context) {
	missions, err := h.service.GetAllMissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch missions"})
		return
	}
	c.JSON(http.StatusOK, missions)
}

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
