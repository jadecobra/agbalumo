package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

type ScraperJob struct {
	repo           domain.ListingRepository
	scraper        *WebsiteScraper
	hoursExtractor domain.HoursExtractor
}

func NewScraperJob(repo domain.ListingRepository, scraper *WebsiteScraper, hoursExtractor domain.HoursExtractor) *ScraperJob {
	return &ScraperJob{
		repo:           repo,
		scraper:        scraper,
		hoursExtractor: hoursExtractor,
	}
}

// EnrichListings finds unenriched listings and runs the scraper against their websites.
func (j *ScraperJob) EnrichListings(ctx context.Context, limit int) (int, error) {
	targets, err := j.repo.FindEnrichmentTargets(ctx, limit)
	if err != nil {
		return 0, err
	}

	successCount := 0
	for _, l := range targets {
		if j.enrichSingle(ctx, l) {
			successCount++
		}
		time.Sleep(2 * time.Second)
	}

	return successCount, nil
}

func (j *ScraperJob) extractListingHours(ctx context.Context, l *domain.Listing) {
	if l.HoursOfOperation == "" || j.hoursExtractor == nil {
		return
	}
	structured, extractErr := j.hoursExtractor.ExtractHours(ctx, l.HoursOfOperation)
	if extractErr == nil {
		l.StructuredHours = structured
	} else {
		slog.Error("[ScraperJob] Failed to extract structured hours", slog.String("id", l.ID), slog.Any("error", extractErr))
	}
}

func (j *ScraperJob) savePartial(ctx context.Context, l domain.Listing) {
	if saveErr := j.repo.Save(ctx, l); saveErr != nil {
		slog.Error("[ScraperJob] Failed to save partial state", slog.String("id", l.ID), slog.Any("error", saveErr))
	}
}

func (j *ScraperJob) enrichSingle(ctx context.Context, l domain.Listing) bool {
	slog.Info("[ScraperJob] Enriching listing", slog.String("id", l.ID), slog.String("title", l.Title), slog.String("url", l.WebsiteURL))

	signals, err := j.scraper.ScrapeListing(ctx, l.WebsiteURL)
	now := time.Now()
	l.EnrichmentAttemptedAt = &now

	j.extractListingHours(ctx, &l)

	if err != nil {
		slog.Error("[ScraperJob] Failed to scrape", slog.String("id", l.ID), slog.Any("error", err))
		j.savePartial(ctx, l)
		return false
	}

	if j.isEmpty(signals) {
		slog.Info("[ScraperJob] No signals found for listing", slog.String("id", l.ID))
		j.savePartial(ctx, l)
		return false
	}

	j.applySignals(&l, signals)

	if err := j.repo.Save(ctx, l); err != nil {
		slog.Error("[ScraperJob] Failed to save", slog.String("id", l.ID), slog.Any("error", err))
		return false
	}
	return true
}

func (j *ScraperJob) isEmpty(s AdaSignals) bool {
	return s.HeatLevel == 0 && s.PaymentMethods == "" && s.MenuURL == "" && s.RegionalSpecialty == ""
}

func (j *ScraperJob) applySignals(l *domain.Listing, signals AdaSignals) {
	l.HeatLevel = signals.HeatLevel
	l.PaymentMethods = signals.PaymentMethods
	l.MenuURL = signals.MenuURL
	if signals.RegionalSpecialty != "" {
		l.RegionalSpecialty = signals.RegionalSpecialty
	}
}
