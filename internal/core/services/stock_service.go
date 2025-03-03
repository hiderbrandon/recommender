package services

import (
	"log"

	"recommender/internal/core/domain"
	"recommender/internal/core/ports"
)

type StockService struct {
	repository ports.StockRepository
	apiClient  ports.StockAPIClient // Usa la interfaz en lugar de una implementación concreta
}

func NewStockService(repo ports.StockRepository, apiClient ports.StockAPIClient) *StockService {
	return &StockService{
		repository: repo,
		apiClient:  apiClient,
	}
}

func (s *StockService) FetchAndStoreStocks() error {
	nextPage := ""
	for {
		response, err := s.apiClient.FetchStocks(nextPage)
		if err != nil {
			return err
		}

		for _, stock := range response.Items {
			existing, _ := s.repository.GetStockByTickerAndTime(stock.Ticker, stock.Time)
			if existing.Ticker == "" {
				if err := s.repository.Create(&stock); err != nil {
					log.Println("Error saving stock:", err)
				}
			}
		}

		if response.NextPage == "" {
			break
		}
		nextPage = response.NextPage
	}
	return nil
}

func (s *StockService) FetchStocks() ([]domain.Stock, error) {
	return s.repository.GetAll()
}

func (s *StockService) AddStock(stock *domain.Stock) error {
	return s.repository.Create(stock)
}
