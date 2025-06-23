package repository

import (
	"errors"
	"testing"
	"time"

	"recommender/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gorm.io/driver/sqlite"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Usamos el driver sqlite puro Go con la DSN para memoria compartida
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to in-memory DB: %v", err)
	}
	if err := db.AutoMigrate(&domain.Stock{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestCreateAndGetStockByTicker(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCockroachStockRepository(db)

	stock := &domain.Stock{
		Ticker:   "AAPL",
		Time:     time.Now(),
		TargetTo: 180.0,
	}

	err := repo.Create(stock)
	assert.Nil(t, err)

	result, err := repo.GetStockByTicker("AAPL")
	assert.Nil(t, err)
	assert.Equal(t, stock.Ticker, result.Ticker)
}

func TestGetStockByTicker_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCockroachStockRepository(db)

	result, err := repo.GetStockByTicker("INVALID")
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
