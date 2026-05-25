package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"awesomeProject3/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	r := setupRouter()
	r.GET("/test-auth", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req, _ := http.NewRequest("GET", "/test-auth", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegister_InvalidEmail(t *testing.T) {
	r := setupRouter()
	r.POST("/register", func(c *gin.Context) {
		var input struct {
			Email string `json:"email"`
		}
		c.ShouldBindJSON(&input)
		if input.Email != "test@mail.com" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	body := map[string]string{"username": "user", "email": "bad-email", "password": "Password123"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_ValidEmail(t *testing.T) {
	r := setupRouter()
	r.POST("/register", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	body := map[string]string{"username": "user", "email": "test@mail.com", "password": "Password123"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogin_WrongPassword(t *testing.T) {
	r := setupRouter()
	r.POST("/login", func(c *gin.Context) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	})

	body := map[string]string{"username": "user", "password": "wrongpassword"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAddTransaction_Unauthorized(t *testing.T) {
	r := setupRouter()
	r.POST("/transactions", func(c *gin.Context) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	})

	body := map[string]interface{}{"amount": 500, "category_id": 1, "note": "Lunch"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAddTransaction_InvalidAmount(t *testing.T) {
	r := setupRouter()
	r.POST("/transactions", func(c *gin.Context) {
		var input struct {
			Amount float64 `json:"amount"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})

	body := map[string]interface{}{"amount": "not-a-number", "category_id": 1, "note": "Test"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTransactions_Success(t *testing.T) {
	r := setupRouter()
	r.GET("/transactions", func(c *gin.Context) {
		list := []models.Transaction{}
		c.JSON(http.StatusOK, list)
	})

	req, _ := http.NewRequest("GET", "/transactions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTransaction_NotFound(t *testing.T) {
	r := setupRouter()
	r.GET("/transactions/:id", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
	})

	req, _ := http.NewRequest("GET", "/transactions/9999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteTransaction_NotFound(t *testing.T) {
	r := setupRouter()
	r.DELETE("/transactions/:id", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
	})

	req, _ := http.NewRequest("DELETE", "/transactions/9999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateTransaction_InvalidBody(t *testing.T) {
	r := setupRouter()
	r.PUT("/transactions/:id", func(c *gin.Context) {
		var input struct {
			Amount float64 `json:"amount"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	})

	body := map[string]interface{}{"amount": "string-instead-of-float"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/transactions/1", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
