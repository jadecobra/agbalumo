package service

import (
	"context"
	"log/slog"

	"github.com/jadecobra/agbalumo/internal/domain"
)

type ScraperJob struct {
	repo    domain.ListingRepository
	scraper *WebsiteScraper
}

func NewScraperJob(repo domain.ListingRepository, scraper *WebsiteScraper) *ScraperJob {
	return &ScraperJob{
		repo:    repo,
		scraper: scraper,
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
	}

	return successCount, nil
}

func (j *ScraperJob) enrichSingle(ctx context.Context, l domain.Listing) bool {
	slog.Info("[ScraperJob] Enriching listing", slog.String("id", l.ID), slog.String("title", l.Title), slog.String("url", l.WebsiteURL))

	signals, err := j.scraper.ScrapeListing(ctx, l.WebsiteURL)
	if err != nil {
		slog.Error("[ScraperJob] Failed to scrape", slog.String("id", l.ID), slog.Any("error", err))
		return false
	}

	if j.isEmpty(signals) {
		slog.Info("[ScraperJob] No signals found for listing", slog.String("id", l.ID))
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
	return s.HeatLevel == 0 && s.PaymentMethods == "" && s.MenuURL == "" && s.TopDish == "" && s.RegionalSpecialty == ""
}

func (j *ScraperJob) applySignals(l *domain.Listing, signals AdaSignals) {
	l.HeatLevel = signals.HeatLevel
	l.PaymentMethods = signals.PaymentMethods
	l.MenuURL = signals.MenuURL
	if signals.TopDish != "" {
		l.TopDish = signals.TopDish
	}
	if signals.RegionalSpecialty != "" {
		l.RegionalSpecialty = signals.RegionalSpecialty
	}
}
