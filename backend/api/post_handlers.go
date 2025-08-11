package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mehulambastha/wrink-ai/backend/internal/models"
	"github.com/mehulambastha/wrink-ai/backend/internal/services"
)

func CreateSuggestion(c *gin.Context) {
	var input models.CreateSuggestionDto

	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserID = c.GetUint("user_id")

	suggestion, err := services.CreateSuggestion(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create suggestion"})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, suggestion)
}
