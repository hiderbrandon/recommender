package repository

import (
	"recommender/internal/core/domain"
	port "recommender/internal/core/ports"

	"time"

	"gorm.io/gorm"
)

type CockroachStockRepository struct {
	db *gorm.DB
}

func NewCockroachStockRepository(db *gorm.DB) port.StockRepository {
	return &CockroachStockRepository{db: db}
}

func (r *CockroachStockRepository) GetAll() ([]domain.Stock, error) {
	var stocks []domain.Stock
	result := r.db.Find(&stocks)
	return stocks, result.Error
}

func (r *CockroachStockRepository) Create(stock *domain.Stock) error {
	return r.db.Create(stock).Error
}

func (r *CockroachStockRepository) GetStockByTickerAndTime(ticker string, t time.Time) (*domain.Stock, error) {
	var stock domain.Stock
	result := r.db.Where("ticker = ? AND time = ?", ticker, t).First(&stock)
	if result.Error != nil {
		return nil, result.Error
	}
	return &stock, nil
}

func (r *CockroachStockRepository) GetTopStocksByTarget(limit int) ([]domain.Stock, error) {
	var stocks []domain.Stock
	result := r.db.Order("target_to DESC").Limit(limit).Find(&stocks)
	return stocks, result.Error
}

func (r *CockroachStockRepository) GetStockByTicker(ticker string) (*domain.Stock, error) {
    var stock domain.Stock
    result := r.db.Where("ticker = ?", ticker).First(&stock)
    if result.Error != nil {
        return nil, result.Error
    }
    return &stock, nil
}
