package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mehulambastha/wrink-ai/backend/internal/models"
	"github.com/mehulambastha/wrink-ai/backend/internal/services"
)

func CreateUser(c *gin.Context) {
	var input models.CreateUserDto

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.CreateUser(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, user)
}

func LoginUser(c *gin.Context) {
	var input models.LoginDto

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := services.LoginUser(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{Token: token})

}
