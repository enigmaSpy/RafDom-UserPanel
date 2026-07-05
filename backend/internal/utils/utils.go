package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RespondWithError(c *gin.Context, code int, message string)  {
	c.JSON(code, gin.H{"error": message})
}
func ParseUUIDParam(c *gin.Context, paramName string) (uuid.UUID, bool){
	paramString := c.Param(paramName)
	parsedUUID,err := uuid.Parse(paramString)
	if err!=nil{
		RespondWithError(c, http.StatusBadRequest, "Nieprawidłowy format parametru: oczekiwane UUID")
		return uuid.Nil, false
	}
	return parsedUUID, true
}