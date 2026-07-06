package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func APILogger() gin.HandlerFunc{
	file, err:=os.OpenFile("api_stats.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil{
		log.Fatal("Nie udało się otworzyć ppliku logów: ", err)
	}

	return func(c *gin.Context){
		startTime := time.Now()
		c.Next()

		latency:=time.Since(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		logEntry := fmt.Sprintf("[%s] %s | %3d | %13v | %15s | %s\n",
				time.Now().Format("2006-01-02 15:05:05"),
				method,
				statusCode,
				latency,
				clientIP,
				path,
		)
		if _, err:= file.WriteString(logEntry); err!=nil{
			log.Print("Błąd zapisu logu: ", err)
		}
		fmt.Print(logEntry)
	}
}