package services

import (
	"log"

	"recommender/internal/core/domain"
	"recommender/internal/core/ports"

	"gorm.io/gorm"
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
	log.Println("📥 Iniciando importación de datos desde la API externa...")

	nextPage := ""
	for {
		// Obtener datos de la API externa
		apiResponse, err := s.apiClient.FetchStocks(nextPage)
		if err != nil {
			return err
		}

		for _, stock := range apiResponse.Items {
			// Verificar si el stock ya existe en la base de datos
			existingStock, err := s.repository.GetStockByTickerAndTime(stock.Ticker, stock.Time)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			// Si no existe, lo insertamos
			if existingStock == nil {
				err = s.repository.Create(&stock)
				if err != nil {
					log.Printf("⚠ Error insertando stock %s: %v\n", stock.Ticker, err)
				} else {
					log.Printf("✅ Stock insertado: %s\n", stock.Ticker)
				}
			} else {
				log.Printf("ℹ Stock %s ya existe en la base de datos, ignorando...\n", stock.Ticker)
			}
		}

		// Si no hay más páginas, terminamos
		if apiResponse.NextPage == "" {
			break
		}
		nextPage = apiResponse.NextPage
	}

	log.Println("✅ Importación completada.")
	return nil
}

func (s *StockService) FetchStocks() ([]domain.Stock, error) {
	return s.repository.GetAll()
}

func (s *StockService) AddStock(stock *domain.Stock) error {
	return s.repository.Create(stock)
}

func (s *StockService) GetTopRecommendedStocks(limit int) ([]domain.Stock, error) {
    return s.repository.GetTopStocksByTarget(limit)
}
