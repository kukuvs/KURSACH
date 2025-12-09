package database

import (
	"buch/src"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Автомиграция - GORM сам создаст таблицы на основе структур!
	log.Println("Running auto migrations...")
	err = db.AutoMigrate(
		&src.Category{},
		&src.Transaction{},
		&src.Account{},
		&src.Transfer{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	DB = db
	return nil
}

// Close больше не нужен так явно для GORM, он управляет пулом сам, но можно оставить пустышку для совместимости main.go
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
