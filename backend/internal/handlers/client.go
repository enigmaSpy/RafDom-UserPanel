package handlers

import (
	"net/http"
	"renovation-api/internal/db"
	"renovation-api/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type CreateClientInput struct{
	Name string `json:"name" binding:"required"`
	Surname string `json:"surname"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone"`
	Password string `json:"password" binding:"required"`
	Address string `json:"address"`
	City string `json:"city"`
}

func CreateClient(c *gin.Context){
	role,_ :=c.Get("role")
	if role != "admin"{
		c.JSON(http.StatusForbidden, gin.H{"error":"brak uprawnień"})
		return
	}
	var input CreateClientInput
	if err:=c.ShouldBindJSON(&input);err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"brak wymaganych danych/błedny format"})
		return
	}
	var count int64
	db.DB.Model(&models.User{}).Where("email = ?", input.Email).Count(&count)
	if count > 0{
		c.JSON(http.StatusConflict, gin.H{"error":"klient z podanym adresem istnieje w bazie"})
		return
	}

	hashedPassword, err:= bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"błąd podczas hashowania"})
		return
	}

	newClient:=models.User{
		Name: input.Name,
		Surname: input.Surname,
		Email: input.Email,
		Phone: input.Phone,
		PasswordHash: string(hashedPassword),
		Address: input.Address,
		City: input.City,
	}

	if err:=db.DB.Create(&newClient).Error; err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Nie udało się utworzyć klienta"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Klient został utworzony",
		"client_id": newClient.ID,
	})
}