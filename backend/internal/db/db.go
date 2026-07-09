package db

import (
	"fmt"
	"log"
	"os"

	"renovation-api/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB 

func Connect(){

	err:= godotenv.Load("../.env")
	if err !=nil{
		log.Println("Brak dostępu do pliku .env")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Warsaw",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err !=nil{
		log.Fatal("Nie można połączyć z bazą \n",err)
	}
	fmt.Printf("GORM: Sukces, połączono z %s\n", os.Getenv("DB_NAME"))
	err = database.AutoMigrate(
		&models.User{},
		&models.Renovation{},
		&models.LaborTask{},
		&models.Transaction{},
		&models.ProgressUpdate{},
		&models.Message{},
	)
	if err !=nil{
		log.Fatal("Błąd migracji db: ", err)
	}
	fmt.Println("Migracja zakończona sukcesem")
	DB = database
}
func Migrate() {
	DB.AutoMigrate(&models.Renovation{}, &models.User{}, /* ... */)
	
	var renovations []models.Renovation
	DB.Where("admin_id IS NULL").Find(&renovations)
	
	var admin models.User
	if err := DB.Where("role = ?", "admin").First(&admin).Error; err == nil {
		for _, r := range renovations {
			r.AdminID = admin.ID
			DB.Save(&r)
		}
	}
}