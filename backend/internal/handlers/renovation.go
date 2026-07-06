package handlers

import (
	"net/http"
	"time"

	"renovation-api/internal/db"
	"renovation-api/internal/models"
	"renovation-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//*------Renovation-------POST
type CreateRenovationInput struct {
	ClientID    string `json:"client_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
// CreateRenovation godoc
// @Summary Tworzy nowy projekt remontowy
// @Description Tworzy nowy remont i przypisuje go do istniejącego klienta. Wymaga roli admin.
// @Tags remonty
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body CreateRenovationInput true "Dane nowego remontu"
// @Success 201 {object} models.Renovation
// @Router /api/admin/renovations [post]
func CreateRenovation(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		utils.RespondWithError(c, http.StatusForbidden, "Brak uprawnień")
		return
	}
	var input CreateRenovationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Nieprawidłowe dane/format")
		return
	}

	clientUUID, err := uuid.Parse(input.ClientID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Nieprawidłowy identyfikator użytkownika")
		return
	}

	var count int64
	db.DB.Model(&models.User{}).Where("id = ? AND role = ?", clientUUID, "client").Count(&count)
	if count == 0 {
		utils.RespondWithError(c, http.StatusNotFound, "Użytkownik nie istnieje w bazie")
		return
	}
	newRenovation := models.Renovation{
		ClientID:    clientUUID,
		Name:        input.Name,
		Description: input.Description,
	}
	if err := db.DB.Create(&newRenovation).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Błąd przy zapisywaniu w bazie")
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":       "Nowy projekt został utworzony",
		"renovation_id": newRenovation.ID,
		"status":        newRenovation.Status,
	})
}

//*----------TaskLabor---------------------POST
type AddLaborTaskInput struct {
	Label     string  `json:"label" binding:"required"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
	Unit      string  `json:"unit" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
	Note      string  `json:"note"`
}
// AddLaborTask godoc
// @Summary Dodaje usługę roboczą do kosztorysu
// @Description Dodaje pozycję roboczą (praca) do istniejącego projektu remontowego. Wymaga roli admin.
// @Tags remonty
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID Remontu"
// @Param input body AddLaborTaskInput true "Dane usługi roboczej"
// @Success 201 {object} models.LaborTask
// @Router /api/admin/renovations/{id}/tasks [post]
func AddLaborTask(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		utils.RespondWithError(c, http.StatusForbidden, "Brak uprawnień")
		return
	}

	renovationUUID, ok := utils.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var input AddLaborTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Brak wymaganych danych")
		return
	}

	var count int64
	db.DB.Model(&models.Renovation{}).Where("id = ?", renovationUUID).Count(&count)
	if count == 0 {
		utils.RespondWithError(c, http.StatusNotFound, "Nie znaleziono projektu do powiązania")
		return
	}

	totalAmount := input.UnitPrice * input.Quantity
	newTask := models.LaborTask{
		RenovationID: renovationUUID,
		Label:        input.Label,
		UnitPrice:    input.UnitPrice,
		Unit:         input.Unit,
		Quantity:     input.Quantity,
		Amount:       totalAmount,
		Note:         input.Note,
	}
	if err := db.DB.Create(&newTask).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Nie udało się powiązać z projektem")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usługa dodana do kosztorysu",
		"task_id": newTask.ID,
		"status":  newTask.Status,
		"amount":  newTask.Amount,
	})
}

// GetRenovationDetails godoc
// @Summary Szczegóły remontu
// @Description Zwraca pełne dane remontu wraz z klientem i listą usług roboczych. Dostęp dla admina lub przypisanego klienta.
// @Tags remonty
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID Remontu"
// @Success 200 {object} models.Renovation
// @Router /api/renovations/{id} [get]
func GetRenovationDetails(c *gin.Context) {
	tokenUserID, _ := c.Get("userID")
	tokenRole, _ := c.Get("role")

	renovationUUID, ok := utils.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var renovation models.Renovation
	if err := db.DB.Preload("Client").Preload("LaborTasks").Where("id = ?", renovationUUID).First(&renovation).Error; err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Nie znaleziono projektu o tym ID")
		return
	}

	if tokenRole == "client" && renovation.ClientID != tokenUserID {
		utils.RespondWithError(c, http.StatusForbidden, "Brak uprawnień")
		return
	}
	c.JSON(http.StatusOK, renovation)
}

//*---AddTransaction----POST
type AddTransactionInput struct {
	Type   string  `json:"type" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
	Note   string  `json:"note"`
}
// AddTransaction godoc
// @Summary Rejestruje transakcję finansową
// @Description Dodaje wpłatę/zakup materiałów lub płatność za pracę do projektu. Wymaga roli admin. Dozwolone typy: material_deposit, material_expense, labour_payment.
// @Tags transakcje
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID Remontu"
// @Param input body AddTransactionInput true "Dane transakcji"
// @Success 200 {object} models.Transaction
// @Router /api/admin/renovations/{id}/transactions [post]
func AddTransaction(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		utils.RespondWithError(c, http.StatusForbidden, "Brak uprawnień")
		return
	}

	renovationUUID, ok := utils.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var input AddTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Brak wymaganych danych")
		return
	}

	validTypes := map[string]bool{"material_deposit": true, "material_expens": true, "labour_payment": true}
	if !validTypes[input.Type] {
		utils.RespondWithError(c, http.StatusBadRequest, "Nieprawidłowy typ transakcji")
		return
	}

	var count int64
	db.DB.Model(&models.Renovation{}).Where("id = ?", renovationUUID).Count(&count)
	if count == 0 {
		utils.RespondWithError(c, http.StatusNotFound, "Remont nie istnieje")
		return
	}
	newTransaction := models.Transaction{
		RenovationID: renovationUUID,
		Type:         input.Type,
		Amount:       input.Amount,
		Note:         input.Note,
		Date:         time.Now(),
	}
	if err := db.DB.Create(&newTransaction).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Wystąpił błąd przy rejestracji transakcji")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":        "Transakcja pomyślnie zarejestrowana",
		"transaction_id": newTransaction.ID,
	})
}

