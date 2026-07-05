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
//*Agregation--GET
type RenovationSummary struct{
	RenovationID uuid.UUID `json:"renovation_id"`
	Status string `json:"status"`
	TotalLaborCost float64 `json:"total_labor_cost"`
	LaborPaid float64 `json:"labor_paid"`
	LaborBalance float64 `json:"labor_balance"`
	MaterialDeposits float64 `json:"material_deposit"`
	MaterialExpenses float64 `json:"material_expenses"`
	MaterialBalance float64 `json:"material_balance"`
}

func GetRenovationSummary(c *gin.Context){
	tokenUserID, _ := c.Get("userID")
	tokenRole, _ := c.Get("role")

	renovationIDParam := c.Param("id")
	renovationUUID, err:=uuid.Parse(renovationIDParam)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe ID remontu"})
		return
	}
	var renovation models.Renovation
	if err:= db.DB.Where("id = ?", renovationUUID).First(&renovation).Error; err !=nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "Projekt nie istnieje"})
		return
	}

	if tokenRole == "client" && renovation.ClientID != tokenUserID{
		c.JSON(http.StatusForbidden, gin.H{"error": "Brak dostępu"})
		return
	}
	var summary RenovationSummary
	summary.RenovationID = renovation.ID
	summary.Status = renovation.Status
	
	//?LaborCostAggregation
	db.DB.Model(&models.LaborTask{}). 
			Select("COALESCE(SUM(amount),0)").
			Where("renovation_id = ?", renovationUUID).
			Scan(&summary.TotalLaborCost) 

	//?LaborPayment
	db.DB.Model(&models.Transaction{}). 
			Select("COALESCE(SUM(amount),0)"). 
			Where("renovation_id = ? AND type = ?", renovationUUID, "labor_payment"). 
			Scan(&summary.LaborPaid)

	//?MaterialDeposit
	db.DB.Model(&models.Transaction{}). 
			Select("COALESCE(SUM(amount),0)"). 
			Where("renovation_id = ? AND type = ?", renovationUUID, "material_deposit"). 
			Scan(&summary.MaterialDeposits)

	//?MaterialExpense
	db.DB.Model(&models.Transaction{}). 
			Select("COALESCE(SUM(amount),0)"). 
			Where("renovation_id = ? AND type = ?", renovationUUID, "material_expense").
			Scan(&summary.MaterialExpenses)

	//?balance
	summary.LaborBalance = summary.TotalLaborCost - summary.LaborPaid
	summary.MaterialBalance = summary.MaterialDeposits - summary.MaterialExpenses

	c.JSON(http.StatusOK, summary)
}

//*---ProgressUpdate
type AddProgressUpdateInput struct{
	Title string `json:"title" binding:"required"`
	Description string `json:"description"`
	Photos []string `json:"photos"`
	LaborTaskID string `json:"labor_task_id"`
}
func AddProgressUpdate(c *gin.Context){
	role,_ :=c.Get("role")
	if role != "admin"{
		c.JSON(http.StatusForbidden, gin.H{"error":"Brak uprawnień"})
		return
	}

	renovationIDParam := c.Param("id")
	renovationUUID, err := uuid.Parse(renovationIDParam)
	if err !=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe ID remontu"})
		return
	}

	var input AddProgressUpdateInput
	if err:=c.ShouldBindJSON(&input); err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brak wymaganych danych wejściowych"})
		return
	}
	var labourTaskUUID *uuid.UUID
	if input.LaborTaskID !=""{
		parsed, err:= uuid.Parse(input.LaborTaskID)
		if err!=nil{
			labourTaskUUID = &parsed
		}
	}
	newProgress :=models.ProgressUpdate{
		RenovationID: renovationUUID,
		LaborTaskID: labourTaskUUID,
		Title: input.Title,
		Description: input.Description,
		Photos: input.Photos,
		Date: time.Now(),
	}
	if err:=db.DB.Create(&newProgress).Error; err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nie udało się zapisać postępu prac"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":     "Dodano wpis w dzienniku prac",
		"progress_id": newProgress.ID,
	})
}
