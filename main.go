package main

import (
	"github.com/gin-gonic/gin"
	"github.com/maxonbejenari/gin-auth/config"
	"github.com/maxonbejenari/gin-auth/controllers"
	"github.com/maxonbejenari/gin-auth/database"
)

func main() {
	router := gin.Default()

	// Load vars from .env
	config.LoadConfig()
	// Connect to DB
	database.ConnectDB()

	// all endPoints
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	router.Run(":" + config.AppConfig.PORT)
}
