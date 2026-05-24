package service_test

import (
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/jadecobra/agbalumo/internal/testutil"
)

func init() {
	testutil.CategorizationServiceConstructor = func(repo domain.ListingRepository) domain.CategorizationService {
		return service.NewCategorizationService(repo, &domain.CategoryCache{})
	}
	testutil.CSVServiceConstructor = func(geo domain.GeocodingService) domain.CSVService {
		csvSvc := service.NewCSVService()
		csvSvc.Geocoding = geo
		return csvSvc
	}
}
