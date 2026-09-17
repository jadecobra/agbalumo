package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/spf13/cobra"
)

var (
	flagSocialPillar    int
	flagSocialListingID string
	flagSocialCity      string
	flagSocialOutput    string
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
		repo := InitRepo()

		var out io.Writer = cmd.OutOrStdout()
		if flagSocialOutput != "" {
			cleanPath := filepath.Clean(flagSocialOutput)
			f, err := os.Create(cleanPath) // #nosec G304 -- user-supplied CLI output path
			ExitOnErr(err, "Failed to create output file")
			defer func() { _ = f.Close() }()
			out = io.MultiWriter(cmd.OutOrStdout(), f)
		}

		err := GenerateSocialDraft(repo, flagSocialPillar, flagSocialListingID, flagSocialCity, out)
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
}

// SocialState tracks rotation offsets and recently spotlighted listings.
type SocialState struct {
	Offsets          map[string]int `json:"offsets,omitempty"`
	RecentlyFeatured []string       `json:"recently_featured,omitempty"`
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

// GenerateSocialDraft queries verified listings and renders a publication-ready post draft.
func GenerateSocialDraft(repo domain.ListingRepository, pillar int, listingID, city string, w io.Writer) error {
	ctx := context.Background()

	rawListings, _, err := repo.FindAll(ctx, string(domain.Food), "", "", 0, 0, 0, "", "", false, 100, 0)
	if err != nil {
		return err
	}

	listings := filterListingsByCity(rawListings, city)
	state := loadSocialState()
	defer func() {
		saveSocialState(state)
	}()

	switch pillar {
	case 1:
		return renderPillar1QualityIndex(listings, city, &state, w)
	case 2:
		return renderPillar2AirportArrival(listings, city, &state, w)
	case 3:
		return renderPillar3MerchantSpotlight(listings, listingID, &state, w)
	case 4:
		return renderPillar4SubMetroCorridor(listings, city, &state, w)
	case 5:
		return renderPillar5CoverageGaps(listings, city, &state, w)
	default:
		return fmt.Errorf("unknown pillar: %d (supported: 1-5)", pillar)
	}
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
	q.Set("utm_source", "cli")
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

func writeSpotItem(b *strings.Builder, idx int, s domain.Listing, campaign string) {
	specialty := s.RegionalSpecialty
	if specialty == "" {
		specialty = "West African"
	}
	b.WriteString(fmt.Sprintf("%d. %s (%s)\n", idx+1, s.Title, s.City))
	b.WriteString(fmt.Sprintf("   ★ %.1f (%d reviews) · %s\n", s.Rating, s.ReviewCount, specialty))
	b.WriteString(fmt.Sprintf("   Reviews & Details: %s\n", buildTrackedURL("/listings/"+s.ID, campaign)))
	if s.ContactPhone != "" {
		b.WriteString(fmt.Sprintf("   Phone: %s\n", s.ContactPhone))
	}
	if s.WebsiteURL != "" {
		b.WriteString(fmt.Sprintf("   Menu/Order: %s\n", s.WebsiteURL))
	}
	b.WriteString("\n")
}

func renderPillar1QualityIndex(listings []domain.Listing, city string, state *SocialState, w io.Writer) error {
	candidates := filterCandidatesByRating(listings, 4.0)
	spots := rotateListings(candidates, "pillar_1_"+strings.ToLower(city), 4, state)
	campaign := "quality_index"

	var b strings.Builder
	b.WriteString("================================================================================\n")
	b.WriteString("[DRAFT - Pillar 1: The DFW African Food Quality Index]\n")
	b.WriteString("Recommended Destination: DFW Diaspora Facebook Groups / Page Feed\n")
	b.WriteString("================================================================================\n")
	b.WriteString("We built agbalumo because landing in a new city or moving across town shouldn't mean gambling on food quality. Right now our strongest network is in Dallas and Fort Worth.\n\n")
	b.WriteString("Here are verified West African spots in DFW where quality is backed by real community reviews:\n\n")

	for idx, s := range spots {
		writeSpotItem(&b, idx, s, campaign)
	}

	b.WriteString("Find verified spots, directions, and direct contact in under 60 seconds:\n")
	b.WriteString(fmt.Sprintf("%s\n\n", buildTrackedURL("/", campaign, [2]string{"city", city})))
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
	exploreCity  string
	exploreLabel string
	addPrompt    string
}

func writeCorridorSpotItem(b *strings.Builder, idx int, s domain.Listing, campaign string) {
	b.WriteString(fmt.Sprintf("%d. %s (%s)\n", idx+1, s.Title, s.City))
	if s.Rating > 0 {
		b.WriteString(fmt.Sprintf("   ★ %.1f (%d reviews)\n", s.Rating, s.ReviewCount))
	}
	detailLabel := "Reviews & Details"
	if campaign == "sub_metro_corridor" {
		detailLabel = "Reviews & Menu"
	}
	b.WriteString(fmt.Sprintf("   %s: %s\n", detailLabel, buildTrackedURL("/listings/"+s.ID, campaign)))
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

func renderCorridorDraft(spots []domain.Listing, cfg corridorConfig, w io.Writer) error {
	var b strings.Builder
	b.WriteString("================================================================================\n")
	b.WriteString(cfg.headerTitle + "\n")
	b.WriteString("Recommended Destination: " + cfg.destination + "\n")
	b.WriteString("================================================================================\n")
	b.WriteString(cfg.introLead + "\n\n")
	b.WriteString(cfg.introList + "\n\n")

	for idx, s := range spots {
		writeCorridorSpotItem(&b, idx, s, cfg.campaign)
	}

	b.WriteString(cfg.exploreLabel + "\n")
	b.WriteString(fmt.Sprintf("%s\n\n", buildTrackedURL("/", cfg.campaign, [2]string{"city", cfg.exploreCity})))
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

func renderPillar2AirportArrival(listings []domain.Listing, city string, state *SocialState, w io.Writer) error {
	candidates := filterAirportCandidates(listings)
	airportSpots := rotateListings(candidates, "pillar_2", 3, state)
	return renderCorridorDraft(airportSpots, corridorConfig{
		campaign:     "airport_corridor",
		headerTitle:  "[DRAFT - Pillar 2: Airport Corridor Cities (Arlington, Grand Prairie, Irving)]",
		destination:  "DFW Diaspora Groups / Community Feed",
		introLead:    "When landing at DFW or navigating the mid-cities corridor, finding African food shouldn't mean driving across the entire metroplex.",
		introList:    "Here are verified spots in the airport corridor cities (Arlington, Grand Prairie, Irving) with direct contact information:",
		exploreLabel: "Explore all airport corridor African food spots in under 60 seconds:",
		exploreCity:  "Arlington",
		addPrompt:    "Know another African-owned spot near the airport corridor we missed? Add it directly to the network:",
	}, w)
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
	for _, id := range state.RecentlyFeatured {
		recentMap[id] = true
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

func selectSpotlight(listings []domain.Listing, listingID string, state *SocialState) (*domain.Listing, error) {
	if len(listings) == 0 {
		return nil, fmt.Errorf("no listings available for spotlight")
	}
	if listingID != "" {
		return findListingByID(listings, listingID)
	}
	return rotateSpotlight(listings, state), nil
}

func recordSpotlightFeatured(state *SocialState, id string) {
	state.RecentlyFeatured = append(state.RecentlyFeatured, id)
	if len(state.RecentlyFeatured) > 50 {
		state.RecentlyFeatured = state.RecentlyFeatured[len(state.RecentlyFeatured)-50:]
	}
}

func renderPillar3MerchantSpotlight(listings []domain.Listing, listingID string, state *SocialState, w io.Writer) error {
	spotlight, err := selectSpotlight(listings, listingID, state)
	if err != nil {
		return err
	}
	recordSpotlightFeatured(state, spotlight.ID)

	campaign := "merchant_spotlight"

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

func renderPillar4SubMetroCorridor(listings []domain.Listing, city string, state *SocialState, w io.Writer) error {
	candidates := filterCollinCandidates(listings)
	collinSpots := rotateListings(candidates, "pillar_4", 4, state)
	return renderCorridorDraft(collinSpots, corridorConfig{
		campaign:     "sub_metro_corridor",
		headerTitle:  "[DRAFT - Pillar 4: Sub-Metro Corridor Guide (Collin County)]",
		destination:  "DFW Diaspora Groups (Plano / Frisco / North Dallas)",
		introLead:    "We don't need to head all the way into Central Dallas when craving authentic West African food.",
		introList:    "Collin County has a trusted cluster of verified African kitchens across Plano, Allen, McKinney, and Frisco:",
		exploreLabel: "Explore all North DFW and Collin County spots in under 60 seconds:",
		exploreCity:  "Plano",
		addPrompt:    "Know another African-owned kitchen in Collin County? Add it directly to the network:",
	}, w)
}

func renderPillar5CoverageGaps(listings []domain.Listing, city string, state *SocialState, w io.Writer) error {
	campaign := "coverage_gaps"
	var b strings.Builder
	b.WriteString("================================================================================\n")
	b.WriteString("[DRAFT - Pillar 5: Radical Transparency & Coverage Gaps]\n")
	b.WriteString("Recommended Destination: DFW Diaspora Groups / Facebook Feed (High Comment Volume)\n")
	b.WriteString("================================================================================\n")
	b.WriteString("We started Agbalumo to map African-owned businesses where we don't have to explain ourselves, starting with food. Right now, Dallas-Fort Worth is our strongest network with verified spots across Dallas, Plano, Arlington, Grand Prairie, and McKinney.\n\n")
	b.WriteString("Here are verified spots mapped so far with community reviews:\n\n")

	cityMap := make(map[string][]domain.Listing)
	for _, l := range listings {
		cityMap[l.City] = append(cityMap[l.City], l)
	}

	var cities []string
	for c := range cityMap {
		cities = append(cities, c)
	}
	sort.Strings(cities)

	for _, c := range cities {
		spots := rotateCitySpots(cityMap[c], c, state)

		b.WriteString(fmt.Sprintf("• %s:\n", c))
		for _, s := range spots {
			ratingStr := ""
			if s.Rating > 0 {
				ratingStr = fmt.Sprintf("★ %.1f (%d reviews) · ", s.Rating, s.ReviewCount)
			}
			link := buildTrackedURL("/listings/"+s.ID, campaign)
			b.WriteString(fmt.Sprintf("  - %s: %s%s\n", s.Title, ratingStr, link))
		}
		b.WriteString("\n")
	}

	b.WriteString("We know there are blind spots in Frisco, Garland, Denton, and Fort Worth.\n\n")
	b.WriteString("Who are we missing? Add your favorite auntie's spot or suya joint directly to the network in under 60 seconds:\n")
	b.WriteString(fmt.Sprintf("%s\n", buildTrackedURL("/", campaign, [2]string{"action", "post"})))
	b.WriteString("================================================================================\n")

	_, err := io.WriteString(w, b.String())
	return err
}
