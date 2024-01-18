package db

import (
	"database/sql"
	"fmt"
	"moapick/db/models"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Client *gorm.DB

func InitDB() {
	var err error

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable", host, user, password, port))
	
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbname))
	
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			panic(err)
		}
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable", host, user, password, port, dbname)
	Client, err = gorm.Open(postgres.Open(dsn), &gorm.Config{

	})

	Client.Logger.LogMode(3)

	if err != nil {
		panic(err)
	}

	Client.AutoMigrate(&models.User{}, &models.Tag{}, &models.Article{})
}

