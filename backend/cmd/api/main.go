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

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "renovation-api/docs"
)
//@title Renovation API
//@version 1.0
//@description System zarządzania projektami remontowymi 
//@host localhost:8081
//@BasePath /api
//@securityDefinition.apikey BearerAuth
//@in header
//@name Authorization
func main() {


	db.Connect()
	seed.SeedAdmin(db.DB)

	if err := os.MkdirAll("uploads", os.ModePerm); err!=nil{
		log.Fatal("Nie udało się utworzyć folderu uploads: ", err)
	}

	router := gin.Default()
	router.Use(middleware.APILogger())
	router.Static("/uploads", "./uploads")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
		adminOnly.POST("/client/create", handlers.CreateClient)
		adminOnly.GET("/client/list", handlers.GetListClient)

		adminOnly.POST("/renovations", handlers.CreateRenovation)
		adminOnly.POST("/renovations/:id/tasks", handlers.AddLaborTask)
		adminOnly.POST("/renovations/:id/transactions", handlers.AddTransaction)
		adminOnly.POST("/renovations/:id/progress", handlers.AddProgressUpdate)
		adminOnly.POST("/renovations/:id/update", handlers.UpdateRenovation)
		adminOnly.PUT("/renovations/:id", handlers.UpdateRenovation)
		adminOnly.DELETE("/renovations/:id", handlers.DeleteRenovation)

		adminOnly.PUT("/tasks/:id", handlers.UpdateLoborTask)
		adminOnly.DELETE("/tasks/:id", handlers.DeleteLaborTask)
		adminOnly.POST("/upload", handlers.UploadFile)
	}
	//*----Authenticated-Router
	authenticated := api.Group("/")
	authenticated.Use(middleware.AuthRequired())
	{
		authenticated.GET("/renovations/:id", handlers.GetRenovationDetails)
		authenticated.GET("/renovations/:id/summary", handlers.GetRenovationSummary)
		authenticated.GET("/renovations/:id/messages", handlers.GetChatHistory)
		authenticated.GET("/renovations/list", handlers.GetListRenovation)

		authenticated.GET("/client/:id", handlers.GetClient)
		authenticated.POST("/client/:id/update", handlers.UpdateClient)
	}
	//*----Server
	fmt.Println("Serwer HTTPS działa na porcie 8081")
	if err:=router.RunTLS(":8081", "cert.pem", "key.pem"); err !=nil{
		log.Fatal("Błąd serwera HTTPS: ", err)
	}
	

}
