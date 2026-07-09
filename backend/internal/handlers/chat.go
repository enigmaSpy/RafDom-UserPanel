package handlers

import (
	"log"
	"net/http"
	"sync"

	"renovation-api/internal/auth"
	"renovation-api/internal/db"
	"renovation-api/internal/models"
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

type Room struct {
	connections map[uuid.UUID]*websocket.Conn
	mu          sync.RWMutex
}

var rooms = make(map[uuid.UUID]*Room) 
var roomsMu sync.RWMutex

type WSMsgInput struct {
	Content string `json:"content"`
}

// ConnectChatWS godoc
// @Summary WebSocket czatu per remont
// @Description Połączenie WebSocket do czatu konkretnego remontu. Wymaga tokena w query param. Backend sam wyznacza odbiorcę.
// @Tags chat
// @Param renovation_id path string true "UUID Remontu"
// @Param token query string true "JWT Token"
// @Router /api/ws/chat/{renovation_id} [get]
func ConnectChatWS(c *gin.Context) {
	
	
	tokenString := c.Query("token")
	if tokenString == "" {
		log.Println("=== WS ERROR: brak tokenu")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Brak tokenu w URL"})
		return
	}

	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Nieważny token")
		return
	}

	renovationUUID, ok := utils.ParseUUIDParam(c, "renovation_id")
	if !ok {
		return
	}
	var renovation models.Renovation
	if err := db.DB.Where("id = ?", renovationUUID).First(&renovation).Error; err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Remont nie istnieje")
		return
	}
	if claims.Role == "client" && renovation.ClientID != claims.UserID {
		utils.RespondWithError(c, http.StatusForbidden, "Brak dostępu do tego remontu")
		return
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Błąd podczas Upgrade do WebSockets:", err)
		return
	}
	defer ws.Close()

	roomsMu.Lock()
	room, exists := rooms[renovationUUID]
	if !exists {
		room = &Room{connections: make(map[uuid.UUID]*websocket.Conn)}
		rooms[renovationUUID] = room
	}
	roomsMu.Unlock()

	room.mu.Lock()
	room.connections[claims.UserID] = ws
	room.mu.Unlock()

	log.Printf("User %s joined room %s", claims.UserID, renovationUUID)

	defer func() {
		room.mu.Lock()
		delete(room.connections, claims.UserID)
		room.mu.Unlock()

		roomsMu.Lock()
		if len(room.connections) == 0 {
			delete(rooms, renovationUUID)
		}
		roomsMu.Unlock()

		log.Printf("User %s left room %s", claims.UserID, renovationUUID)
	}()

	for {
		var input WSMsgInput

		if err := ws.ReadJSON(&input); err != nil {
			log.Printf("Błąd odczytu WS dla user %s: %v", claims.UserID, err)
			break
		}

		if input.Content == "" {
			ws.WriteJSON(gin.H{"error": "Pusta wiadomość"})
			continue
		}

		var receiverID uuid.UUID
		if claims.Role == "admin" {
			receiverID = renovation.ClientID
		} else {
			receiverID = renovation.AdminID
		}

		newMsg := models.Message{
			RenovationID: renovationUUID,
			SenderID:     claims.UserID,
			ReceiverID:   receiverID,
			Content:      input.Content,
		}
		if err := db.DB.Create(&newMsg).Error; err != nil {
			log.Printf("Błąd zapisu wiadomości: %v", err)
			ws.WriteJSON(gin.H{"error": "Błąd zapisu wiadomości"})
			continue
		}

		room.mu.RLock()
		receiverConn, isOnline := room.connections[receiverID]
		room.mu.RUnlock()

		if isOnline {
			if err := receiverConn.WriteJSON(newMsg); err != nil {
				log.Printf("Nie udało się dostarczyć do %s: %v", receiverID, err)
			}
		}

		ws.WriteJSON(gin.H{
			"status":  "sent",
			"message": newMsg,
		})
	}
}

// GetChatHistory godoc
// @Summary Historia czatu
// @Description Zwraca wszystkie wiadomości dla danego remontu. Dostęp dla admina lub przypisanego klienta.
// @Tags chat
// @Security BearerAuth
// @Param id path string true "UUID Remontu"
// @Success 200 {array} models.Message
// @Failure 403 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Router /api/renovations/{id}/messages [get]
func GetChatHistory(c *gin.Context) {
	tokenUserID, _ := c.Get("userID")
	tokenRole, _ := c.Get("role")

	renovationUUID, ok := utils.ParseUUIDParam(c, "id")
	if !ok {
		return
	}

	var renovation models.Renovation
	if err := db.DB.Where("id = ?", renovationUUID).First(&renovation).Error; err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Remont nie istnieje")
		return
	}

	if tokenRole == "client" && renovation.ClientID != tokenUserID.(uuid.UUID) {
		utils.RespondWithError(c, http.StatusForbidden, "Brak uprawnień")
		return
	}

	var messages []models.Message
	if err := db.DB.Where("renovation_id = ?", renovationUUID).Order("created_at asc").Find(&messages).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Błąd ładowania historii czatu")
		return
	}

	c.JSON(http.StatusOK, messages)
}