package cats

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service CatService
}

type CatService interface {
	CreateCat(req CreateCatRequest) (int64, error)
	GetAllCats() ([]Cat, error)
	GetCat(id int64) (*Cat, error)
	UpdateSalary(id int64, salary float64) error
	DeleteCat(id int64) error
}

func NewHandler(service CatService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateCat(c *gin.Context) {
	var req CreateCatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.CreateCat(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) ListCats(c *gin.Context) {
	cats, err := h.service.GetAllCats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cats"})
		return
	}
	c.JSON(http.StatusOK, cats)
}

func (h *Handler) GetCat(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cat, err := h.service.GetCat(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cat"})
		return
	}
	if cat == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cat not found with id " + c.Param("id")})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *Handler) UpdateSalary(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cat id"})
		return
	}

	var req UpdateSalaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateSalary(id, req.Salary)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "cat not found with id " + c.Param("id")})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update salary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "salary updated successfully"})
}

func (h *Handler) DeleteCat(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.DeleteCat(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete cat"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cat deleted"})
}
