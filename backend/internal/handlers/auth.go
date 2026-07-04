package handlers

import (
	"net/http"
	"renovation-api/internal/auth"
	"renovation-api/internal/db"
	"renovation-api/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct{
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context){
	var input LoginInput

	if err:= c.ShouldBindJSON(&input); err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"nieprawidłowe dane wejściowe"})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err!=nil{
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Nieprawidłowy adres email"})
		return
	}

	if err:= bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err!=nil{
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Nieprawidłowe hasło"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Błąd serwera przy tworzeniu tokena"})
	}
	c.JSON(http.StatusOK, gin.H{
		"token":token,
		"role":user.Role,
	})
}