package main

import (
	"auth/internal/entity"
	"auth/internal/handler"
	"auth/internal/repository"
	"auth/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// AutoMigrate Book entity
	if err := db.AutoMigrate(&entity.Book{}, &entity.User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Database successful migrate")

	// Dependency injection
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	router := gin.Default()
	handler.SetupBookRoutes(router, bookHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
