package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"recommender/internal/core/domain"
	"recommender/internal/core/ports"
	"strconv"
	"strings"
	"time"
)

type ExternalStockAPI struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewExternalStockAPI() ports.StockAPIClient {
	apiURL := strings.TrimSpace(os.Getenv("API_URL"))
	if apiURL == "" {
		panic("API_URL environment variable not set")
	}
	apiKey := strings.TrimSpace(os.Getenv("API_KEY"))
	if apiKey == "" {
		panic("API_KEY environment variable not set")
	}
	return &ExternalStockAPI{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: apiURL,
		apiKey:  apiKey,
	}
}

func parsePrice(priceStr string) (float64, error) {
	//elimina caracteres no necesarios 
	priceStr = strings.TrimSpace(strings.Replace(priceStr, "$", "", -1))
	priceStr = strings.Replace(priceStr, ",", "", -1)

	if priceStr == "" {
		return 0, fmt.Errorf("parsePrice: recibido string vacío")
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		fmt.Printf("⚠ Error convirtiendo TargetFrom: '%s'\n", priceStr)
		return 0, fmt.Errorf("parsePrice: error convirtiendo '%s' a float64", priceStr)
	}

	return price, nil
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr) // Convertir a `time.Time`
}

func (a *ExternalStockAPI) FetchStocks(nextPage string) (*domain.APIResponse, error) {
	url := a.baseURL
	if nextPage != "" {
		url = fmt.Sprintf("%s?next_page=%s", a.baseURL, nextPage)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Agregar encabezados, incluyendo User-Agent similar al de curl
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "curl/7.68.0") // Opcional, para imitar la petición de curl

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var apiResponseDTO domain.APIResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&apiResponseDTO)
	if err != nil {
		return nil, err
	}
	// Convertir `StockDTO` a `Stock`
	var stocks []domain.Stock
	for _, stockDTO := range apiResponseDTO.Items {
		targetFrom, err := parsePrice(stockDTO.TargetFrom)
		if err != nil {
			return nil, fmt.Errorf("error parsing TargetFrom")
		}

		targetTo, err := parsePrice(stockDTO.TargetTo)
		if err != nil {
			return nil, fmt.Errorf("error parsing TargetTo")
		}

		parsedTime, err := parseTime(stockDTO.Time)
		if err != nil {
			return nil, fmt.Errorf("error parsing Time")
		}

		stocks = append(stocks, domain.Stock{
			Ticker:     stockDTO.Ticker,
			TargetFrom: targetFrom,
			TargetTo:   targetTo,
			Company:    stockDTO.Company,
			Brokerage:  stockDTO.Brokerage,
			Action:     stockDTO.Action,
			RatingFrom: stockDTO.RatingFrom,
			RatingTo:   stockDTO.RatingTo,
			Time:       parsedTime,
		})
	}

	return &domain.APIResponse{
		Items:    stocks,
		NextPage: apiResponseDTO.NextPage,
	}, nil
}
