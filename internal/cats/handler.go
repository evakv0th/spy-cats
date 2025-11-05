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

// CreateCat creates a new spy cat
// @Summary      Create a spy cat
// @Description  Create a new spy cat with the provided information
// @Tags         cats
// @Accept       json
// @Produce      json
// @Param        cat  body      CreateCatRequest  true  "Cat information"
// @Success      201  {object}  map[string]int64  "Successfully created cat"
// @Failure      400  {object}  map[string]string "Invalid input"
// @Router       /cats [post]
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

// ListCats retrieves all spy cats
// @Summary      List all spy cats
// @Description  Get a list of all spy cats in the system
// @Tags         cats
// @Produce      json
// @Success      200  {array}   Cat "List of cats"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /cats [get]
func (h *Handler) ListCats(c *gin.Context) {
	cats, err := h.service.GetAllCats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cats"})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// GetCat retrieves a specific spy cat by ID
// @Summary      Get a spy cat
// @Description  Get a spy cat by its ID
// @Tags         cats
// @Produce      json
// @Param        id   path      int  true  "Cat ID"
// @Success      200  {object}  Cat "Cat information"
// @Failure      404  {object}  map[string]string "Cat not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /cats/{id} [get]
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

// UpdateSalary updates a spy cat's salary
// @Summary      Update cat salary
// @Description  Update the salary of a specific spy cat
// @Tags         cats
// @Accept       json
// @Produce      json
// @Param        id      path      int                  true  "Cat ID"
// @Param        salary  body      UpdateSalaryRequest  true  "New salary information"
// @Success      200     {object}  map[string]string    "Salary updated successfully"
// @Failure      400     {object}  map[string]string    "Invalid input"
// @Failure      404     {object}  map[string]string    "Cat not found"
// @Failure      500     {object}  map[string]string    "Internal server error"
// @Router       /cats/{id}/salary [patch]
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

// DeleteCat deletes a spy cat
// @Summary      Delete a spy cat
// @Description  Delete a spy cat by its ID
// @Tags         cats
// @Param        id  path      int  true  "Cat ID"
// @Success      200 {object}  map[string]string "Cat deleted successfully"
// @Failure      500 {object}  map[string]string "Internal server error"
// @Router       /cats/{id} [delete]
func (h *Handler) DeleteCat(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.service.DeleteCat(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete cat"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cat deleted"})
}
