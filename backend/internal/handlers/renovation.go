package handlers

import (
	"net/http"
	"time"

	"renovation-api/internal/db"
	"renovation-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//*------Renovation-------POST
type CreateRenovationInput struct{
	ClientID string`json:"client_id" binding:"required"`
	Name string`json:"name" binding:"required"`
	Description string `json:"description"`
}

func CreateRenovation(c *gin.Context){
	role, _:=c.Get("role")
	if role!="admin"{
		c.JSON(http.StatusForbidden, gin.H{"error":"Brak uprawnień"})
		return
	}
	var input CreateRenovationInput
	if err:=c.ShouldBindJSON(&input); err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Nieprawidłowe dane/format"})
		return
	}

	clientUUID, err:=uuid.Parse(input.ClientID)
	if err !=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Nieprawidłowy identyfikator użytkownika"})
		return
	}

	var count int64
	db.DB.Model(&models.User{}).Where("id = ? AND role = ?", clientUUID, "client").Count(&count)
	if count == 0{
		c.JSON(http.StatusNotFound, gin.H{"error":"Użytkownik nie istnieje w bazie"})
		return 
	}
	newRenovation := models.Renovation{
		ClientID: clientUUID,
		Name: input.Name,
		Description: input.Description,
	}
	if err:=db.DB.Create(&newRenovation).Error; err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Błąd przy zapisywaniu w bazie"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Nowy projekt został utworzony",
		"renovation_id": newRenovation.ID,
		"status": newRenovation.Status,
	})
}

//*----------TaskLabor---------------------POST
type AddLaborTaskInput struct{
	Label string `json:"label" binding:"required"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
	Unit string `json:"unit" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required"`
	Note string `json:"note"`
}

func AddLaborTask(c *gin.Context){
	role,_:=c.Get("role")
	if role !="admin"{
		c.JSON(http.StatusForbidden, gin.H{"error":"Brak uprawnień"})
		return
	}

	renovationIDParam := c.Param("id")
	renovationUUID, err := uuid.Parse(renovationIDParam)
	if err !=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Nieprawidłowy format ID"})
		return
	}

	var input AddLaborTaskInput
	if err:=c.ShouldBindJSON(&input); err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Brak wymaganych danych"})
		return
	}

	var count int64
	db.DB.Model(&models.Renovation{}).Where("id = ?", renovationUUID).Count(&count)
	if count == 0{
		c.JSON(http.StatusNotFound, gin.H{"error":"Nie znaleziono projektu do powiązania"})
		return
	}

	totalAmount := input.UnitPrice*input.Quantity
	newTask := models.LaborTask{
		RenovationID: renovationUUID,
		Label: input.Label,
		UnitPrice: input.UnitPrice,
		Unit: input.Unit,
		Quantity: input.Quantity,
		Amount: totalAmount,
		Note: input.Note,
	}
	if err := db.DB.Create(&newTask).Error; err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Nie udało się powiązać z projektem"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usługa dodana do kosztorysu",
		"task_id": newTask.ID,
		"status": newTask.Status,
		"amount": newTask.Amount,
	})
}

//*---RenovationDetails-----GET
func GetRenovationDetails(c *gin.Context){
	tokenUserID, _:=c.Get("userID")
	tokenRole,_:=c.Get("role")

	renovationIDParam := c.Param("id")
	renovationUUID, err := uuid.Parse(renovationIDParam)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Nieprawidłowy identyfikator"})
		return
	}

	var renovation models.Renovation
	if err := db.DB.Preload("Client").Preload("LaborTasks").Where("id = ?", renovationUUID).First(&renovation).Error; err !=nil{
		c.JSON(http.StatusNotFound, gin.H{"error":"Nie znaleziono projektu o tym ID","id":renovationUUID})
		return
	}

	if tokenRole == "client" && renovation.ClientID != tokenUserID{
		c.JSON(http.StatusForbidden, gin.H{"error":"Brak uprawnień"})
		return
	}
	c.JSON(http.StatusOK, renovation)
}
//*---AddTransaction----POST
type AddTransactionInput struct{
	Type string `json:"type" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
	Note string `json:"note"`
}
func AddTransaction(c *gin.Context){
	role, _:= c.Get("role")
	if role !="admin"{
		c.JSON(http.StatusForbidden, gin.H{"error":"Brak uprawnień"})
		return 
	}

	renovationIDParam :=c.Param("id")
	renovationUUID, err:=uuid.Parse(renovationIDParam)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Nieprawidłowy format ID"})
		return
	}

	var input AddTransactionInput
	if err :=c.ShouldBindJSON(&input);err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Brak wymaganych danych"})
		return
	}

	validTypes := map[string]bool{"material_deposit": true, "material_expens":true, "labour_payment":true}
	if !validTypes[input.Type]{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowy typ transakcji"})
		return
	}

	var count int64
	db.DB.Model(&models.Renovation{}).Where("id = ?", renovationUUID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Remont nie istnieje"})
		return
	}
	newTransaction:=models.Transaction{
		RenovationID: renovationUUID,
		Type: input.Type,
		Amount: input.Amount,
		Note: input.Note,
		Date: time.Now(),
	}
	if err:=db.DB.Create(&newTransaction).Error; err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wystąpił błąd przy rejestracji transakcji"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":        "Transakcja pomyślnie zarejestrowana",
		"transaction_id": newTransaction.ID,
	})
}