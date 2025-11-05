package main

import (
	"log"

	_ "spy-cats/docs" // Import docs for swagger
	"spy-cats/internal/cats"
	"spy-cats/internal/database"
	"spy-cats/internal/middleware"
	"spy-cats/internal/missions"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Spy Cats API
// @version         1.0
// @description     A mission management system for spy cats
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   http://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	r := gin.Default()
	r.Use(middleware.LoggingMiddleware())

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		cats.RegisterRoutes(api.Group("/cats"), db)
		missions.RegisterRoutes(api.Group("/missions"), db)
	}

	log.Println("Server started on port 8080")
	log.Println("Swagger UI available at: http://localhost:8080/swagger/index.html")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
