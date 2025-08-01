package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	dsn := os.Getenv("DATABASE_DSN")

	if dsn == "" {
		dsn = "host=db user=wrink-dev password=wrink-dev-123 dbname=wrink-dev port=5432 sslmode=disable"
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to load database: ", err)
	}

	log.Println("Connected to database")
}
