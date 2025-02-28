package repositories

import (
	"recommender/config"
	"recommender/models"
)

func GetAllStocks() ([]models.Stock, error) {
	var stocks []models.Stock
	result := config.DB.Find(&stocks)
	return stocks, result.Error
}

func CreateStock(stock *models.Stock) error {
	return config.DB.Create(stock).Error
}
