package main

import (
	"log"

	"spy-cats/internal/cats"
	"spy-cats/internal/database"
	"spy-cats/internal/middleware"
	"spy-cats/internal/missions"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	r := gin.Default()
	r.Use(middleware.LoggingMiddleware())

	api := r.Group("/api")
	{
		cats.RegisterRoutes(api.Group("/cats"), db)
		missions.RegisterRoutes(api.Group("/missions"), db)
	}

	log.Println("Server started on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
