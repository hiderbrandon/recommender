package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"recommender/internal/core/domain"
	"recommender/internal/core/services"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

//implementa todos los métodos requeridos por services.StockService
type fakeStockRepository struct{}

func (f *fakeStockRepository) GetAll(limit, offset int) ([]domain.Stock, error) {
	return []domain.Stock{
		{
			ID:         1,
			Ticker:     "AAPL",
			Company:    "Apple Inc.",
			Brokerage:  "JP Morgan",
			Action:     "target raised by",
			RatingFrom: "Neutral",
			RatingTo:   "Buy",
			TargetFrom: 150,
			TargetTo:   180,
			Time:       time.Now(),
		},
	}, nil
}

func (f *fakeStockRepository) Create(stock *domain.Stock) error {
	return nil
}

func (f *fakeStockRepository) GetStockByTickerAndTime(ticker string, t time.Time) (*domain.Stock, error) {
	return nil, nil
}

func (f *fakeStockRepository) GetTopStocksByTarget(limit int) ([]domain.Stock, error) {
	return nil, nil
}

func (f *fakeStockRepository) GetStockByTicker(ticker string) (*domain.Stock, error) {
	return nil, nil
}

func (f *fakeStockRepository) GetRecentStocks(limit int) ([]domain.Stock, error) {
	return nil, nil
}

type fakeStockAPIClient struct{}

func (f *fakeStockAPIClient) FetchStocks(nextPage string) (*domain.APIResponse, error) {
	return nil, nil
}

func TestGetStocks_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear servicio con dependencias fake
	repo := &fakeStockRepository{}
	client := &fakeStockAPIClient{}
	service := services.NewStockService(repo, client)

	handler := NewStockHandler(service)

	// Configurar ruta
	router := gin.New()
	router.GET("/stocks", handler.GetStocks)

	// Ejecutar petición simulada
	req, _ := http.NewRequest("GET", "/stocks?limit=1&offset=0", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verificar respuesta
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "AAPL")
}

func TestPostStock_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear servicio con dependencias fake
	repo := &fakeStockRepository{}
	client := &fakeStockAPIClient{}
	service := services.NewStockService(repo, client)

	handler := NewStockHandler(service)

	// Configurar ruta
	router := gin.New()
	router.POST("/stocks", handler.PostStock)

	// Crear JSON válido para el stock
	jsonData := `{
		"ticker": "GOOGL",
		"company": "Alphabet Inc.",
		"brokerage": "Goldman Sachs",
		"action": "target raised by",
		"rating_from": "Hold",
		"rating_to": "Buy",
		"target_from": 2500.0,
		"target_to": 2800.0
	}`

	// Ejecutar petición simulada
	req, _ := http.NewRequest("POST", "/stocks", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verificar respuesta
	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "GOOGL")
	assert.Contains(t, resp.Body.String(), "Alphabet Inc.")
}

func TestPostStock_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear servicio con dependencias fake
	repo := &fakeStockRepository{}
	client := &fakeStockAPIClient{}
	service := services.NewStockService(repo, client)

	handler := NewStockHandler(service)

	// Configurar ruta
	router := gin.New()
	router.POST("/stocks", handler.PostStock)

	// JSON inválido
	invalidJSON := `{"ticker": "INVALID"` // JSON malformado

	// Ejecutar petición simulada
	req, _ := http.NewRequest("POST", "/stocks", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verificar respuesta de error
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid request")
}

func TestGetRecommendations_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear servicio con dependencias fake mejoradas
	repo := &fakeStockRepositoryWithRecommendations{}
	client := &fakeStockAPIClient{}
	service := services.NewStockService(repo, client)

	handler := NewStockHandler(service)

	// Configurar ruta
	router := gin.New()
	router.GET("/stocks/recommendations", handler.GetRecommendations)

	// Ejecutar petición simulada
	req, _ := http.NewRequest("GET", "/stocks/recommendations", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verificar respuesta
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "TSLA")
	assert.Contains(t, resp.Body.String(), "NVDA")
}

func TestGetStockByTicker_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear servicio con dependencias fake mejoradas
	repo := &fakeStockRepositoryWithTicker{}
	client := &fakeStockAPIClient{}
	service := services.NewStockService(repo, client)

	handler := NewStockHandler(service)

	// Configurar ruta
	router := gin.New()
	router.GET("/stocks/:ticker", handler.GetStockByTicker)

	// Ejecutar petición simulada
	req, _ := http.NewRequest("GET", "/stocks/MSFT", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verificar respuesta
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "MSFT")
	assert.Contains(t, resp.Body.String(), "Microsoft")
}

func TestGetStockByTicker_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Crear servicio con dependencias fake que retorna error
	repo := &fakeStockRepositoryNotFound{}
	client := &fakeStockAPIClient{}
	service := services.NewStockService(repo, client)

	handler := NewStockHandler(service)

	// Configurar ruta
	router := gin.New()
	router.GET("/stocks/:ticker", handler.GetStockByTicker)

	// Ejecutar petición simulada con ticker inexistente
	req, _ := http.NewRequest("GET", "/stocks/NOTFOUND", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verificar respuesta de error
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Contains(t, resp.Body.String(), "Stock not found")
}

