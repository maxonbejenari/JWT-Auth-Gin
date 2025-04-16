package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/maxonbejenari/gin-auth/config"
	"github.com/maxonbejenari/gin-auth/database"
	"github.com/maxonbejenari/gin-auth/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

func Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	// extract data form request body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	// hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to hash the password",
		})
	}

	// if all is good we create new user
	newUser := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already in use",
		})
		return
	}

	result := database.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "User creation failed",
		})
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": newUser.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	log.Println("Jwt:", config.AppConfig.JWTSecret)

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Error with generating token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "Success",
		"token": tokenString,
	})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
	}

	// find the user in db
	var user models.User
	err := database.DB.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Register first",
		})
	}

	// compare the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"err": "Invalid email or password",
		})
	}

	// generate the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Error with generating token",
		})
	}

	c.JSON(http.StatusAccepted, gin.H{
		"msg":   "Welcome Back",
		"token": tokenString,
	})
}
