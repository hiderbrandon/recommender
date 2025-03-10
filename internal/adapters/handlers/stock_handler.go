package handlers

import (
	"net/http"
	"recommender/internal/core/domain"
	"recommender/internal/core/services"

	"strconv"
	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	service *services.StockService
}

func NewStockHandler(service *services.StockService) *StockHandler {
	return &StockHandler{service: service}
}

func (h *StockHandler) GetStocks(c *gin.Context) {
	// Leer parámetros de la query
	limit := 10 // Valor por defecto
	offset := 0 // Valor por defecto

	if l, exists := c.GetQuery("limit"); exists {
		parsedLimit, err := strconv.Atoi(l)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o, exists := c.GetQuery("offset"); exists {
		parsedOffset, err := strconv.Atoi(o)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	stocks, err := h.service.FetchStocks(limit, offset)
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

func (h *StockHandler) GetStockByTicker(c *gin.Context) {
	ticker := c.Param("ticker") // Obtener el ticker de la URL

	stock, err := h.service.GetStockByTicker(ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	c.JSON(http.StatusOK, stock)
}
