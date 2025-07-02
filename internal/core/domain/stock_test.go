package domain_test

import (
    "encoding/json"
    "testing"
    "time"

    "recommender/internal/core/domain"
)

func TestCreateStock(t *testing.T) {
    now := time.Now()
    s := domain.Stock{
        ID:         1,
        Ticker:     "AAPL",
        Company:    "Apple Inc.",
        Brokerage:  "Goldman Sachs",
        Action:     "Buy",
        RatingFrom: "Neutral",
        RatingTo:   "Buy",
        TargetFrom: 150.00,
        TargetTo:   180.00,
        Time:       now,
    }

    if s.Ticker != "AAPL" {
        t.Errorf("Expected Ticker to be 'AAPL', got '%s'", s.Ticker)
    }

    if s.Company != "Apple Inc." {
        t.Errorf("Expected Company to be 'Apple Inc.', got '%s'", s.Company)
    }

    if s.Time != now {
        t.Errorf("Expected Time to be '%v', got '%v'", now, s.Time)
    }
}

func TestAPIResponseMarshalling(t *testing.T) {
    now := time.Now()
    response := domain.APIResponse{
        Items: []domain.Stock{
            {
                ID:        1,
                Ticker:    "GOOG",
                Company:   "Google LLC",
                Brokerage: "Morgan Stanley",
                Action:    "Hold",
                Time:      now,
            },
        },
        NextPage: "page2",
    }

    data, err := json.Marshal(response)
    if err != nil {
        t.Fatalf("Error marshaling APIResponse: %v", err)
    }

    var decoded domain.APIResponse
    err = json.Unmarshal(data, &decoded)
    if err != nil {
        t.Fatalf("Error unmarshaling APIResponse: %v", err)
    }

    if decoded.NextPage != "page2" {
        t.Errorf("Expected NextPage 'page2', got '%s'", decoded.NextPage)
    }

    if len(decoded.Items) != 1 || decoded.Items[0].Ticker != "GOOG" {
        t.Errorf("Decoded item incorrect: %+v", decoded.Items)
    }
}
