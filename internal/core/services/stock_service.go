package services

import (
	"sort"
	"math"
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

func calculateScore(stock domain.Stock) float64 {
	// 1. Cambio en el precio objetivo
	targetFrom := stock.TargetFrom
	targetTo := stock.TargetTo
	targetChange := (targetTo - targetFrom) / targetFrom

	// 2. Cambio en la calificación
	ratingMap := map[string]int{
		"Sell":   -1,
		"Neutral": 0,
		"Buy":    1,
	}

	ratingFrom, okFrom := ratingMap[stock.RatingFrom]
	ratingTo, okTo := ratingMap[stock.RatingTo]

	ratingChange := 0
	if okFrom && okTo {
		ratingChange = ratingTo - ratingFrom
	}

	ratingScore := 0.0
	switch ratingChange {
	case 2:
		ratingScore = 3  // Sell -> Buy
	case 1:
		ratingScore = 2  // Neutral -> Buy
	case -1:
		ratingScore = -2 // Buy -> Neutral
	case -2:
		ratingScore = -3 // Buy -> Sell
	}

	// 3. Peso del bróker
	brokerWeights := map[string]float64{
		"The Goldman Sachs Group": 1.5,
		"JP Morgan":               1.4,
	}
	brokerWeight, exists := brokerWeights[stock.Brokerage]
	if !exists {
		brokerWeight = 1.0
	}

	// 4. Recencia de la recomendación
	now := time.Now().UTC()
	daysOld := now.Sub(stock.Time).Hours() / 24
	timeWeight := math.Max(0.5, 1-(daysOld/30)) // Reduce el peso si la recomendación es vieja

	// Fórmula de puntuación
	score := (targetChange*10 + ratingScore) * brokerWeight * timeWeight
	return score
}