// Repositorio fake mejorado para recomendaciones
type fakeStockRepositoryWithRecommendations struct{}

func (f *fakeStockRepositoryWithRecommendations) GetAll(limit, offset int) ([]domain.Stock, error) {
	return []domain.Stock{}, nil
}

func (f *fakeStockRepositoryWithRecommendations) Create(stock *domain.Stock) error {
	return nil
}

func (f *fakeStockRepositoryWithRecommendations) GetStockByTickerAndTime(ticker string, t time.Time) (*domain.Stock, error) {
	return nil, nil
}

func (f *fakeStockRepositoryWithRecommendations) GetTopStocksByTarget(limit int) ([]domain.Stock, error) {
	return []domain.Stock{
		{
			ID:         1,
			Ticker:     "TSLA",
			Company:    "Tesla Inc.",
			Brokerage:  "Morgan Stanley",
			Action:     "target raised by",
			RatingFrom: "Hold",
			RatingTo:   "Strong Buy",
			TargetFrom: 800,
			TargetTo:   1000,
			Time:       time.Now(),
		},
		{
			ID:         2,
			Ticker:     "NVDA",
			Company:    "NVIDIA Corporation",
			Brokerage:  "Bank of America",
			Action:     "target raised by",
			RatingFrom: "Buy",
			RatingTo:   "Strong Buy",
			TargetFrom: 450,
			TargetTo:   500,
			Time:       time.Now(),
		},
	}, nil
}

func (f *fakeStockRepositoryWithRecommendations) GetStockByTicker(ticker string) (*domain.Stock, error) {
	return nil, nil
}

func (f *fakeStockRepositoryWithRecommendations) GetRecentStocks(limit int) ([]domain.Stock, error) {
	return []domain.Stock{
		{
			ID:         1,
			Ticker:     "TSLA",
			Company:    "Tesla Inc.",
			Brokerage:  "Morgan Stanley",
			Action:     "target raised by",
			RatingFrom: "Hold",
			RatingTo:   "Strong Buy",
			TargetFrom: 800,
			TargetTo:   1000,
			Time:       time.Now(),
		},
		{
			ID:         2,
			Ticker:     "NVDA",
			Company:    "NVIDIA Corporation",
			Brokerage:  "Bank of America",
			Action:     "target raised by",
			RatingFrom: "Buy",
			RatingTo:   "Strong Buy",
			TargetFrom: 450,
			TargetTo:   500,
			Time:       time.Now(),
		},
	}, nil
}

// Repositorio fake para búsqueda por ticker
type fakeStockRepositoryWithTicker struct{}

func (f *fakeStockRepositoryWithTicker) GetAll(limit, offset int) ([]domain.Stock, error) {
	return []domain.Stock{}, nil
}

func (f *fakeStockRepositoryWithTicker) Create(stock *domain.Stock) error {
	return nil
}

func (f *fakeStockRepositoryWithTicker) GetStockByTickerAndTime(ticker string, t time.Time) (*domain.Stock, error) {
	return nil, nil
}

func (f *fakeStockRepositoryWithTicker) GetTopStocksByTarget(limit int) ([]domain.Stock, error) {
	return nil, nil
}

func (f *fakeStockRepositoryWithTicker) GetStockByTicker(ticker string) (*domain.Stock, error) {
	if ticker == "MSFT" {
		return &domain.Stock{
			ID:         3,
			Ticker:     "MSFT",
			Company:    "Microsoft Corporation",
			Brokerage:  "JP Morgan",
			Action:     "target raised by",
			RatingFrom: "Neutral",
			RatingTo:   "Buy",
			TargetFrom: 300,
			TargetTo:   350,
			Time:       time.Now(),
		}, nil
	}
	return nil, errors.New("stock not found")
}

func (f *fakeStockRepositoryWithTicker) GetRecentStocks(limit int) ([]domain.Stock, error) {
	return nil, nil
}

// Repositorio fake que simula stock no encontrado
type fakeStockRepositoryNotFound struct{}

func (f *fakeStockRepositoryNotFound) GetAll(limit, offset int) ([]domain.Stock, error) {
	return []domain.Stock{}, nil
}

func (f *fakeStockRepositoryNotFound) Create(stock *domain.Stock) error {
	return nil
}

func (f *fakeStockRepositoryNotFound) GetStockByTickerAndTime(ticker string, t time.Time) (*domain.Stock, error) {
	return nil, nil
}

func (f *fakeStockRepositoryNotFound) GetTopStocksByTarget(limit int) ([]domain.Stock, error) {
	return nil, nil
}

func (f *fakeStockRepositoryNotFound) GetStockByTicker(ticker string) (*domain.Stock, error) {
	return nil, errors.New("stock not found")
}

func (f *fakeStockRepositoryNotFound) GetRecentStocks(limit int) ([]domain.Stock, error) {
	return nil, nil
}