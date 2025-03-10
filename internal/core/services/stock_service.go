package services

import (
	"log"

	"recommender/internal/core/domain"
	"recommender/internal/core/ports"

	"time"

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
		var err error
		var apiResponse *domain.APIResponse

		// Reintentar hasta 3 veces en caso de fallo
		for attempts := 1; attempts <= 3; attempts++ {
			apiResponse, err = s.apiClient.FetchStocks(nextPage)
			if err == nil {
				break
			}

			log.Printf("⚠ Error en FetchStocks (Intento %d/3): %v", attempts, err)
			time.Sleep(time.Duration(attempts) * time.Second)
		}

		if err != nil {
			log.Printf("❌ Falló la importación de stocks en la página '%s' después de 3 intentos. Continuando con la siguiente...", nextPage)
			nextPage = "" // Forzar fin del bucle si hay fallo total
			continue
		}

		// Procesar los datos obtenidos
		for _, stock := range apiResponse.Items {
			existingStock, err := s.repository.GetStockByTickerAndTime(stock.Ticker, stock.Time)
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Printf("⚠ Error verificando existencia de %s: %v", stock.Ticker, err)
				continue
			}

			if existingStock == nil {
				err = s.repository.Create(&stock)
				if err != nil {
					log.Printf("⚠ Error insertando stock %s: %v", stock.Ticker, err)
				} else {
					log.Printf("✅ Stock insertado: %s", stock.Ticker)
				}
			} else {
				log.Printf("ℹ Stock %s ya existe en la base de datos, ignorando...", stock.Ticker)
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

func (s *StockService) FetchStocks(limit, offset int) ([]domain.Stock, error) {
	return s.repository.GetAll(limit, offset)
}

func (s *StockService) AddStock(stock *domain.Stock) error {
	return s.repository.Create(stock)
}

func (s *StockService) GetTopRecommendedStocks(limit int) ([]domain.Stock, error) {
	return s.repository.GetTopStocksByTarget(limit)
}

func (s *StockService) GetStockByTicker(ticker string) (*domain.Stock, error) {
	return s.repository.GetStockByTicker(ticker)
}
