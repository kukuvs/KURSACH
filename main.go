package main

import (
	"buch/src/api"
	"buch/src/database"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
}

func main() {
	// Инициализация подключения к БД
	if err := database.Init(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.Close()

	// Настройка маршрутов (роутера)
	router := api.SetupRoutes()

	// Получение порта из переменных окружения
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)

	// Запуск сервера
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
