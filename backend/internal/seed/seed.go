package seed 

import (
	"fmt"
	"log"
	"renovation-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB){
	var count int64

	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

	if count==0{
		hashedPassword, err:= bcrypt.GenerateFromPassword([]byte("root"), bcrypt.DefaultCost)
		if err !=nil{
			log.Fatal("Błąd hashowania: ", err)
		}
		admin := models.User{
			Name: "Kamil",
			Email: "admin@root.io",
			PasswordHash: string(hashedPassword),
			Role: "admin",
		}

		if err:=db.Create(&admin).Error; err!=nil{
			log.Fatal("Nie udało sie zapisać admina: ",err)
		}
		fmt.Println("Utworzono konto-> admin@root.io | root")
	}else{
		fmt.Println("Admin istnie w bazie")
	}
}