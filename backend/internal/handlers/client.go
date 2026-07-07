package handlers

import (
	"net/http"
	"renovation-api/internal/db"
	"renovation-api/internal/models"
	"renovation-api/internal/utils"
	"time"

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

func GetListClient(c *gin.Context){
	role, _ := c.Get("role")
	if role != "admin"{
		utils.RespondWithError(c, http.StatusUnauthorized, "Brak dostępu")
		return
	}
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil{
		utils.RespondWithError(c, http.StatusInternalServerError, "Błąd przy pobieraniu użytkowników")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func GetClient(c *gin.Context){
	role, _ := c.Get("role")
	tokenUserID, _ := c.Get("userID")
	userUUID, ok := utils.ParseUUIDParam(c, "id")
	if !ok{
		utils.RespondWithError(c, http.StatusBadRequest, "Break użytkownika o podanym id")
		return
	}
	if role == "client" && tokenUserID != userUUID{
		utils.RespondWithError(c, http.StatusForbidden, "Brak dostępu")
		return
	}
	var user models.User
	if err := db.DB.Where("id = ?", userUUID).First(&user).Error; err!=nil{
		utils.RespondWithError(c, http.StatusBadRequest, "Użytkownik nie istnieje")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})

}

type UpdateUserInput struct{
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Password string    `json:"password"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	PostalCode   string    `json:"postal_code"`
}
func UpdateClient(c *gin.Context){
	role, _ := c.Get("role")
	tokenUserID, _ := c.Get("userID")
	userUUID, ok := utils.ParseUUIDParam(c, "id")
	if !ok{
		return
	}
	
	if role != "admin" && tokenUserID.(string) !=userUUID.String(){
		utils.RespondWithError(c, http.StatusForbidden, "Brak dostępu")
		return
	}
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err !=nil{
		utils.RespondWithError(c, http.StatusBadRequest, "Nieprawidłowe dane")
		return
	}
	var user models.User
	if err:=db.DB.First(&user, "id = ?", userUUID).Error; err !=nil{
		utils.RespondWithError(c, http.StatusNotFound, "Nie znaleziono użytkownika")
		return 
	}
	if input.Email != user.Email{
		var existingUser models.User
		if err := db.DB.Where("email = ? AND id = ?", input.Email, userUUID).First(&existingUser).Error; err == nil{
			utils.RespondWithError(c, http.StatusConflict, "Email już zajęty")
			return
		}
	}

	updates := map[string]interface{}{
        "name":        input.Name,
        "surname":     input.Surname,
        "email":       input.Email,
        "phone":       input.Phone,
        "address":     input.Address,
        "city":        input.City,
        "postal_code": input.PostalCode,
        "updated_at":  time.Now(),
    }
	if input.Password !=""{
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err !=nil{
			utils.RespondWithError(c, http.StatusInternalServerError, "Błąd podczas hashowania")
			return
		}
		updates["password_hash"]=string(hashedPassword)
	}

	
	if err := db.DB.Model(&user).Omit("ID", "Role").Updates(updates).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Błąd podczas aktualizacji")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message":"Pomyślnie zaktualizowano dane użytkwnika"})
}

func DeleteClient(c *gin.Context){
	role, _ := c.Get("role")
	if role != "admin"{
		utils.RespondWithError(c, http.StatusUnauthorized, "Brak dostępu")
		return
	}

}