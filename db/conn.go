package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
func InitDB() {
	var err error

  envError := godotenv.Load()
	if envError != nil {
		panic("cannot load env")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})


	if err != nil {
		panic(err)
	}
}

