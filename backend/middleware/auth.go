package middleware

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mehulambastha/wrink-ai/backend/config"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Inside middleware")
		authHeader := c.GetHeader("Authorization")
		log.Println("Auth header is: ", authHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			return
		}

		log.Println("Second part of string: ", parts[1])

		userID, err := validateToken(parts[1])

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		log.Println("User id is: ", userID)
		log.Println("And its type is: ", reflect.TypeOf(userID))

		c.Set("user_id", userID)
		c.Next()
	}
}

func validateToken(tokenString string) (uint, error) {
	log.Println("inside validation of token")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if x, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Fatalln("Invalid signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

		} else {
			log.Println("Parsing complete all okay. _ is: ", x)
		}
		return []byte(config.JWTSecret), nil
	})
	log.Println("The token is this: ", token)

	if err != nil {
		return 0, err
	}

	// token.Valid already checks expiration
	if !token.Valid {
		log.Fatalln("Invalid token")
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Fatalln("Invalid token claims")
		return 0, fmt.Errorf("invalid token claims")
	}

	log.Println("Claims are : ", claims)

	sub, ok := claims["sub"].(float64)
	if !ok {
		log.Fatalln("Invalid sub claim")
		return 0, fmt.Errorf("invalid sub claim")
	}

	return uint(sub), nil
}
