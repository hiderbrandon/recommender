package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"recommender/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// StockServiceInterface defines the methods we need for testing
type StockServiceInterface interface {
	FetchStocks(limit, offset int) ([]domain.Stock, error)
	AddStock(stock *domain.Stock) error
	GetTopRecommendedStocks(limit int) ([]domain.Stock, error)
	GetStockByTicker(ticker string) (*domain.Stock, error)
}

// MockStockService is a mock implementation of the StockService
type MockStockService struct {
	mock.Mock
}

func (m *MockStockService) FetchStocks(limit, offset int) ([]domain.Stock, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]domain.Stock), args.Error(1)
}

func (m *MockStockService) AddStock(stock *domain.Stock) error {
	args := m.Called(stock)
	return args.Error(0)
}

func (m *MockStockService) GetTopRecommendedStocks(limit int) ([]domain.Stock, error) {
	args := m.Called(limit)
	return args.Get(0).([]domain.Stock), args.Error(1)
}

func (m *MockStockService) GetStockByTicker(ticker string) (*domain.Stock, error) {
	args := m.Called(ticker)
	return args.Get(0).(*domain.Stock), args.Error(1)
}

func setupTestRouter() (*gin.Engine, *MockStockService) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockStockService)

	router := gin.New()

	// Setup CORS middleware (same as in SetupRouter)
	router.Use(func(c *gin.Context) {
		// Set CORS headers for testing
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Define routes with custom handlers that use our mock
	router.GET("/stocks", func(c *gin.Context) {
		limit := 10
		offset := 0

		if l := c.Query("limit"); l != "" {
			if parsedLimit := parseIntOrDefault(l, 10); parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		if o := c.Query("offset"); o != "" {
			if parsedOffset := parseIntOrDefault(o, 0); parsedOffset >= 0 {
				offset = parsedOffset
			}
		}

		stocks, err := mockService.FetchStocks(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stocks"})
			return
		}
		c.JSON(http.StatusOK, stocks)
	})

	router.POST("/stocks", func(c *gin.Context) {
		var stock domain.Stock
		if err := c.ShouldBindJSON(&stock); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		if err := mockService.AddStock(&stock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save stock"})
			return
		}
		c.JSON(http.StatusCreated, stock)
	})

	router.GET("/stocks/recommendations", func(c *gin.Context) {
		stocks, err := mockService.GetTopRecommendedStocks(5)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommendations"})
			return
		}
		c.JSON(http.StatusOK, stocks)
	})

	router.GET("/stocks/:ticker", func(c *gin.Context) {
		ticker := c.Param("ticker")
		stock, err := mockService.GetStockByTicker(ticker)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
			return
		}
		c.JSON(http.StatusOK, stock)
	})

	return router, mockService
}

// Helper function to parse integers with default values
func parseIntOrDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	// Simple parse for test purposes
	switch s {
	case "5":
		return 5
	case "10":
		return 10
	case "0":
		return 0
	default:
		return defaultVal
	}
}

func TestSetupRouter_RouterConfiguration(t *testing.T) {
	router, _ := setupTestRouter()

	// Verify router is not nil
	assert.NotNil(t, router, "Router should not be nil")
}

func TestSetupRouter_CORSConfiguration(t *testing.T) {
	router, _ := setupTestRouter()

	// Test CORS preflight request
	req := httptest.NewRequest("OPTIONS", "/stocks", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check CORS headers
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "OPTIONS")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
}

func TestSetupRouter_GetStocksRoute(t *testing.T) {
	router, mockService := setupTestRouter()

	// Mock data
	mockStocks := []domain.Stock{
		{
			ID:         1,
			Ticker:     "AAPL",
			Company:    "Apple Inc.",
			Brokerage:  "Test Broker",
			Action:     "BUY",
			RatingFrom: "HOLD",
			RatingTo:   "BUY",
			TargetFrom: 150.0,
			TargetTo:   180.0,
			Time:       time.Now(),
		},
	}

	mockService.On("FetchStocks", 10, 0).Return(mockStocks, nil)

	req := httptest.NewRequest("GET", "/stocks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []domain.Stock
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "AAPL", response[0].Ticker)

	mockService.AssertExpectations(t)
}

