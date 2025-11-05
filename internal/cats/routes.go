package cats

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	rg.POST("/", handler.CreateCat)
	rg.GET("/", handler.ListCats)
	rg.GET("/:id", handler.GetCat)
	rg.PATCH("/:id/salary", handler.UpdateSalary)
	rg.DELETE("/:id", handler.DeleteCat)
}
