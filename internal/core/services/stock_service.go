package services

import (
	"log"
	"math"
	"sort"

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
	stocks, err := s.repository.GetRecentStocks(100) // Tomamos un grupo grande para filtrar mejor
	if err != nil {
		return nil, err
	}

	// Calcular la puntuación de cada stock
	type ScoredStock struct {
		Stock domain.Stock
		Score float64
	}
	var scoredStocks []ScoredStock

	for _, stock := range stocks {
		score := calculateScore(stock)
		scoredStocks = append(scoredStocks, ScoredStock{Stock: stock, Score: score})
	}

	// Ordenar por puntuación de mayor a menor
	sort.Slice(scoredStocks, func(i, j int) bool {
		return scoredStocks[i].Score > scoredStocks[j].Score
	})

	// Tomar los primeros "limit" elementos
	var topStocks []domain.Stock
	for i := 0; i < limit && i < len(scoredStocks); i++ {
		topStocks = append(topStocks, scoredStocks[i].Stock)
	}

	return topStocks, nil
}

func (s *StockService) GetStockByTicker(ticker string) (*domain.Stock, error) {
	return s.repository.GetStockByTicker(ticker)
}

// calculateScore ahora delega responsabilidades a subfunciones.
func calculateScore(stock domain.Stock) float64 {
	priceImpact := calculatePriceImpact(stock)
	ratingImpact := calculateRatingImpact(stock)
	brokerageWeight := getBrokerageWeight(stock.Brokerage)

	return (priceImpact + ratingImpact) * brokerageWeight
}

// calculatePriceImpact calcula el impacto por cambio de precio
func calculatePriceImpact(stock domain.Stock) float64 {
	percentageChange := (stock.TargetTo - stock.TargetFrom) / stock.TargetFrom * 100
	return math.Round(percentageChange*100) / 100 // Redondeo a dos decimales
}

// calculateRatingImpact evalúa el impacto por cambio de rating
func calculateRatingImpact(stock domain.Stock) float64 {
	ratingScores := map[string]float64{
		"Sell":    -2,
		"Neutral": 0,
		"Buy":     2,
	}
	fromScore := ratingScores[stock.RatingFrom]
	toScore := ratingScores[stock.RatingTo]
	return toScore - fromScore
}

// getBrokerageWeight devuelve el peso asociado a la corredora
func getBrokerageWeight(brokerage string) float64 {
	weights := map[string]float64{
		"The Goldman Sachs Group": 1.5,
		"JP Morgan":               1.4,
		"Morgan Stanley":          1.3,
		"Others":                  1.0,
	}
	weight, exists := weights[brokerage]
	if !exists {
		return weights["Others"]
	}
	return weight
}
