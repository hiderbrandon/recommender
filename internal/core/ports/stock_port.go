package ports

import (
	"recommender/internal/core/domain"
	"time"
)

type StockRepository interface {
	GetAll(limit, offset int) ([]domain.Stock, error)
	Create(stock *domain.Stock) error
	GetStockByTickerAndTime(ticker string, t time.Time) (*domain.Stock, error)
	GetTopStocksByTarget(limit int) ([]domain.Stock, error)
	GetStockByTicker(ticker string) (*domain.Stock, error) // 🔹 Agregar esta línea

}
