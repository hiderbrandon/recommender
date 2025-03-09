package handlers

import (
	"net/http"
	"recommender/internal/core/domain"
	"recommender/internal/core/services"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	service *services.StockService
}

func NewStockHandler(service *services.StockService) *StockHandler {
	return &StockHandler{service: service}
}

func (h *StockHandler) GetStocks(c *gin.Context) {
	stocks, err := h.service.FetchStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stocks"})
		return
	}
	c.JSON(http.StatusOK, stocks)
}

func (h *StockHandler) PostStock(c *gin.Context) {
	var stock domain.Stock
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := h.service.AddStock(&stock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save stock"})
		return
	}
	c.JSON(http.StatusCreated, stock)
}
func (h *StockHandler) GetRecommendations(c *gin.Context) {
	limit := 5 // Número de acciones recomendadas
	stocks, err := h.service.GetTopRecommendedStocks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommendations"})
		return
	}
	c.JSON(http.StatusOK, stocks)
}
