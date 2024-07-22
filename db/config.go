package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := os.Getenv("PINTAS_DB_DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	db.AutoMigrate(&Url{}, &User{}, &UserToken{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	return db
}
