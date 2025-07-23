package handler

import (
	"auth/internal/auth"
	"auth/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupBookRoutes sets up only the book routes
func SetupBookRoutes(router *gin.Engine, bookHandler *BookHandler, userHandler *UserHandler) {
	router.POST("/books", func(c *gin.Context) {
		var book model.BookRes
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		msg, err := bookHandler.AddBook(book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": msg, "book": book})
	})
	router.GET("/books", func(c *gin.Context) {
		data, err := bookHandler.GetAllBook()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"data": data})
	})

	router.PUT("/books/:id", func(c *gin.Context) {
		id := c.Param("id")

		var dataToUpdate model.ChangeData
		if err := c.ShouldBindJSON(&dataToUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if id != "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "id is required"})
			return
		}

		result, err := bookHandler.EditBook(id, dataToUpdate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"result": result})
	})

	router.POST("/login", func(c *gin.Context) {
		var dataToLogin auth.LoginUser

		if err := c.ShouldBindJSON(&dataToLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := userHandler.Login(dataToLogin)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"result": result})
	})

	router.POST("/register", func(c *gin.Context) {
		var dataToRegister auth.RegisterUser
		if err := c.ShouldBindJSON(&dataToRegister); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		result, err := userHandler.Register(dataToRegister)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"result": result})
	})
}
