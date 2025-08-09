package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mehulambastha/wrink-ai/backend/internal/database"
	"github.com/mehulambastha/wrink-ai/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(input models.CreateUserDto) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	newUser := models.User{
		Name:           input.Name,
		Email:          input.Email,
		HashedPassword: string(hashedPassword),
	}

	result := database.DB.Create(&newUser)

	if result.Error != nil {
		return nil, result.Error
	}

	return &newUser, nil
}

func LoginUser(input models.LoginDto) (string, error) {
	var user models.User

	// find if user exists
	result := database.DB.Where("email = ?", input.Email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", errors.New("Invalid Credentials")
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(input.Password))

	if err != nil {
		return "", errors.New("Invalid Credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	jwtString := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtString))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
