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
	"renovation-api/internal/db"
	"renovation-api/internal/models"
)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var activeClients = make(map[uuid.UUID]*websocket.Conn)
var clientsMu sync.RWMutex

type WSMsgInput struct{
	RenovationID string `json:"renovation_id"`
	ReceiverID 	 string `json:"receiver_id"`
	Content 	 string `json:"content"`
}
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
		var input WSMsgInput

		if err:=ws.ReadJSON(&input); err !=nil{
			break
		}
		log.Printf("Przetwarzanie wiadomości: %v", input.Content)

		renovationUUID, err1 := uuid.Parse(input.RenovationID)
		receiverUUID, err2 := uuid.Parse(input.ReceiverID)
		if err1 != nil || err2 !=nil{
			ws.WriteJSON(gin.H{"error":"nieprawidłowy format id"})
			continue
		}

		newMsg := models.Message{
			RenovationID: renovationUUID,
			SenderID: claims.UserID,
			ReceiverID: receiverUUID,
			Content: input.Content,
		}
		if err := db.DB.Create(&newMsg).Error; err !=nil{
			ws.WriteJSON(gin.H{"error":"Błąd zapisu danych"})
			continue
		}
		clientsMu.RLock()
		receiverConn, isOnline := activeClients[receiverUUID]
		clientsMu.RUnlock()

		if isOnline{
			err = receiverConn.WriteJSON(newMsg)
			if err !=nil{
				log.Println("Nie udało się dostarczyć wiadomości: ", err)
			}
		}
		ws.WriteJSON(gin.H{
			"status":"sent",
			"message": newMsg,
		})
	}
}


func GetChatHistory(c *gin.Context){
	tokenUserID, _:=c.Get("userID")
	tokenRole, _:=c.Get("role")

	renovationUUID, ok :=utils.ParseUUIDParam(c, "id")
	if !ok{
		return
	}
	if tokenRole == "client"{
		var renovation models.Renovation
		if err := db.DB.First(&renovation, "id =?", renovationUUID).Error; err!=nil{
			utils.RespondWithError(c, http.StatusNotFound, "Remont nie istnieje")
			return
		}
		if renovation.ClientID.String() != tokenUserID.(string){
			utils.RespondWithError(c, http.StatusForbidden, "Brak uprawnień")
			return
		}
	}
	var messages []models.Message
	if err:=db.DB.Where("renovation_id = ?", renovationUUID).Order("created_at asc").Find(&messages).Error; err !=nil{
		utils.RespondWithError(c, http.StatusInternalServerError, "Błąd ładowania histroii czatu")
		return
	} 
	c.JSON(http.StatusOK, messages)
}