package services

import (
	"math"
	"testing"
	"time"

	"recommender/internal/core/domain"
)

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name  string
		stock domain.Stock
		want  float64
	}{
		{
			name: "Cambio de Neutral a Buy con aumento en precio",
			stock: domain.Stock{
				TargetFrom: 100.0,
				TargetTo:   120.0, // +20% de cambio
				RatingFrom: "Neutral",
				RatingTo:   "Buy",                                       // Mejora en rating (+2)
				Brokerage:  "The Goldman Sachs Group",                   // Peso 1.5
				Time:       time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC), // Fecha fija
			},
			want: 3, // Calculado con la fórmula real
		},
		{
			name: "Cambio de Buy a Neutral con disminución en precio",
			stock: domain.Stock{
				TargetFrom: 150.0,
				TargetTo:   140.0, // -6.67% de cambio
				RatingFrom: "Buy",
				RatingTo:   "Neutral",                                    // Degradación (-2)
				Brokerage:  "JP Morgan",                                  // Peso 1.4
				Time:       time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC), // Fecha fija
			},
			want: -1.86, // Calculado con la fórmula real
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateScore(tt.stock)
			if math.Abs(got-tt.want) > 0.1 { // Comparación con margen de error
				t.Errorf("calculateScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
