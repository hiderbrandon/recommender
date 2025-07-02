package clients

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"recommender/internal/core/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePrice_Success(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "Basic price",
			input:    "100.50",
			expected: 100.50,
		},
		{
			name:     "Price with dollar sign",
			input:    "$1,234.56",
			expected: 1234.56,
		},
		{
			name:     "Price with commas",
			input:    "1,000",
			expected: 1000,
		},
		{
			name:     "Price with spaces",
			input:    " 250.75 ",
			expected: 250.75,
		},
		{
			name:     "Zero price",
			input:    "0",
			expected: 0,
		},
		{
			name:     "Complex price",
			input:    " $2,500.99 ",
			expected: 2500.99,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parsePrice(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParsePrice_Error(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "Empty string",
			input: "",
		},
		{
			name:  "Only spaces",
			input: "   ",
		},
		{
			name:  "Invalid format",
			input: "abc",
		},
		{
			name:  "Mixed invalid characters",
			input: "12.34.56",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parsePrice(tc.input)
			assert.Error(t, err)
		})
	}
}

func TestParseTime_Success(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "Valid RFC3339 time",
			input:    "2023-12-25T10:30:00Z",
			expected: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		},
		{
			name:     "Valid RFC3339 time with timezone",
			input:    "2023-06-15T14:45:30-05:00",
			expected: time.Date(2023, 6, 15, 19, 45, 30, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseTime(tc.input)
			assert.NoError(t, err)
			assert.True(t, tc.expected.Equal(result))
		})
	}
}

func TestParseTime_Error(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "Invalid format",
			input: "2023-12-25",
		},
		{
			name:  "Empty string",
			input: "",
		},
		{
			name:  "Wrong format",
			input: "Dec 25, 2023",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseTime(tc.input)
			assert.Error(t, err)
		})
	}
}

func TestNewExternalStockAPI_Success(t *testing.T) {
	// Set environment variables
	os.Setenv("API_URL", "https://api.example.com")
	os.Setenv("API_KEY", "test-api-key")
	defer func() {
		os.Unsetenv("API_URL")
		os.Unsetenv("API_KEY")
	}()

	client := NewExternalStockAPI()
	assert.NotNil(t, client)

	// Verify it's the correct type
	externalAPI, ok := client.(*ExternalStockAPI)
	assert.True(t, ok)
	assert.Equal(t, "https://api.example.com", externalAPI.baseURL)
	assert.Equal(t, "test-api-key", externalAPI.apiKey)
	assert.NotNil(t, externalAPI.client)
}

func TestNewExternalStockAPI_PanicOnMissingAPIURL(t *testing.T) {
	// Ensure API_URL is not set
	os.Unsetenv("API_URL")
	os.Setenv("API_KEY", "test-key")
	defer os.Unsetenv("API_KEY")

	assert.Panics(t, func() {
		NewExternalStockAPI()
	})
}

func TestNewExternalStockAPI_PanicOnMissingAPIKey(t *testing.T) {
	// Ensure API_KEY is not set
	os.Setenv("API_URL", "https://api.example.com")
	os.Unsetenv("API_KEY")
	defer os.Unsetenv("API_URL")

	assert.Panics(t, func() {
		NewExternalStockAPI()
	})
}

func TestNewExternalStockAPI_HandlesWhitespace(t *testing.T) {
	// Set environment variables with whitespace
	os.Setenv("API_URL", "  https://api.example.com  ")
	os.Setenv("API_KEY", "  test-api-key  ")
	defer func() {
		os.Unsetenv("API_URL")
		os.Unsetenv("API_KEY")
	}()

	client := NewExternalStockAPI()
	externalAPI, ok := client.(*ExternalStockAPI)
	require.True(t, ok)

	assert.Equal(t, "https://api.example.com", externalAPI.baseURL)
	assert.Equal(t, "test-api-key", externalAPI.apiKey)
}

