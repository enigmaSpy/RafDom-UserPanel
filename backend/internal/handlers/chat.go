package handlers

import (
	"log"
	"net/http"
	"sync"

	"renovation-api/internal/auth"
	"renovation-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var activeClients = make(map[uuid.UUID]*websocket.Conn)
var clientsMu sync.RWMutex

func ConnectChatWS(c *gin.Context){
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Brak tokenu w URL"})
		return
	}

	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Nieważny token")
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Błąd podczas Upgrade do WebSockets:", err)
		return
	}
	defer ws.Close() 

	clientsMu.Lock()
	activeClients[claims.UserID] = ws
	clientsMu.Unlock()
	//!Do usunięcia po teście
	log.Printf("Użytkownik %s połączył się z czatem", claims.UserID)

	defer func(){
		clientsMu.Lock()
		delete(activeClients, claims.UserID)
		clientsMu.Unlock()
		log.Printf("Użytkownik %s rozłączył się", claims.UserID)
	}()

	for{
		var incomingMsg map[string]interface{}
		err:=ws.ReadJSON(&incomingMsg)
		if err !=nil{
			break
		}
		log.Printf("Otrzymano wiadomość: %v", incomingMsg)

		ws.WriteJSON(gin.H{
			"status": "serwer potwierdza odbiór",
			"your_message": incomingMsg,
		})
	}
}