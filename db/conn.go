package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"moapick/db/models"
)

var Client *gorm.DB

func InitDB() {
	var err error

	envError := godotenv.Load()
	if envError != nil {
		panic("cannot load env")
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// Connect to the PostgreSQL server without specifying the database name
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable", host, user, password, port))
	if err != nil {
		panic(err)
	}
	// Create the database if it does not exist
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbname))
	if err != nil {
		// Ignore the error if the database already exists
		if !strings.Contains(err.Error(), "already exists") {
			panic(err)
		}
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable", host, user, password, port, dbname)
	Client, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	Client.AutoMigrate(&models.User{})
}