//*Agregation--GET
type RenovationSummary struct {
	RenovationID     uuid.UUID `json:"renovation_id"`
	Status           string    `json:"status"`
	TotalLaborCost   float64   `json:"total_labor_cost"`
	LaborPaid        float64   `json:"labor_paid"`
	LaborBalance     float64   `json:"labor_balance"`
	MaterialDeposits float64   `json:"material_deposit"`
	MaterialExpenses float64   `json:"material_expenses"`
	MaterialBalance  float64   `json:"material_balance"`
}
// GetRenovationSummary godoc
// @Summary Zwraca podsumowanie finansowe remontu
// @Description Agreguje koszty pracy, płatności, wpłaty na materiały i wydatki. Zwraca salda. Dostęp dla admina lub przypisanego klienta.
// @Tags remonty
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID Remontu"
// @Success 200 {object} models.Renovation
// @Router /api/renovations/{id}/summary [get]
func GetRenovationSummary(c *gin.Context) {
	tokenUserID, _ := c.Get("userID")
	tokenRole, _ := c.Get("role")

	renovationUUID, ok := utils.ParseUUIDParam(c, "id")
	if !ok {
		return
	}
	var renovation models.Renovation
	if err := db.DB.Where("id = ?", renovationUUID).First(&renovation).Error; err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Projekt nie istnieje")
		return
	}

	if tokenRole == "client" && renovation.ClientID != tokenUserID {
		utils.RespondWithError(c, http.StatusForbidden, "Brak dostępu")
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
type AddProgressUpdateInput struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	LaborTaskID string   `json:"labor_task_id"`
}
// AddProgressUpdate godoc
// @Summary Dodaje wpis w dzienniku prac
// @Description Rejestruje aktualizację postępu prac w projekcie. Opcjonalnie można powiązać z konkretnym zadaniem roboczym. Wymaga roli admin.
// @Tags postęp prac
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID Remontu"
// @Param input body AddProgressUpdateInput true "Dane aktualizacji postępu"
// @Success 201 {object} models.ProgressUpdate
// @Router /api/admin/renovations/{id}/progress [post]
func AddProgressUpdate(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		utils.RespondWithError(c, http.StatusForbidden, "Brak uprawnień")
		return
	}

	renovationUUID, ok := utils.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var input AddProgressUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Brak wymaganych danych wejściowych")
		return
	}
	var labourTaskUUID *uuid.UUID
	if input.LaborTaskID != "" {
		parsed, err := uuid.Parse(input.LaborTaskID)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Nieprawidłowy identyfikator zadania roboczego")
			return
		}
		labourTaskUUID = &parsed
	}
	newProgress := models.ProgressUpdate{
		RenovationID: renovationUUID,
		LaborTaskID:  labourTaskUUID,
		Title:        input.Title,
		Description:  input.Description,
		Photos:       input.Photos,
		Date:         time.Now(),
	}
	if err := db.DB.Create(&newProgress).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Nie udało się zapisać postępu prac")
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":     "Dodano wpis w dzienniku prac",
		"progress_id": newProgress.ID,
	})
}