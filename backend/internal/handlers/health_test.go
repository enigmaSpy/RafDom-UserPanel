package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)
func TestHealthEndpoint(t *testing.T){
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/api/health", func(c *gin.Context){
		c.String(http.StatusOK, "El Psy Kongroo")
	})

	w := httptest.NewRecorder()

	req, _:=http.NewRequest(http.MethodGet, "/api/health", nil)

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK{
		t.Errorf("Oczekiwano statusu 200, otrzymano %d", w.Code)
	}
	expectedBody := "El Psy Kongroo"
	if w.Body.String()!=expectedBody{
		t.Errorf("Oczekiwano '%s', otrzymano '%s'", expectedBody, w.Body.String())
	}
}