package handlers

import (
	"net/http"
	"renovation-api/internal/db"
	"renovation-api/internal/models"

	"github.com/gin-gonic/gin"
)
type DashboardStats struct{
	TotalProjects     int64   `json:"total_projects"`
	ActiveProjects    int64   `json:"active_projects"`
	CompletedProjects int64   `json:"completed_projects"`
	TotalClients      int64   `json:"total_clients"`
	TotalIncome       float64 `json:"total_income"`
}
func GetAdminDashboard(c *gin.Context){
	var stats DashboardStats
	db.DB.Model(&models.Renovation{}).Count(&stats.TotalProjects)
	db.DB.Model(&models.Renovation{}).Where("status = ?", "in_progress").Count(&stats.ActiveProjects)
	db.DB.Model(&models.Renovation{}).Where("status = ?", "completed").Count(&stats.CompletedProjects)
	
	db.DB.Model(&models.User{}).Where("role = ?", "client").Count(&stats.TotalClients)

	db.DB.Model(&models.Transaction{}).
		Where("type = ?", "labor_payment").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.TotalIncome)
	c.JSON(http.StatusOK, stats)
}