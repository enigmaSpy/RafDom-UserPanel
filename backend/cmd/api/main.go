package main

import (
	"fmt"
	"log"

	"renovation-api/internal/db"
	"renovation-api/internal/handlers"
	"renovation-api/internal/seed"
	"renovation-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main(){
	db.Connect()

	seed.SeedAdmin(db.DB)

	router:=gin.Default()
	api:=router.Group("/api")
	{
		api.POST("/login", handlers.Login)

		api.GET("/health", func(c *gin.Context){
			c.String(200, "El Psy Kongroo")
		})
	}

	protected:=api.Group("/admin")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/me", func(c *gin.Context){
			userID, _ := c.Get("userID")
			role,_ := c.Get("role")

			c.JSON(200, gin.H{
				"message": "Autoryzacja udana",
				"twoje_id": userID,
				"twoja_rola": role,
			})
		})

		protected.POST("/clients", handlers.CreateClient)
		protected.POST("/renovations", handlers.CreateRenovation)
		protected.POST("/renovations/:id/tasks", handlers.AddLaborTask)
	}

	fmt.Println("Serwer działa na porcie 8081")
	if err:=router.Run(":8081"); err!=nil{
		log.Fatal("Błąd serwera: ", err)
	}

}