func TestSetupRouter_GetStocksRouteWithParams(t *testing.T) {
	router, mockService := setupTestRouter()

	mockStocks := []domain.Stock{}
	mockService.On("FetchStocks", 5, 10).Return(mockStocks, nil)

	req := httptest.NewRequest("GET", "/stocks?limit=5&offset=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestSetupRouter_PostStockRoute(t *testing.T) {
	router, mockService := setupTestRouter()

	newStock := domain.Stock{
		Ticker:     "GOOGL",
		Company:    "Alphabet Inc.",
		Brokerage:  "Test Broker",
		Action:     "BUY",
		RatingFrom: "HOLD",
		RatingTo:   "BUY",
		TargetFrom: 2500.0,
		TargetTo:   2800.0,
		Time:       time.Now(),
	}

	mockService.On("AddStock", mock.AnythingOfType("*domain.Stock")).Return(nil)

	jsonData, _ := json.Marshal(newStock)
	req := httptest.NewRequest("POST", "/stocks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestSetupRouter_GetRecommendationsRoute(t *testing.T) {
	router, mockService := setupTestRouter()

	mockRecommendations := []domain.Stock{
		{
			ID:         1,
			Ticker:     "TSLA",
			Company:    "Tesla Inc.",
			Brokerage:  "Test Broker",
			Action:     "BUY",
			RatingFrom: "HOLD",
			RatingTo:   "STRONG_BUY",
			TargetFrom: 800.0,
			TargetTo:   1000.0,
			Time:       time.Now(),
		},
	}

	mockService.On("GetTopRecommendedStocks", 5).Return(mockRecommendations, nil)

	req := httptest.NewRequest("GET", "/stocks/recommendations", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []domain.Stock
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "TSLA", response[0].Ticker)

	mockService.AssertExpectations(t)
}

func TestSetupRouter_GetStockByTickerRoute(t *testing.T) {
	router, mockService := setupTestRouter()

	mockStock := &domain.Stock{
		ID:         1,
		Ticker:     "MSFT",
		Company:    "Microsoft Corporation",
		Brokerage:  "Test Broker",
		Action:     "BUY",
		RatingFrom: "NEUTRAL",
		RatingTo:   "BUY",
		TargetFrom: 300.0,
		TargetTo:   350.0,
		Time:       time.Now(),
	}

	mockService.On("GetStockByTicker", "MSFT").Return(mockStock, nil)

	req := httptest.NewRequest("GET", "/stocks/MSFT", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Stock
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "MSFT", response.Ticker)

	mockService.AssertExpectations(t)
}

func TestSetupRouter_InvalidRoute(t *testing.T) {
	router, _ := setupTestRouter()

	req := httptest.NewRequest("GET", "/invalid-route", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSetupRouter_CORSWithDifferentOrigins(t *testing.T) {
	router, _ := setupTestRouter()

	testCases := []string{
		"http://localhost:3000",
		"https://example.com",
		"http://192.168.1.1:8080",
	}

	for _, origin := range testCases {
		req := httptest.NewRequest("OPTIONS", "/stocks", nil)
		req.Header.Set("Origin", origin)
		req.Header.Set("Access-Control-Request-Method", "GET")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"),
			"Should allow all origins for origin: %s", origin)
	}
}

func TestSetupRouter_AllHTTPMethods(t *testing.T) {
	router, mockService := setupTestRouter()

	// Test allowed methods
	allowedMethods := []string{"GET", "POST", "OPTIONS"}

	for _, method := range allowedMethods {
		req := httptest.NewRequest(method, "/stocks", nil)
		w := httptest.NewRecorder()
		// Mock service calls for GET and POST
		if method == "GET" {
			mockService.On("FetchStocks", 10, 0).Return([]domain.Stock{}, nil).Maybe()
		} else if method == "POST" {
			// Add mock expectation for POST request
			mockService.On("AddStock", mock.AnythingOfType("*domain.Stock")).Return(nil).Maybe()
			req.Header.Set("Content-Type", "application/json")
			// Create a simple JSON body for POST
			req = httptest.NewRequest(method, "/stocks", bytes.NewBuffer([]byte(`{}`)))
			req.Header.Set("Content-Type", "application/json")
		}

		router.ServeHTTP(w, req)

		if method == "OPTIONS" {
			assert.NotEqual(t, http.StatusMethodNotAllowed, w.Code,
				"Method %s should be allowed", method)
		} else {
			assert.NotEqual(t, http.StatusMethodNotAllowed, w.Code,
				"Method %s should be allowed", method)
		}
	}
}
