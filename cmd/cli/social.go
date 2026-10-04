package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/spf13/cobra"
)

var (
	flagSocialPillar       int
	flagSocialListingID    string
	flagSocialCity         string
	flagSocialOutput       string
	flagSocialFailBadLinks bool
	flagSocialStrict       bool
)

var SocialCmd = &cobra.Command{
	Use:   "social",
	Short: "Deterministic social media content engine",
	Long: `Generate pre-formatted, verified social media post drafts directly from 
the database following agbalumo core story principles.`,
}

var socialDraftCmd = &cobra.Command{
	Use:   "draft",
	Short: "Generate a publication-ready social draft for a specific pillar",
	Example: `  # Generate a post for Pillar 1 (Quality Index) in Dallas
  agbalumo social draft --pillar 1

  # Generate an Airport corridor dispatch (Pillar 2)
  agbalumo social draft --pillar 2

  # Spotlight a specific merchant (Pillar 3)
  agbalumo social draft --pillar 3 --listing-id cli-12345

  # Save the draft directly to a text file for quick editing
  agbalumo social draft --pillar 5 --output post_draft.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := InitSocialRepo()
		ExitOnErr(err, "Failed to initialize database for social draft")

		var out io.Writer = cmd.OutOrStdout()
		if flagSocialOutput != "" {
			cleanPath := filepath.Clean(flagSocialOutput)
			f, createErr := os.Create(cleanPath) // #nosec G304 -- user-supplied CLI output path
			ExitOnErr(createErr, "Failed to create output file")
			defer func() { _ = f.Close() }()
			out = io.MultiWriter(cmd.OutOrStdout(), f)
		}

		var opts []SocialOption
		if flagSocialFailBadLinks {
			opts = append(opts, WithFailBadLinks(true))
		}
		if flagSocialStrict {
			opts = append(opts, WithStrictMode(true))
		}

		err = GenerateSocialDraft(repo, flagSocialPillar, flagSocialListingID, flagSocialCity, out, opts...)
		ExitOnErr(err, "Failed to generate social draft")

		if flagSocialOutput != "" {
			cmd.Printf("\n[Draft also saved to: %s]\n", flagSocialOutput)
		}
	},
}

func init() {
	SocialCmd.AddCommand(socialDraftCmd)

	socialDraftCmd.Flags().IntVarP(&flagSocialPillar, "pillar", "p", 1, "Content pillar: 1-5")
	socialDraftCmd.Flags().StringVarP(&flagSocialListingID, "listing-id", "l", "", "Target listing ID (for pillar 3 spotlight; rotates if omitted)")
	socialDraftCmd.Flags().StringVar(&flagSocialCity, "city", "Dallas", "Target city or metro anchor")
	socialDraftCmd.Flags().StringVarP(&flagSocialOutput, "output", "o", "", "Optional output file path to save the draft")
	socialDraftCmd.Flags().BoolVar(&flagSocialFailBadLinks, "fail-bad-links", false, "Fail draft immediately if any listing deep link returns non-2xx (default: omit link with loud stderr; pillar 3 always fails)")
	socialDraftCmd.Flags().BoolVar(&flagSocialStrict, "strict", false, "Fail the draft instead of warning if any lint check fails")
}

// FeaturedSpot records a listing featured in social drafts with its timestamp.
type FeaturedSpot struct {
	FeaturedAt time.Time `json:"featured_at"`
	ID         string    `json:"id"`
}

// SocialState tracks rotation offsets and recently spotlighted listings.
type SocialState struct {
	Offsets          map[string]int `json:"offsets,omitempty"`
	RecentlyFeatured []FeaturedSpot `json:"recently_featured,omitempty"`
}

func parseRawRecentlyFeatured(raw []json.RawMessage) []FeaturedSpot {
	spots := make([]FeaturedSpot, 0, len(raw))
	for _, item := range raw {
		var spot FeaturedSpot
		if err := json.Unmarshal(item, &spot); err == nil && spot.ID != "" {
			spots = append(spots, spot)
			continue
		}
		var strID string
		if err := json.Unmarshal(item, &strID); err == nil && strID != "" {
			spots = append(spots, FeaturedSpot{ID: strID})
		}
	}
	return spots
}

// UnmarshalJSON handles both legacy []string and current []FeaturedSpot formats.
func (s *SocialState) UnmarshalJSON(data []byte) error {
	type Alias SocialState
	var aux struct {
		*Alias
		RawRecently []json.RawMessage `json:"recently_featured,omitempty"`
	}
	aux.Alias = (*Alias)(s)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if len(aux.RawRecently) > 0 {
		s.RecentlyFeatured = parseRawRecentlyFeatured(aux.RawRecently)
	}
	return nil
}

func getStateFilePath() string {
	if p := os.Getenv("AGBALUMO_SOCIAL_STATE"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), ".agbalumo_social_state.json")
	}
	return filepath.Join(home, ".agbalumo", "social_state.json")
}

func loadSocialState() SocialState {
	path := filepath.Clean(getStateFilePath())
	data, err := os.ReadFile(path) // #nosec G304 -- local CLI state path
	if err != nil {
		return SocialState{Offsets: make(map[string]int)}
	}
	var s SocialState
	if err := json.Unmarshal(data, &s); err != nil {
		return SocialState{Offsets: make(map[string]int)}
	}
	if s.Offsets == nil {
		s.Offsets = make(map[string]int)
	}
	return s
}

func saveSocialState(s SocialState) {
	path := filepath.Clean(getStateFilePath())
	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0750)
	data, err := json.MarshalIndent(s, "", "  ")
	if err == nil {
		_ = os.WriteFile(path, data, 0600) // #nosec G304 -- local CLI state path
	}
}

func rotateListings(candidates []domain.Listing, key string, limit int, state *SocialState) []domain.Listing {
	if len(candidates) == 0 {
		return nil
	}
	if len(candidates) <= limit {
		return candidates
	}
	offset := state.Offsets[key] % len(candidates)
	result := make([]domain.Listing, limit)
	for i := 0; i < limit; i++ {
		result[i] = candidates[(offset+i)%len(candidates)]
	}
	state.Offsets[key] = (offset + limit) % len(candidates)
	return result
}

func rotateCitySpots(spots []domain.Listing, city string, state *SocialState) []domain.Listing {
	if len(spots) <= 1 {
		return spots
	}
	key := "pillar_5_" + strings.ToLower(city)
	offset := state.Offsets[key] % len(spots)
	rotated := make([]domain.Listing, len(spots))
	for i := 0; i < len(spots); i++ {
		rotated[i] = spots[(offset+i)%len(spots)]
	}
	state.Offsets[key] = (offset + 1) % len(spots)
	return rotated
}

func isDFW(city string) bool {
	c := strings.ToLower(strings.TrimSpace(city))
	return c == "dallas" || c == "plano" || c == "arlington" || c == "grand prairie" || c == "irving" || c == "allen" || c == "mckinney" || c == "fort worth" || c == "dfw"
}

func filterListingsByCity(raw []domain.Listing, city string) []domain.Listing {
	if city == "" {
		return raw
	}
	var filtered []domain.Listing
	dfwTarget := isDFW(city)
	for _, l := range raw {
		if (dfwTarget && isDFW(l.City)) || strings.EqualFold(l.City, city) {
			filtered = append(filtered, l)
		}
	}
	if len(filtered) == 0 {
		return raw
	}
	return filtered
}

func isProhibitedTesterDB(path string) bool {
	clean := filepath.Clean(path)
	return clean == filepath.Clean(domain.DefaultDatabaseURL) ||
		clean == ".tester/data/agbalumo.db" ||
		strings.HasSuffix(clean, filepath.Join(".tester", "data", "agbalumo.db"))
}

// GetSocialDatabaseURL resolves the database URL for social drafts.
// It prefers AGBALUMO_SOCIAL_DB, falls back to DATABASE_URL, or defaults to local prod_snapshot.db.
// It explicitly prohibits using .tester/data/agbalumo.db for shipping deep links.
func GetSocialDatabaseURL() (string, error) {
	if db := os.Getenv("AGBALUMO_SOCIAL_DB"); db != "" {
		if isProhibitedTesterDB(db) {
			return "", fmt.Errorf("AGBALUMO_SOCIAL_DB points to prohibited local tester database (%s); prod-parity database required", db)
		}
		return db, nil
	}
	if db := os.Getenv(domain.EnvKeyDatabaseURL); db != "" {
		if isProhibitedTesterDB(db) {
			return "", fmt.Errorf("DATABASE_URL points to prohibited local tester database (%s); prod-parity database required (prefer AGBALUMO_SOCIAL_DB)", db)
		}
		return db, nil
	}
	defaultSnapshot := filepath.Join(".tester", "data", "prod_snapshot.db")
	if _, err := os.Stat(defaultSnapshot); err == nil {
		return defaultSnapshot, nil
	}
	return "", fmt.Errorf("social draft requires a prod-parity database with matching UUIDs: %s is prohibited; set AGBALUMO_SOCIAL_DB or sync to %s", domain.DefaultDatabaseURL, defaultSnapshot)
}

// InitSocialRepo connects to the prod-parity database for social drafts.
func InitSocialRepo() (*sqlite.SQLiteRepository, error) {
	dbPath, err := GetSocialDatabaseURL()
	if err != nil {
		return nil, err
	}
	repo, err := sqlite.NewSQLiteRepository(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open social database at %s: %w", dbPath, err)
	}
	return repo, nil
}

// LinkVerifier checks whether a listing deep link is live and reachable.
type LinkVerifier interface {
	Verify(ctx context.Context, listingID string) (ok bool, statusCode int, err error)
}

// LinkVerifierFunc allows plain functions to act as LinkVerifier.
type LinkVerifierFunc func(ctx context.Context, listingID string) (bool, int, error)

func (f LinkVerifierFunc) Verify(ctx context.Context, listingID string) (bool, int, error) {
	return f(ctx, listingID)
}

// HTTPLinkVerifier performs live HTTP probes (HEAD, falling back to GET) against agbalumo deep links.
type HTTPLinkVerifier struct {
	Client  *http.Client
	BaseURL string
}

func NewHTTPLinkVerifier(baseURL string, client *http.Client) *HTTPLinkVerifier {
	if baseURL == "" {
		baseURL = os.Getenv("AGBALUMO_BASE_URL")
	}
	if baseURL == "" {
		baseURL = os.Getenv(domain.EnvKeyAppURL)
	}
	if baseURL == "" {
		baseURL = "https://agbalumo.com"
	}
	if client == nil {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	return &HTTPLinkVerifier{BaseURL: strings.TrimRight(baseURL, "/"), Client: client}
}

func (h *HTTPLinkVerifier) probeHEAD(ctx context.Context, targetURL string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, targetURL, nil)
	if err != nil {
		return false
	}
	resp, doErr := h.Client.Do(req)
	if doErr != nil {
		return false
	}
	_ = resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func (h *HTTPLinkVerifier) probeGET(ctx context.Context, targetURL string) (bool, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return false, 0, err
	}
	resp, doErr := h.Client.Do(req)
	if doErr != nil {
		return false, 0, doErr
	}
	_ = resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300, resp.StatusCode, nil
}

func (h *HTTPLinkVerifier) Verify(ctx context.Context, listingID string) (bool, int, error) {
	targetURL := fmt.Sprintf("%s/listings/%s", h.BaseURL, listingID)
	if h.probeHEAD(ctx, targetURL) {
		return true, http.StatusOK, nil
	}
	return h.probeGET(ctx, targetURL)
}

type socialConfig struct {
	verifier     LinkVerifier
	stderr       io.Writer
	failBadLinks bool
	strict       bool
}

type SocialOption func(*socialConfig)

func WithLinkVerifier(v LinkVerifier) SocialOption {
	return func(c *socialConfig) {
		c.verifier = v
	}
}

func WithStderr(w io.Writer) SocialOption {
	return func(c *socialConfig) {
		c.stderr = w
	}
}

func WithFailBadLinks(fail bool) SocialOption {
	return func(c *socialConfig) {
		c.failBadLinks = fail
	}
}

func WithStrictMode(strict bool) SocialOption {
	return func(c *socialConfig) {
		c.strict = strict
	}
}

func handleBadLink(stderr io.Writer, s domain.Listing, deepLink string, statusCode int, vErr error, failBadLinks bool) error {
	var detail string
	if vErr != nil {
		detail = vErr.Error()
	} else {
		detail = fmt.Sprintf("%s returned status %d", deepLink, statusCode)
	}

	if failBadLinks {
		return fmt.Errorf("deep link verification failed for listing %q (%s): %s", s.ID, s.Title, detail)
	}
	if stderr != nil {
		_, _ = fmt.Fprintf(stderr, "WARNING: deep link verification failed for listing %q (%s): %s; omitting link from draft\n", s.ID, s.Title, detail)
	}
	return nil
}

func verifyListingLink(ctx context.Context, verifier LinkVerifier, stderr io.Writer, s domain.Listing, campaign string, failBadLinks bool) (bool, error) {
	if verifier == nil {
		return true, nil
	}
	ok, statusCode, vErr := verifier.Verify(ctx, s.ID)
	if !ok || vErr != nil {
		deepLink := buildTrackedURL("/listings/"+s.ID, campaign)
		err := handleBadLink(stderr, s, deepLink, statusCode, vErr, failBadLinks)
		return false, err
	}
	return true, nil
}

type draftListingRef struct {
	ID    string
	Title string
	City  string
	URL   string
}

type draftData struct {
	RenderedText string
	City         string
	Listings     []draftListingRef
	BlindSpots   []string
	AllListings  []domain.Listing
	Pillar       int
}

func cleanTrackingParams(rawURL string) (bool, string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false, rawURL
	}
	q := u.Query()
	if len(q) == 0 {
		return false, rawURL
	}
	hasTracking := false
	for k := range q {
		lowerK := strings.ToLower(k)
		if lowerK == "v" || lowerK == "fbclid" || lowerK == "gclid" || strings.HasPrefix(lowerK, "utm_") {
			hasTracking = true
			break
		}
	}
	if !hasTracking {
		return false, rawURL
	}
	cleaned := &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
	}
	return true, cleaned.String()
}

func isRawBuilderDomain(rawURL string) (bool, string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false, ""
	}
	host := strings.ToLower(u.Hostname())
	if strings.HasSuffix(host, ".myshopify.com") {
		return true, "raw myshopify"
	}
	if strings.HasSuffix(host, ".wixsite.com") || host == "wixsite.com" {
		return true, "wixsite.com subdomain"
	}
	return false, ""
}

func extractDraftBody(rendered string) string {
	parts := strings.Split(rendered, "================================================================================\n")
	if len(parts) >= 3 {
		return strings.TrimSpace(parts[2])
	}
	return strings.TrimSpace(rendered)
}

func lintFormatID(id string) string {
	if strings.Contains(id, "-") && len(id) > 8 {
		parts := strings.Split(id, "-")
		if len(parts[0]) == 8 {
			return parts[0]
		}
	}
	return id
}

func lintMessyMenuLinks(listings []draftListingRef) []string {
	var warnings []string
	for _, l := range listings {
		if l.URL == "" {
			continue
		}
		shortID := lintFormatID(l.ID)
		if isBuilder, builderKind := isRawBuilderDomain(l.URL); isBuilder {
			warnings = append(warnings, fmt.Sprintf("lint: menu link is %s — %s (%s)", builderKind, l.Title, shortID))
		}
		if hasTracking, cleaned := cleanTrackingParams(l.URL); hasTracking {
			warnings = append(warnings, fmt.Sprintf("lint: menu link carries tracking params — %s (%s) (cleaned: %s)", l.Title, shortID, cleaned))
		}
	}
	return warnings
}

func lintCitiesConflict(listings []draftListingRef, blindSpots []string) []string {
	var warnings []string
	listedCities := make(map[string]bool)
	for _, l := range listings {
		if l.City != "" {
			listedCities[strings.ToLower(strings.TrimSpace(l.City))] = true
		}
	}
	for _, c := range blindSpots {
		cLower := strings.ToLower(strings.TrimSpace(c))
		if listedCities[cLower] {
			warnings = append(warnings, fmt.Sprintf("lint: city %s is both listed and called a blind spot", c))
		}
	}
	return warnings
}

func lintRepeatFeatures(listings []draftListingRef, recentlyFeatured []FeaturedSpot) []string {
	var warnings []string
	if len(recentlyFeatured) == 0 {
		return nil
	}
	cutoff := time.Now().UTC().Add(-72 * time.Hour)
	for _, l := range listings {
		for _, feat := range recentlyFeatured {
			if feat.ID == l.ID && (feat.FeaturedAt.IsZero() || feat.FeaturedAt.After(cutoff)) {
				warnings = append(warnings, fmt.Sprintf("lint: spot %s was featured in the last 3 days", l.ID))
				break
			}
		}
	}
	return warnings
}

func lintDraftLength(renderedText string) []string {
	var warnings []string
	linkCount := strings.Count(renderedText, "https://agbalumo.com/listings/")
	if linkCount > 10 {
		warnings = append(warnings, fmt.Sprintf("lint: draft has more than 10 listing links (%d links)", linkCount))
	}
	body := extractDraftBody(renderedText)
	if len(body) > 2000 {
		warnings = append(warnings, fmt.Sprintf("lint: draft body exceeds 2000 characters (%d chars)", len(body)))
	}
	return warnings
}

func runDraftLinter(draft draftData, state *SocialState) []string {
	var warnings []string
	warnings = append(warnings, lintMessyMenuLinks(draft.Listings)...)
	warnings = append(warnings, lintCitiesConflict(draft.Listings, draft.BlindSpots)...)
	if state != nil {
		warnings = append(warnings, lintRepeatFeatures(draft.Listings, state.RecentlyFeatured)...)
	}
	warnings = append(warnings, lintDraftLength(draft.RenderedText)...)
	return warnings
}

func recordDraftFeatured(state *SocialState, listings []draftListingRef) {
	now := time.Now().UTC()
	for _, l := range listings {
		state.RecentlyFeatured = append(state.RecentlyFeatured, FeaturedSpot{
			FeaturedAt: now,
			ID:         l.ID,
		})
	}
	if len(state.RecentlyFeatured) > 100 {
		state.RecentlyFeatured = state.RecentlyFeatured[len(state.RecentlyFeatured)-100:]
	}
}

func dispatchPillarRender(ctx context.Context, repo domain.ListingRepository, pillar int, listings []domain.Listing, listingID, city string, state *SocialState, w io.Writer, cfg *socialConfig, draft *draftData) error {
	switch pillar {
	case 1:
		return renderPillar1QualityIndex(ctx, listings, city, state, w, cfg, draft)
	case 2:
		return renderPillar2AirportArrival(ctx, listings, city, state, w, cfg, draft)
	case 3:
		return renderPillar3MerchantSpotlight(ctx, repo, listings, listingID, state, w, cfg, draft)
	case 4:
		return renderPillar4SubMetroCorridor(ctx, listings, city, state, w, cfg, draft)
	case 5:
		return renderPillar5CoverageGaps(ctx, listings, city, state, w, cfg, draft)
	default:
		return fmt.Errorf("unknown pillar: %d (supported: 1-5)", pillar)
	}
}

func handleDraftWarnings(warnings []string, stderr io.Writer, strict bool) error {
	if len(warnings) == 0 {
		return nil
	}
	if stderr != nil {
		for _, wMsg := range warnings {
			_, _ = fmt.Fprintln(stderr, wMsg)
		}
	}
	if strict {
		return fmt.Errorf("lint check failed: %s", strings.Join(warnings, "; "))
	}
	return nil
}

// GenerateSocialDraft queries verified listings and renders a publication-ready post draft.
func GenerateSocialDraft(repo domain.ListingRepository, pillar int, listingID, city string, w io.Writer, opts ...SocialOption) error {
	ctx := context.Background()

	cfg := &socialConfig{stderr: os.Stderr}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.verifier == nil && os.Getenv("AGBALUMO_SKIP_LINK_VERIFY") != "true" {
		cfg.verifier = NewHTTPLinkVerifier("", nil)
	}

	rawListings, _, findErr := repo.FindAll(ctx, string(domain.Food), "", "", 0, 0, 0, "", "", false, 100, 0)
	if findErr != nil {
		return findErr
	}

	listings := filterListingsByCity(rawListings, city)
	state := loadSocialState()

	var buf bytes.Buffer
	draft := draftData{
		Pillar:      pillar,
		City:        city,
		AllListings: listings,
	}

	if renderErr := dispatchPillarRender(ctx, repo, pillar, listings, listingID, city, &state, &buf, cfg, &draft); renderErr != nil {
		return renderErr
	}

	draft.RenderedText = buf.String()

	warnings := runDraftLinter(draft, &state)
	if warnErr := handleDraftWarnings(warnings, cfg.stderr, cfg.strict); warnErr != nil {
		return warnErr
	}

	recordDraftFeatured(&state, draft.Listings)
	saveSocialState(state)

	_, writeErr := io.WriteString(w, draft.RenderedText)
	return writeErr
}

func buildTrackedURL(path, campaign string, extraParams ...[2]string) string {
	baseURL := "https://agbalumo.com"
	u, err := url.Parse(baseURL + path)
	if err != nil {
		return baseURL + path
	}
	q := u.Query()
	for _, p := range extraParams {
		if p[0] != "" && p[1] != "" {
			q.Set(p[0], p[1])
		}
	}
	q.Set("utm_source", "facebook")
	q.Set("utm_medium", "social")
	if campaign != "" {
		q.Set("utm_campaign", strings.ToLower(campaign))
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func filterCandidatesByRating(listings []domain.Listing, minRating float64) []domain.Listing {
	var candidates []domain.Listing
	for _, l := range listings {
		if l.Rating >= minRating {
			candidates = append(candidates, l)
		}
	}
	if len(candidates) == 0 {
		return listings
	}
	return candidates
}

func writeSpotItem(b *strings.Builder, idx int, s domain.Listing, campaign string, includeLink bool) {
	specialty := s.RegionalSpecialty
	if specialty == "" {
		specialty = "West African"
	}
	b.WriteString(fmt.Sprintf("%d. %s (%s)\n", idx+1, s.Title, s.City))
	b.WriteString(fmt.Sprintf("   ★ %.1f (%d reviews) · %s\n", s.Rating, s.ReviewCount, specialty))
	if includeLink {
		b.WriteString(fmt.Sprintf("   Reviews & Details: %s\n", buildTrackedURL("/listings/"+s.ID, campaign)))
	}
	if s.ContactPhone != "" {
		b.WriteString(fmt.Sprintf("   Phone: %s\n", s.ContactPhone))
	}
	if s.WebsiteURL != "" {
		b.WriteString(fmt.Sprintf("   Menu/Order: %s\n", s.WebsiteURL))
	}
	b.WriteString("\n")
}

func renderPillar1QualityIndex(ctx context.Context, listings []domain.Listing, city string, state *SocialState, w io.Writer, cfg *socialConfig, draft *draftData) error {
	candidates := filterCandidatesByRating(listings, 4.0)
	spots := rotateListings(candidates, "pillar_1_"+strings.ToLower(city), 4, state)
	campaign := "quality_index"

	for _, s := range spots {
		draft.Listings = append(draft.Listings, draftListingRef{
			ID:    s.ID,
			Title: s.Title,
			City:  s.City,
			URL:   s.WebsiteURL,
		})
	}

	var b strings.Builder
	b.WriteString("================================================================================\n")
	b.WriteString("[DRAFT - Pillar 1: The DFW African Food Quality Index]\n")
	b.WriteString("Recommended Destination: DFW Diaspora Facebook Groups / Page Feed\n")
	b.WriteString("================================================================================\n")
	b.WriteString("We built agbalumo because landing in a new city or moving across town shouldn't mean gambling on food quality. Right now our strongest network is in Dallas and Fort Worth.\n\n")
	b.WriteString("Here are verified West African spots in DFW where quality is backed by real community reviews:\n\n")

	for idx, s := range spots {
		includeLink, err := verifyListingLink(ctx, cfg.verifier, cfg.stderr, s, campaign, cfg.failBadLinks)
		if err != nil {
			return err
		}
		writeSpotItem(&b, idx, s, campaign, includeLink)
	}

	b.WriteString("Find verified spots, directions, and direct contact in under 60 seconds:\n")
	b.WriteString(fmt.Sprintf("%s\n\n", buildTrackedURL("/", campaign)))
	b.WriteString("If we missed your trusted spot in DFW, add it directly to the network in under 60 seconds:\n")
	b.WriteString(fmt.Sprintf("%s\n", buildTrackedURL("/", campaign, [2]string{"action", "post"})))
	b.WriteString("================================================================================\n")

	_, err := io.WriteString(w, b.String())
	return err
}

type corridorConfig struct {
	campaign     string
	headerTitle  string
	destination  string
	introLead    string
	introList    string
	exploreLabel string
	addPrompt    string
}

func writeCorridorSpotItem(b *strings.Builder, idx int, s domain.Listing, campaign string, includeLink bool) {
	b.WriteString(fmt.Sprintf("%d. %s (%s)\n", idx+1, s.Title, s.City))
	if s.Rating > 0 {
		b.WriteString(fmt.Sprintf("   ★ %.1f (%d reviews)\n", s.Rating, s.ReviewCount))
	}
	if includeLink {
		detailLabel := "Reviews & Details"
		if campaign == "sub_metro_corridor" {
			detailLabel = "Reviews & Menu"
		}
		b.WriteString(fmt.Sprintf("   %s: %s\n", detailLabel, buildTrackedURL("/listings/"+s.ID, campaign)))
	}
	if s.ContactPhone != "" {
		phoneLabel := "Phone"
		if campaign == "airport_corridor" {
			phoneLabel = "Direct Phone"
		}
		b.WriteString(fmt.Sprintf("   %s: %s\n", phoneLabel, s.ContactPhone))
	}
	if s.WebsiteURL != "" {
		b.WriteString(fmt.Sprintf("   Online Order: %s\n", s.WebsiteURL))
	}
	b.WriteString("\n")
}

func renderCorridorDraft(ctx context.Context, spots []domain.Listing, cfg corridorConfig, w io.Writer, sCfg *socialConfig, draft *draftData) error {
	for _, s := range spots {
		draft.Listings = append(draft.Listings, draftListingRef{
			ID:    s.ID,
			Title: s.Title,
			City:  s.City,
			URL:   s.WebsiteURL,
		})
	}

	var b strings.Builder
	b.WriteString("================================================================================\n")
	b.WriteString(cfg.headerTitle + "\n")
	b.WriteString("Recommended Destination: " + cfg.destination + "\n")
	b.WriteString("================================================================================\n")
	b.WriteString(cfg.introLead + "\n\n")
	b.WriteString(cfg.introList + "\n\n")

	for idx, s := range spots {
		includeLink, err := verifyListingLink(ctx, sCfg.verifier, sCfg.stderr, s, cfg.campaign, sCfg.failBadLinks)
		if err != nil {
			return err
		}
		writeCorridorSpotItem(&b, idx, s, cfg.campaign, includeLink)
	}

	b.WriteString(cfg.exploreLabel + "\n")
	b.WriteString(fmt.Sprintf("%s\n\n", buildTrackedURL("/", cfg.campaign)))
	b.WriteString(cfg.addPrompt + "\n")
	b.WriteString(fmt.Sprintf("%s\n", buildTrackedURL("/", cfg.campaign, [2]string{"action", "post"})))
	b.WriteString("================================================================================\n")

	_, err := io.WriteString(w, b.String())
	return err
}

func filterAirportCandidates(listings []domain.Listing) []domain.Listing {
	var candidates []domain.Listing
	for _, l := range listings {
		c := strings.ToLower(l.City)
		if c == "arlington" || c == "grand prairie" || c == "irving" {
			candidates = append(candidates, l)
		}
	}
	if len(candidates) == 0 {
		return listings
	}
	return candidates
}

func renderPillar2AirportArrival(ctx context.Context, listings []domain.Listing, city string, state *SocialState, w io.Writer, cfg *socialConfig, draft *draftData) error {
	candidates := filterAirportCandidates(listings)
	airportSpots := rotateListings(candidates, "pillar_2", 3, state)
	return renderCorridorDraft(ctx, airportSpots, corridorConfig{
		campaign:     "airport_corridor",
		headerTitle:  "[DRAFT - Pillar 2: Airport Corridor Cities (Arlington, Grand Prairie, Irving)]",
		destination:  "DFW Diaspora Groups / Community Feed",
		introLead:    "When landing at DFW or navigating the mid-cities corridor, finding African food shouldn't mean driving across the entire metroplex.",
		introList:    "Here are verified spots in the airport corridor cities (Arlington, Grand Prairie, Irving) with direct contact information:",
		exploreLabel: "Explore all airport corridor African food spots in under 60 seconds:",
		addPrompt:    "Know another African-owned spot near the airport corridor we missed? Add it directly to the network:",
	}, w, cfg, draft)
}

func findListingByID(listings []domain.Listing, id string) (*domain.Listing, error) {
	for _, l := range listings {
		if l.ID == id {
			chosen := l
			return &chosen, nil
		}
	}
	return nil, fmt.Errorf("listing with ID %q not found", id)
}

func rotateSpotlight(listings []domain.Listing, state *SocialState) *domain.Listing {
	recentMap := make(map[string]bool)
	for _, spot := range state.RecentlyFeatured {
		recentMap[spot.ID] = true
	}

	var unfeatured []domain.Listing
	for _, l := range listings {
		if !recentMap[l.ID] {
			unfeatured = append(unfeatured, l)
		}
	}

	pool := unfeatured
	if len(pool) == 0 {
		state.RecentlyFeatured = nil
		pool = listings
	}

	offset := state.Offsets["pillar_3"] % len(pool)
	chosen := pool[offset]
	state.Offsets["pillar_3"] = (offset + 1) % len(pool)
	return &chosen
}

func selectSpotlight(ctx context.Context, repo domain.ListingRepository, listings []domain.Listing, listingID string, state *SocialState) (*domain.Listing, error) {
	if listingID != "" {
		if l, err := findListingByID(listings, listingID); err == nil {
			return l, nil
		}
		if repo != nil {
			l, err := repo.FindByID(ctx, listingID)
			if err == nil && l.ID != "" {
				return &l, nil
			}
		}
		return nil, fmt.Errorf("listing with ID %q not found", listingID)
	}
	if len(listings) == 0 {
		return nil, fmt.Errorf("no listings available for spotlight")
	}
	return rotateSpotlight(listings, state), nil
}

func renderPillar3MerchantSpotlight(ctx context.Context, repo domain.ListingRepository, listings []domain.Listing, listingID string, state *SocialState, w io.Writer, cfg *socialConfig, draft *draftData) error {
	spotlight, err := selectSpotlight(ctx, repo, listings, listingID, state)
	if err != nil {
		return err
	}

	campaign := "merchant_spotlight"

	// Verify deep link before emit: for Pillar 3, non-2xx must fail the draft!
	includeLink, err := verifyListingLink(ctx, cfg.verifier, cfg.stderr, *spotlight, campaign, true)
	if err != nil {
		return err
	}
	if !includeLink {
		return fmt.Errorf("deep link verification failed for listing %q (%s)", spotlight.ID, spotlight.Title)
	}

	draft.Listings = append(draft.Listings, draftListingRef{
		ID:    spotlight.ID,
		Title: spotlight.Title,
		City:  spotlight.City,
		URL:   spotlight.WebsiteURL,
	})

	var b strings.Builder
	b.WriteString("================================================================================\n")
	b.WriteString("[DRAFT - Pillar 3: Merchant Spotlight]\n")
	b.WriteString("Recommended Destination: Facebook Page / Tag Venue on Instagram / X\n")
	b.WriteString("================================================================================\n")
	b.WriteString(fmt.Sprintf("Spotlight: %s (%s, TX)\n", spotlight.Title, spotlight.City))
	b.WriteString(fmt.Sprintf("★ %.1f rating across %d community reviews.\n\n", spotlight.Rating, spotlight.ReviewCount))

	specialty := spotlight.RegionalSpecialty
	if specialty == "" {
		specialty = "West African"
	}
	b.WriteString(fmt.Sprintf("Serving authentic %s dishes with verified contact:\n", specialty))
	if spotlight.ContactPhone != "" {
		b.WriteString(fmt.Sprintf("• Direct Line: %s\n", spotlight.ContactPhone))
	}
	if spotlight.WebsiteURL != "" {
		b.WriteString(fmt.Sprintf("• Menu/Ordering: %s\n", spotlight.WebsiteURL))
	}
	b.WriteString("\nView reviews, hours, and directions on agbalumo:\n")
	b.WriteString(fmt.Sprintf("%s\n\n", buildTrackedURL("/listings/"+spotlight.ID, campaign)))
	b.WriteString(fmt.Sprintf("Tagging %s — thank you for serving the diaspora.\n", spotlight.Title))
	b.WriteString("================================================================================\n")

	_, err = io.WriteString(w, b.String())
	return err
}

func isCollinCounty(city string) bool {
	switch strings.ToLower(city) {
	case "plano", "allen", "mckinney", "frisco":
		return true
	default:
		return false
	}
}

func filterCollinCandidates(listings []domain.Listing) []domain.Listing {
	var candidates []domain.Listing
	for _, l := range listings {
		if isCollinCounty(l.City) {
			candidates = append(candidates, l)
		}
	}
	if len(candidates) == 0 {
		return listings
	}
	return candidates
}

func renderPillar4SubMetroCorridor(ctx context.Context, listings []domain.Listing, city string, state *SocialState, w io.Writer, cfg *socialConfig, draft *draftData) error {
	candidates := filterCollinCandidates(listings)
	collinSpots := rotateListings(candidates, "pillar_4", 4, state)
	return renderCorridorDraft(ctx, collinSpots, corridorConfig{
		campaign:     "sub_metro_corridor",
		headerTitle:  "[DRAFT - Pillar 4: Sub-Metro Corridor Guide (Collin County)]",
		destination:  "DFW Diaspora Groups (Plano / Frisco / North Dallas)",
		introLead:    "We don't need to head all the way into Central Dallas when craving authentic West African food.",
		introList:    "Collin County has a trusted cluster of verified African kitchens across Plano, Allen, McKinney, and Frisco:",
		exploreLabel: "Explore all North DFW and Collin County spots in under 60 seconds:",
		addPrompt:    "Know another African-owned kitchen in Collin County? Add it directly to the network:",
	}, w, cfg, draft)
}

func renderCitySpotItem(b *strings.Builder, s domain.Listing, campaign string, includeLink bool) {
	ratingStr := ""
	if s.Rating > 0 {
		ratingStr = fmt.Sprintf("★ %.1f (%d reviews)", s.Rating, s.ReviewCount)
	}
	if !includeLink {
		if ratingStr != "" {
			b.WriteString(fmt.Sprintf("  - %s: %s\n", s.Title, ratingStr))
		} else {
			b.WriteString(fmt.Sprintf("  - %s\n", s.Title))
		}
		return
	}
	if ratingStr != "" {
		ratingStr += " · "
	}
	link := buildTrackedURL("/listings/"+s.ID, campaign)
	b.WriteString(fmt.Sprintf("  - %s: %s%s\n", s.Title, ratingStr, link))
}

func joinNatural(items []string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	case 2:
		return items[0] + " and " + items[1]
	default:
		return strings.Join(items[:len(items)-1], ", ") + ", and " + items[len(items)-1]
	}
}

type cityGroup struct {
	name  string
	spots []domain.Listing
}

func groupAndSortCities(listings []domain.Listing) []cityGroup {
	cityMap := make(map[string][]domain.Listing)
	for _, l := range listings {
		c := strings.TrimSpace(l.City)
		if c == "" {
			continue
		}
		cityMap[c] = append(cityMap[c], l)
	}

	groups := make([]cityGroup, 0, len(cityMap))
	for c, spots := range cityMap {
		groups = append(groups, cityGroup{name: c, spots: spots})
	}

	sort.SliceStable(groups, func(i, j int) bool {
		if len(groups[i].spots) != len(groups[j].spots) {
			return len(groups[i].spots) > len(groups[j].spots)
		}
		return groups[i].name < groups[j].name
	})
	return groups
}

func sortCitySpots(spots []domain.Listing) []domain.Listing {
	sorted := make([]domain.Listing, len(spots))
	copy(sorted, spots)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Rating != sorted[j].Rating {
			return sorted[i].Rating > sorted[j].Rating
		}
		if sorted[i].ReviewCount != sorted[j].ReviewCount {
			return sorted[i].ReviewCount > sorted[j].ReviewCount
		}
		return sorted[i].Title < sorted[j].Title
	})
	return sorted
}

func renderCitySpotWithCount(ctx context.Context, b *strings.Builder, g cityGroup, state *SocialState, campaign string, cfg *socialConfig, draft *draftData) error {
	spots := sortCitySpots(g.spots)
	leadSpot := rotateCitySpots(spots, g.name, state)[0]

	draft.Listings = append(draft.Listings, draftListingRef{
		ID:    leadSpot.ID,
		Title: leadSpot.Title,
		City:  leadSpot.City,
		URL:   leadSpot.WebsiteURL,
	})

	b.WriteString(fmt.Sprintf("• %s:\n", g.name))
	includeLink, err := verifyListingLink(ctx, cfg.verifier, cfg.stderr, leadSpot, campaign, cfg.failBadLinks)
	if err != nil {
		return err
	}
	renderCitySpotItem(b, leadSpot, campaign, includeLink)
	if len(g.spots) > 1 {
		b.WriteString(fmt.Sprintf("  +%d more in %s\n", len(g.spots)-1, g.name))
	}
	b.WriteString("\n")
	return nil
}

func qualifyBlindSpots(listings []domain.Listing) []string {
	watchlist := []string{
		"Frisco",
		"Garland",
		"Denton",
		"Mesquite",
		"Richardson",
		"Carrollton",
		"Lewisville",
		"Fort Worth",
	}

	cityCounts := make(map[string]int)
	for _, l := range listings {
		cityCounts[strings.ToLower(strings.TrimSpace(l.City))]++
	}

	var qualifying []string
	for _, wCity := range watchlist {
		if cityCounts[strings.ToLower(wCity)] <= 1 {
			qualifying = append(qualifying, wCity)
		}
	}
	return qualifying
}

func writePillar5Intro(b *strings.Builder, displayedGroups []cityGroup) {
	b.WriteString("================================================================================\n")
	b.WriteString("[DRAFT - Pillar 5: Radical Transparency & Coverage Gaps]\n")
	b.WriteString("Recommended Destination: DFW Diaspora Groups / Facebook Feed (High Comment Volume)\n")
	b.WriteString("================================================================================\n")

	if len(displayedGroups) == 0 {
		b.WriteString("We started Agbalumo to map African-owned businesses where we don't have to explain ourselves, starting with food. Right now, Dallas-Fort Worth is our strongest network.\n\n")
		b.WriteString("Here are verified spots mapped so far with community reviews:\n\n")
		return
	}

	names := make([]string, len(displayedGroups))
	for i, g := range displayedGroups {
		names[i] = g.name
	}
	b.WriteString(fmt.Sprintf("We started Agbalumo to map African-owned businesses where we don't have to explain ourselves, starting with food. Right now, Dallas-Fort Worth is our strongest network with verified spots across %s.\n\n", joinNatural(names)))
	b.WriteString("Here are verified spots mapped so far with community reviews:\n\n")
}

func renderPillar5CoverageGaps(ctx context.Context, listings []domain.Listing, city string, state *SocialState, w io.Writer, cfg *socialConfig, draft *draftData) error {
	campaign := "coverage_gaps"
	groups := groupAndSortCities(listings)

	maxCities := 8
	if len(groups) < maxCities {
		maxCities = len(groups)
	}
	displayedGroups := groups[:maxCities]
	extraGroups := groups[maxCities:]

	var b strings.Builder
	writePillar5Intro(&b, displayedGroups)

	for _, g := range displayedGroups {
		if err := renderCitySpotWithCount(ctx, &b, g, state, campaign, cfg, draft); err != nil {
			return err
		}
	}

	if len(extraGroups) > 0 {
		names := make([]string, len(extraGroups))
		for i, g := range extraGroups {
			names[i] = g.name
		}
		b.WriteString(fmt.Sprintf("Also mapped in %s.\n\n", joinNatural(names)))
	}

	blindSpots := qualifyBlindSpots(listings)
	draft.BlindSpots = blindSpots
	if len(blindSpots) > 0 {
		b.WriteString(fmt.Sprintf("We know there are blind spots in %s.\n\n", joinNatural(blindSpots)))
	}

	b.WriteString("Find verified spots, directions, and direct contact in under 60 seconds:\n")
	b.WriteString(fmt.Sprintf("%s\n\n", buildTrackedURL("/", campaign)))
	b.WriteString("Who are we missing? Add your favorite auntie's spot or suya joint directly to the network in under 60 seconds:\n")
	b.WriteString(fmt.Sprintf("%s\n", buildTrackedURL("/", campaign, [2]string{"action", "post"})))
	b.WriteString("================================================================================\n")

	_, err := io.WriteString(w, b.String())
	return err
}
