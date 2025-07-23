package main

import (
	"auth/internal/entity"
	"auth/internal/handler"
	"auth/internal/repository"
	"auth/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sevlyar/go-daemon"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {

	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	log.Println("Database successful migrate")
	if err != nil {
		log.Fatal(err)
	}
	// AutoMigrate Book entity
	if err := db.AutoMigrate(&entity.Book{}, &entity.User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Dependency injection
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(*authRepo)
	authHandler := handler.NewUserHandler(*authService)

	router := gin.Default()
	handler.SetupBookRoutes(router, bookHandler, authHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runAsDaemon() {
	cntxt := &daemon.Context{
		PidFileName: "sample.pid",
		PidFilePerm: 0644,
		LogFileName: "sample.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-daemon sample]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	log.Print("Daemon started")
}
