package handlers

import (
	"net/http"
	"net/http/httptest"
	"recommender/internal/core/domain"
	"recommender/internal/core/services"
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