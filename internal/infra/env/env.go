package env

import (
	"log/slog"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
)

// AppEnv centralizes all application-wide dependencies.
// This simplifies handler constructors and facilitates AI agent interaction.
type AppEnv struct {
	DB                domain.ListingRepository
	Cfg               *config.Config
	Logger            *slog.Logger
	CSVService        domain.CSVService
	GeocodingSvc      domain.GeocodingService
	ImageSvc          domain.ImageService
	ListingSvc        domain.ListingService
	CategorizationSvc domain.CategorizationService
	MetricsSvc        domain.MetricsService
	CatCache          *domain.CategoryCache
}

// CategoryCache is moved to domain/category.go to avoid circular dependencies
// and because it's a domain-specific caching entity.

func NewAppEnv(db domain.ListingRepository, cfg *config.Config, logger *slog.Logger, csv domain.CSVService, geo domain.GeocodingService, img domain.ImageService, listing domain.ListingService, cat domain.CategorizationService, metrics domain.MetricsService) *AppEnv {
	return &AppEnv{
		DB:                db,
		Cfg:               cfg,
		Logger:            logger,
		CSVService:        csv,
		GeocodingSvc:      geo,
		ImageSvc:          img,
		ListingSvc:        listing,
		CategorizationSvc: cat,
		MetricsSvc:        metrics,
		CatCache:          &domain.CategoryCache{},
	}
}
