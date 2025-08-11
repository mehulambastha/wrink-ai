package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/mehulambastha/wrink-ai/backend/api"
	"github.com/mehulambastha/wrink-ai/backend/config"
	"github.com/mehulambastha/wrink-ai/backend/internal/database"
	"github.com/mehulambastha/wrink-ai/backend/internal/models"
	"github.com/mehulambastha/wrink-ai/backend/middleware"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config.LoadConfig()
	database.ConnectDB()
}

func main() {
	err := database.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Suggestion{}, &models.SearchResult{})

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	v1 := router.Group("/api/v1")

	{
		v1.POST("/register", api.CreateUser)
		v1.POST("/login", api.LoginUser)
	}

	postRoutes := v1.Group("/post")
	postRoutes.Use(middleware.AuthMiddleware())
	{
		postRoutes.POST("suggestion", api.CreateSuggestion)
	}

	router.Run(":5000")
}
