package rdb

import (
	"fmt"
	"log"
	"os"

	"github.com/dopeCape/schduler/internal/models"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

func NewDB() error {
	dbURl := os.Getenv("TURSO_DATABASE_URL")
	dbAuthToken := os.Getenv("TURSO_AUTH_TOKEN")
	url := fmt.Sprint(dbURl, "?authToken=", dbAuthToken)
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "libsql",
		DSN:        url,
	}), &gorm.Config{})
	// here add auto migrate code

	db.AutoMigrate(&models.ApiKey{})

	db.AutoMigrate(&models.Task{})

	db.AutoMigrate(&models.Execution{})

	if err != nil {
		log.Fatalf("db connection failur: %v", err)
	}
	fmt.Println("Connected to db..")
	return nil

}

func GetDb() (*gorm.DB, error) {
	dbURl := os.Getenv("TURSO_DATABASE_URL")
	dbAuthToken := os.Getenv("TURSO_AUTH_TOKEN")

	url := fmt.Sprint(dbURl, "?authToken=", dbAuthToken)
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "libsql",
		DSN:        url,
	}), &gorm.Config{})

	// here add auto migrate code

	if err != nil {
		log.Fatalf("db connection failur: %v", err)
	}
	fmt.Println("Connected to db..")
	return db, nil
}
