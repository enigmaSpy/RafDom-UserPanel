package middleware

import (
	"net/http"
	"strings"
	"renovation-api/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc{
	return func(c *gin.Context){
		authHeader := c.GetHeader("Authorization")
		if authHeader == ""{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"brak nagłówka"})
			return 
		}

		parts := strings.Split(authHeader, " ")
		if len(parts)!=2 ||parts[0]!="Bearer"{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"nieprawidłowy format nagłówka"})
			return 
		}

		tokenString := parts[1]
		claims, err:=auth.ValidateToken(tokenString)
		if err != nil{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"Nieprawidłowy lub wygasły token"})
			return 
		}
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}