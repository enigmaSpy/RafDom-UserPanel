package main

import (
	"fmt"
	"log"
	"os"

	"renovation-api/internal/db"
	"renovation-api/internal/handlers"
	"renovation-api/internal/middleware"
	"renovation-api/internal/seed"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()
	seed.SeedAdmin(db.DB)

	if err := os.MkdirAll("uploads", os.ModePerm); err!=nil{
		log.Fatal("Nie udało się utworzyć folderu uploads: ", err)
	}

	router := gin.Default()
	router.Static("/uploads", "./uploads")
	//*---Default-Router
	api := router.Group("/api")
	{
		api.POST("/login", handlers.Login)
		api.GET("/ws/chat", handlers.ConnectChatWS)
		api.GET("/health", func(c *gin.Context) {
			c.String(200, "El Psy Kongroo")
		})
	}
	//*----Admin-Router
	adminOnly := api.Group("/admin")
	adminOnly.Use(middleware.AuthRequired())
	{
		adminOnly.GET("/me", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			role, _ := c.Get("role")

			c.JSON(200, gin.H{
				"message":    "Autoryzacja udana",
				"twoje_id":   userID,
				"twoja_rola": role,
			})
		})
		adminOnly.POST("/clients", handlers.CreateClient)
		adminOnly.POST("/renovations", handlers.CreateRenovation)
		adminOnly.POST("/renovations/:id/tasks", handlers.AddLaborTask)
		adminOnly.POST("/renovations/:id/transactions", handlers.AddTransaction)
		adminOnly.POST("/upload", handlers.UploadFile)
		adminOnly.POST("/renovations/:id/progress", handlers.AddProgressUpdate)
	}
	//*----Authenticated-Router
	authenticated := api.Group("/")
	authenticated.Use(middleware.AuthRequired())
	{
		authenticated.GET("/renovations/:id", handlers.GetRenovationDetails)
		authenticated.GET("/renovations/:id/summary", handlers.GetRenovationSummary)
	}
	//*----Server
	fmt.Println("Serwer działa na porcie 8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal("Błąd serwera: ", err)
	}

}
