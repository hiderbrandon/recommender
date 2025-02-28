package handlers

import (
	"net/http"
	"recommender/internal/services"
	"recommender/models"

	"github.com/gin-gonic/gin"
)

func GetStocks(c *gin.Context) {
	stocks, err := services.FetchStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve stocks"})
		return
	}
	c.JSON(http.StatusOK, stocks)
}

func PostStock(c *gin.Context) {
	var stock models.Stock
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := services.AddStock(&stock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save stock"})
		return
	}
	c.JSON(http.StatusCreated, stock)
}
