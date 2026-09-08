package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/spf13/cobra"
)

var (
	flagSocialPillar   int
	flagSocialPlatform string
	flagSocialCity     string
	flagSocialOutput   string
)

var SocialCmd = &cobra.Command{
	Use:   "social",
	Short: "Deterministic social media content engine",
	Long: `Generate pre-formatted, verified social media post drafts directly from 
the database for Facebook, X, and WhatsApp following agbalumo core story principles.`,
}

var socialDraftCmd = &cobra.Command{
	Use:   "draft",
	Short: "Generate a publication-ready social draft for a specific pillar",
	Example: `  # Generate a Facebook post for Pillar 1 (Quality Index) in Dallas
  agbalumo social draft --pillar 1 --platform facebook

  # Generate an Airport arrival dispatch (Pillar 2)
  agbalumo social draft --pillar 2 --platform facebook

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

		err := GenerateSocialDraft(repo, flagSocialPillar, flagSocialPlatform, flagSocialCity, out)
		ExitOnErr(err, "Failed to generate social draft")

		if flagSocialOutput != "" {
			cmd.Printf("\n[Draft also saved to: %s]\n", flagSocialOutput)
		}
	},
}

func init() {
	SocialCmd.AddCommand(socialDraftCmd)

	socialDraftCmd.Flags().IntVar(&flagSocialPillar, "pillar", 1, "Content pillar: 1-5")
	socialDraftCmd.Flags().StringVar(&flagSocialPlatform, "platform", "facebook", "Target platform: facebook, x, whatsapp")
	socialDraftCmd.Flags().StringVar(&flagSocialCity, "city", "Dallas", "Target city or metro anchor")
	socialDraftCmd.Flags().StringVar(&flagSocialOutput, "output", "", "Optional output file path to save the draft")
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
func GenerateSocialDraft(repo domain.ListingRepository, pillar int, platform, city string, w io.Writer) error {
	ctx := context.Background()

	rawListings, _, err := repo.FindAll(ctx, string(domain.Food), "", "", 0, 0, 0, "", "", false, 100, 0)
	if err != nil {
		return err
	}

	listings := filterListingsByCity(rawListings, city)

	switch pillar {
	case 1:
		return renderPillar1QualityIndex(listings, platform, city, w)
	case 2:
		return renderPillar2AirportArrival(listings, platform, city, w)
	case 3:
		return renderPillar3MerchantSpotlight(listings, platform, city, w)
	case 4:
		return renderPillar4SubMetroCorridor(listings, platform, city, w)
	case 5:
		return renderPillar5CoverageGaps(listings, platform, city, w)
	default:
		return fmt.Errorf("unknown pillar: %d (supported: 1-5)", pillar)
	}
}

func getTopSpots(listings []domain.Listing) []domain.Listing {
	var top []domain.Listing
	for _, l := range listings {
		if l.Rating >= 4.0 {
			top = append(top, l)
		}
	}
	if len(top) == 0 && len(listings) > 0 {
		top = listings
	}
	sort.Slice(top, func(i, j int) bool {
		if top[i].Rating != top[j].Rating {
			return top[i].Rating > top[j].Rating
		}
		return top[i].ReviewCount > top[j].ReviewCount
	})
	if len(top) > 4 {
		top = top[:4]
	}
	return top
}

func renderPillar1QualityIndex(listings []domain.Listing, platform, city string, w io.Writer) error {
	topSpots := getTopSpots(listings)

	var b strings.Builder
	b.WriteString("================================================================================\n")
	b.WriteString(fmt.Sprintf("[DRAFT: %s - Pillar 1: The DFW African Food Quality Index]\n", strings.ToUpper(platform)))
	b.WriteString("Recommended Destination: DFW Diaspora Facebook Groups / Page Feed\n")
	b.WriteString("================================================================================\n")
	b.WriteString("We built agbalumo because landing in a new city or moving across town shouldn't mean gambling on food quality. Right now our strongest network is in Dallas and Fort Worth.\n\n")
	b.WriteString("Here are verified West African spots in DFW where quality is backed by real community reviews:\n\n")

	for idx, s := range topSpots {
		specialty := s.RegionalSpecialty
		if specialty == "" {
			specialty = "West African"
		}
		b.WriteString(fmt.Sprintf("%d. %s (%s)\n", idx+1, s.Title, s.City))
		b.WriteString(fmt.Sprintf("   ★ %.1f (%d reviews) · %s\n", s.Rating, s.ReviewCount, specialty))
		if s.ContactPhone != "" {
			b.WriteString(fmt.Sprintf("   Phone: %s\n", s.ContactPhone))
		}
		if s.WebsiteURL != "" {
			b.WriteString(fmt.Sprintf("   Menu/Order: %s\n", s.WebsiteURL))
		}
		b.WriteString("\n")
	}

	b.WriteString("Find verified spots, directions, and direct contact in under 60 seconds:\n")
	b.WriteString(fmt.Sprintf("https://agbalumo.com/?city=%s\n\n", city))
	b.WriteString("If we missed your trusted spot in DFW, let us know so we can verify and add it.\n")
	b.WriteString("================================================================================\n")

	_, err := io.WriteString(w, b.String())
	return err
}

func getAirportSpots(listings []domain.Listing) []domain.Listing {
	var spots []domain.Listing
	for _, l := range listings {
		c := strings.ToLower(l.City)
		if c == "arlington" || c == "grand prairie" || c == "irving" {
			spots = append(spots, l)
		}
	}
	if len(spots) == 0 && len(listings) > 0 {
		spots = listings
	}
	if len(spots) > 3 {
		spots = spots[:3]
	}
	return spots
}

func renderPillar2AirportArrival(listings []domain.Listing, platform, city string, w io.Writer) error {
	airportSpots := getAirportSpots(listings)

	var b strings.Builder
	b.WriteString("================================================================================\n")
	b.WriteString(fmt.Sprintf("[DRAFT: %s - Pillar 2: Airport & Late-Night Arrival Dispatch]\n", strings.ToUpper(platform)))
	b.WriteString("Recommended Destination: DFW Diaspora Groups / LinkedIn / X\n")
	b.WriteString("================================================================================\n")
	b.WriteString("We know the feeling of landing at DFW or Love Field after 8 PM on a weekday and just wanting food that tastes like home—without waiting 45 minutes or guessing if the kitchen is still open.\n\n")
	b.WriteString("Here are verified spots within 20 minutes of DFW airport terminals with active kitchens and direct phone ordering:\n\n")

	for idx, s := range airportSpots {
		b.WriteString(fmt.Sprintf("%d. %s (%s)\n", idx+1, s.Title, s.City))
		if s.ContactPhone != "" {
			b.WriteString(fmt.Sprintf("   Direct Phone: %s\n", s.ContactPhone))
		}
		if s.WebsiteURL != "" {
			b.WriteString(fmt.Sprintf("   Online Order: %s\n", s.WebsiteURL))
		}
		b.WriteString("\n")
	}

	b.WriteString("Explore all airport-area African food spots in under 60 seconds:\n")
	b.WriteString("https://agbalumo.com/?city=Arlington\n")
	b.WriteString("================================================================================\n")

	_, err := io.WriteString(w, b.String())
	return err
}

func renderPillar3MerchantSpotlight(listings []domain.Listing, platform, city string, w io.Writer) error {
	if len(listings) == 0 {
		return fmt.Errorf("no listings available for spotlight")
	}
	spotlight := listings[0]
	for _, l := range listings {
		if l.Rating > spotlight.Rating {
			spotlight = l
		}
	}

	var b strings.Builder
	b.WriteString("================================================================================\n")
	b.WriteString(fmt.Sprintf("[DRAFT: %s - Pillar 3: Merchant Reciprocity Spotlight]\n", strings.ToUpper(platform)))
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
	b.WriteString("\nView details, hours, and directions on agbalumo:\n")
	b.WriteString(fmt.Sprintf("https://agbalumo.com/listings/%s\n\n", spotlight.ID))
	b.WriteString(fmt.Sprintf("Tagging %s — thank you for serving the diaspora.\n", spotlight.Title))
	b.WriteString("================================================================================\n")

	_, err := io.WriteString(w, b.String())
	return err
}

func renderPillar4SubMetroCorridor(listings []domain.Listing, platform, city string, w io.Writer) error {
	var collinSpots []domain.Listing
	for _, l := range listings {
		c := strings.ToLower(l.City)
		if c == "plano" || c == "allen" || c == "mckinney" || c == "frisco" {
			collinSpots = append(collinSpots, l)
		}
	}
	if len(collinSpots) == 0 && len(listings) > 0 {
		collinSpots = listings
	}
	if len(collinSpots) > 4 {
		collinSpots = collinSpots[:4]
	}

	var b strings.Builder
	b.WriteString("================================================================================\n")
	b.WriteString(fmt.Sprintf("[DRAFT: %s - Pillar 4: Sub-Metro Corridor Guide (Collin County)]\n", strings.ToUpper(platform)))
	b.WriteString("Recommended Destination: DFW Diaspora Groups (Plano / Frisco / North Dallas)\n")
	b.WriteString("================================================================================\n")
	b.WriteString("We don't need to drive 45 minutes down to Central Dallas on a Thursday night just for good food.\n\n")
	b.WriteString("Collin County has a trusted cluster of verified West African kitchens right in Plano, Allen, and McKinney:\n\n")

	for idx, s := range collinSpots {
		b.WriteString(fmt.Sprintf("%d. %s (%s) · ★ %.1f\n", idx+1, s.Title, s.City, s.Rating))
		if s.ContactPhone != "" {
			b.WriteString(fmt.Sprintf("   Phone: %s\n", s.ContactPhone))
		}
	}

	b.WriteString("\nExplore all North DFW and Collin County spots in under 60 seconds:\n")
	b.WriteString("https://agbalumo.com/?city=Plano\n")
	b.WriteString("================================================================================\n")

	_, err := io.WriteString(w, b.String())
	return err
}

func renderPillar5CoverageGaps(listings []domain.Listing, platform, city string, w io.Writer) error {
	var b strings.Builder
	b.WriteString("================================================================================\n")
	b.WriteString(fmt.Sprintf("[DRAFT: %s - Pillar 5: Radical Transparency & Coverage Gaps]\n", strings.ToUpper(platform)))
	b.WriteString("Recommended Destination: DFW Diaspora Groups / Facebook Feed (High Comment Volume)\n")
	b.WriteString("================================================================================\n")
	b.WriteString("We started Agbalumo to map African-owned businesses where we don't have to explain ourselves, starting with food. Right now, Dallas-Fort Worth is our strongest network with verified spots across Dallas, Plano, Arlington, Grand Prairie, and McKinney.\n\n")
	b.WriteString("Here are spots we have verified so far:\n")

	cityMap := make(map[string][]string)
	for _, l := range listings {
		cityMap[l.City] = append(cityMap[l.City], l.Title)
	}

	for c, titles := range cityMap {
		b.WriteString(fmt.Sprintf("• %s: %s\n", c, strings.Join(titles, ", ")))
	}

	b.WriteString("\nWe know there are blind spots in Frisco, Garland, Denton, and Fort Worth.\n\n")
	b.WriteString("Who are we missing? Tell us your favorite auntie's spot or suya joint below so we can verify and add them to the network:\n")
	b.WriteString("https://agbalumo.com\n")
	b.WriteString("================================================================================\n")

	_, err := io.WriteString(w, b.String())
	return err
}
