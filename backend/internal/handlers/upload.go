package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context){
	err := c.Request.ParseMultipartForm(5<<20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plik jest zbyt duży lub format żądania jest niepoprawny"+err.Error()})
		return
	}
	file, err := c.FormFile("file")
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Brak pliku w rządzniu"})
		return 
	}
	ext := strings.ToLower((filepath.Ext(file.Filename)))
	allowedExt := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".pdf": true}

	if !allowedExt[ext]{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Niedozwolony format pliku"})
		return
	}

	newFilename := uuid.NewString()+ext
	savePath := filepath.Join("uploads", newFilename) 

	if err:=c.SaveUploadedFile(file,savePath); err !=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Nie udało się zapisać pliku na serwerze"})
		return
	}

	fileURL := fmt.Sprintf("/uploads/%s", newFilename)
	c.JSON(http.StatusOK, gin.H{
		"message": "Plik wgrany pomyślnie",
		"url":     fileURL,
	})
}