func TestFetchStocks_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer test-api-key")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "curl/7.68.0", r.Header.Get("User-Agent"))

		// Mock response
		mockResponse := domain.APIResponseDTO{
			Items: []domain.StockDTO{
				{
					Ticker:     "AAPL",
					Company:    "Apple Inc.",
					Brokerage:  "JP Morgan",
					Action:     "target raised by",
					RatingFrom: "Hold",
					RatingTo:   "Buy",
					TargetFrom: "$150.00",
					TargetTo:   "$180.50",
					Time:       "2023-12-25T10:30:00Z",
				},
				{
					Ticker:     "GOOGL",
					Company:    "Alphabet Inc.",
					Brokerage:  "Goldman Sachs",
					Action:     "upgraded by",
					RatingFrom: "Neutral",
					RatingTo:   "Strong Buy",
					TargetFrom: "$2,500.00",
					TargetTo:   "$2,800.75",
					Time:       "2023-12-26T15:45:00Z",
				},
			},
			NextPage: "page2",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Set up environment
	os.Setenv("API_URL", server.URL)
	os.Setenv("API_KEY", "test-api-key")
	defer func() {
		os.Unsetenv("API_URL")
		os.Unsetenv("API_KEY")
	}()

	client := NewExternalStockAPI()
	response, err := client.FetchStocks("")

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Items, 2)
	assert.Equal(t, "page2", response.NextPage)

	// Verify first stock
	stock1 := response.Items[0]
	assert.Equal(t, "AAPL", stock1.Ticker)
	assert.Equal(t, "Apple Inc.", stock1.Company)
	assert.Equal(t, "JP Morgan", stock1.Brokerage)
	assert.Equal(t, "target raised by", stock1.Action)
	assert.Equal(t, "Hold", stock1.RatingFrom)
	assert.Equal(t, "Buy", stock1.RatingTo)
	assert.Equal(t, 150.00, stock1.TargetFrom)
	assert.Equal(t, 180.50, stock1.TargetTo)
	expectedTime1 := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	assert.True(t, expectedTime1.Equal(stock1.Time))

	// Verify second stock
	stock2 := response.Items[1]
	assert.Equal(t, "GOOGL", stock2.Ticker)
	assert.Equal(t, 2500.00, stock2.TargetFrom)
	assert.Equal(t, 2800.75, stock2.TargetTo)
}

func TestFetchStocks_WithNextPage(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify next_page parameter
		nextPage := r.URL.Query().Get("next_page")
		assert.Equal(t, "page2", nextPage)

		mockResponse := domain.APIResponseDTO{
			Items:    []domain.StockDTO{},
			NextPage: "",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Set up environment
	os.Setenv("API_URL", server.URL)
	os.Setenv("API_KEY", "test-api-key")
	defer func() {
		os.Unsetenv("API_URL")
		os.Unsetenv("API_KEY")
	}()

	client := NewExternalStockAPI()
	response, err := client.FetchStocks("page2")

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestFetchStocks_HTTPError(t *testing.T) {
	// Create mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Set up environment
	os.Setenv("API_URL", server.URL)
	os.Setenv("API_KEY", "test-api-key")
	defer func() {
		os.Unsetenv("API_URL")
		os.Unsetenv("API_KEY")
	}()

	client := NewExternalStockAPI()
	response, err := client.FetchStocks("")

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "API returned status: 500")
}

func TestFetchStocks_InvalidJSON(t *testing.T) {
	// Create mock server with invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	// Set up environment
	os.Setenv("API_URL", server.URL)
	os.Setenv("API_KEY", "test-api-key")
	defer func() {
		os.Unsetenv("API_URL")
		os.Unsetenv("API_KEY")
	}()

	client := NewExternalStockAPI()
	response, err := client.FetchStocks("")

	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestFetchStocks_InvalidPriceFormat(t *testing.T) {
	// Create mock server with invalid price format
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockResponse := domain.APIResponseDTO{
			Items: []domain.StockDTO{
				{
					Ticker:     "AAPL",
					Company:    "Apple Inc.",
					Brokerage:  "JP Morgan",
					Action:     "target raised by",
					RatingFrom: "Hold",
					RatingTo:   "Buy",
					TargetFrom: "invalid-price", // Invalid price format
					TargetTo:   "$180.50",
					Time:       "2023-12-25T10:30:00Z",
				},
			},
			NextPage: "",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Set up environment
	os.Setenv("API_URL", server.URL)
	os.Setenv("API_KEY", "test-api-key")
	defer func() {
		os.Unsetenv("API_URL")
		os.Unsetenv("API_KEY")
	}()

	client := NewExternalStockAPI()
	response, err := client.FetchStocks("")

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "error parsing TargetFrom")
}

func TestFetchStocks_InvalidTimeFormat(t *testing.T) {
	// Create mock server with invalid time format
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockResponse := domain.APIResponseDTO{
			Items: []domain.StockDTO{
				{
					Ticker:     "AAPL",
					Company:    "Apple Inc.",
					Brokerage:  "JP Morgan",
					Action:     "target raised by",
					RatingFrom: "Hold",
					RatingTo:   "Buy",
					TargetFrom: "$150.00",
					TargetTo:   "$180.50",
					Time:       "invalid-time-format", // Invalid time format
				},
			},
			NextPage: "",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Set up environment
	os.Setenv("API_URL", server.URL)
	os.Setenv("API_KEY", "test-api-key")
	defer func() {
		os.Unsetenv("API_URL")
		os.Unsetenv("API_KEY")
	}()

	client := NewExternalStockAPI()
	response, err := client.FetchStocks("")

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "error parsing Time")
}