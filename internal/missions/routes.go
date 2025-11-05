package missions

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, db *sql.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.POST("/", handler.CreateMission)
	r.GET("/", handler.GetAllMissions)
	r.GET("/:id", handler.GetMissionByID)
	r.PUT("/:id/assign", handler.AssignCat)
	r.DELETE("/:id", handler.DeleteMission)
	r.PATCH("/:id/complete", handler.MarkMissionComplete)

	r.POST("/:id/targets", handler.AddTarget)
	r.PATCH("/targets/:targetId", handler.UpdateTarget)
	r.DELETE("/targets/:targetId", handler.DeleteTarget)
}
