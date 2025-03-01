package services

import (
    "recommender/internal/core/domain"
    "recommender/internal/core/ports"
)

type StockService struct {
    repository port.StockRepository
}

func NewStockService(repo port.StockRepository) *StockService {
    return &StockService{repository: repo}
}

func (s *StockService) FetchStocks() ([]domain.Stock, error) {
    return s.repository.GetAll()
}

func (s *StockService) AddStock(stock *domain.Stock) error {
    return s.repository.Create(stock)
}
