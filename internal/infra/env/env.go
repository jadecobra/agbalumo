package env

import (
	"log/slog"
	"sync"
	"time"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
)

// AppEnv centralizes all application-wide dependencies.
// This simplifies handler constructors and facilitates AI agent interaction.
type AppEnv struct {
	DB           domain.ListingRepository
	Cfg          *config.Config
	Logger       *slog.Logger
	CSVService   domain.CSVService
	GeocodingSvc domain.GeocodingService
	ImageSvc     domain.ImageService
	ListingSvc   domain.ListingService
	CatCache     *CategoryCache
}

// CategoryCache is a simple thread-safe cache for category data.
type CategoryCache struct {
	mu         sync.RWMutex
	categories []domain.CategoryData
	expiration time.Time
}

func (c *CategoryCache) Get() ([]domain.CategoryData, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if time.Now().After(c.expiration) {
		return nil, false
	}
	return c.categories, true
}

func (c *CategoryCache) Set(categories []domain.CategoryData, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.categories = categories
	c.expiration = time.Now().Add(ttl)
}

func NewAppEnv(db domain.ListingRepository, cfg *config.Config, logger *slog.Logger, csv domain.CSVService, geo domain.GeocodingService, img domain.ImageService, listing domain.ListingService) *AppEnv {
	return &AppEnv{
		DB:           db,
		Cfg:          cfg,
		Logger:       logger,
		CSVService:   csv,
		GeocodingSvc: geo,
		ImageSvc:     img,
		ListingSvc:   listing,
		CatCache:     &CategoryCache{},
	}
}
