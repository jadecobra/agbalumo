package service

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type AdaSignals struct {
	PaymentMethods    string
	MenuURL           string
	TopDish           string
	RegionalSpecialty string
	HeatLevel         int
}

type WebsiteScraper struct {
	client *http.Client
}

func NewWebsiteScraper() *WebsiteScraper {
	return &WebsiteScraper{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *WebsiteScraper) ScrapeListing(ctx context.Context, websiteURL string) (AdaSignals, error) {
	if websiteURL == "" {
		return AdaSignals{}, nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", websiteURL, nil)
	if err != nil {
		return AdaSignals{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; AgbalumoBot/1.0; +https://agbalumo.com)")

	resp, err := s.client.Do(req)
	if err != nil {
		return AdaSignals{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return AdaSignals{}, nil
	}

	// Limit reader to 512KB to prevent memory exhaustion
	limitBody := io.LimitReader(resp.Body, 512*1024)

	return s.parseHTML(limitBody, websiteURL), nil
}

type scrapeState struct {
	regionalCounts   map[string]int
	regionalKeywords map[string][]string
	currentAnchorURL string
	foundPayments    []string
	heatKeywords     []string
	paymentKeywords  []string
	heatCount        int
	inAnchor         bool
}

func (s *WebsiteScraper) parseHTML(r io.Reader, baseURL string) AdaSignals {
	var signals AdaSignals
	z := html.NewTokenizer(r)
	parsedBase, _ := url.Parse(baseURL)

	state := &scrapeState{
		heatKeywords:    []string{"spicy", "hot", "pepper", "habanero", "scotch bonnet", "chili"},
		paymentKeywords: []string{"zelle", "venmo", "cashapp", "cash app"},
		regionalCounts:  make(map[string]int),
		regionalKeywords: map[string][]string{
			"Nigerian":      {"nigeria", "jollof", "egusi", "suya", "lagos", "naija"},
			"Ghanaian":      {"ghana", "waakye", "shito", "kenkey", "accra"},
			"Senegalese":    {"senegal", "thieboudienne", "yassa", "dakar"},
			"Ethiopian":     {"ethiopia", "injera", "wat", "addis"},
			"South African": {"south africa", "braai", "biltong", "bobotie"},
		},
	}

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		s.processToken(z, tt, parsedBase, state, &signals)
	}

	signals.HeatLevel = s.mapHeatLevel(state.heatCount)
	signals.PaymentMethods = strings.Join(state.foundPayments, ", ")
	signals.RegionalSpecialty = s.inferRegionalSpecialty(state)
	return signals
}

func (s *WebsiteScraper) inferRegionalSpecialty(state *scrapeState) string {
	maxCount := 0
	bestRegion := ""
	for region, count := range state.regionalCounts {
		if count > maxCount {
			maxCount = count
			bestRegion = region
		}
	}
	return bestRegion
}

func (s *WebsiteScraper) processToken(z *html.Tokenizer, tt html.TokenType, base *url.URL, state *scrapeState, signals *AdaSignals) {
	switch tt {
	case html.StartTagToken, html.SelfClosingTagToken:
		s.handleTag(z, base, state, signals)
	case html.EndTagToken:
		tn, _ := z.TagName()
		if string(tn) == "a" {
			state.inAnchor = false
			state.currentAnchorURL = ""
		}
	case html.TextToken:
		s.handleText(z, state, base, signals)
	}
}

func (s *WebsiteScraper) handleTag(z *html.Tokenizer, base *url.URL, state *scrapeState, signals *AdaSignals) {
	tn, hasAttr := z.TagName()
	tagName := string(tn)

	switch tagName {
	case "h1", "h2":
		s.handleHeading(z, signals)
	case "a":
		if hasAttr {
			s.handleAnchor(z, base, state, signals)
		}
	}
}

func (s *WebsiteScraper) handleHeading(z *html.Tokenizer, signals *AdaSignals) {
	if signals.TopDish != "" {
		return
	}
	z.Next()
	text := strings.TrimSpace(string(z.Text()))
	if s.isLikelySignature(text) {
		signals.TopDish = text
	}
}

func (s *WebsiteScraper) handleAnchor(z *html.Tokenizer, base *url.URL, state *scrapeState, signals *AdaSignals) {
	for {
		key, val, more := z.TagAttr()
		if string(key) == "href" {
			link := string(val)
			if s.isMenuLink(link) {
				signals.MenuURL = s.resolveURL(base, link)
			}
			state.inAnchor = true
			state.currentAnchorURL = link
		}
		if !more {
			break
		}
	}
}

func (s *WebsiteScraper) handleText(z *html.Tokenizer, state *scrapeState, base *url.URL, signals *AdaSignals) {
	text := strings.ToLower(string(z.Text()))

	s.checkMenuText(text, state, base, signals)
	s.checkHeatKeywords(text, state)
	s.checkPaymentKeywords(text, state)
	s.checkRegionalKeywords(text, state)
}

func (s *WebsiteScraper) checkMenuText(text string, state *scrapeState, base *url.URL, signals *AdaSignals) {
	if state.inAnchor && state.currentAnchorURL != "" && signals.MenuURL == "" {
		if s.isMenuText(text) {
			signals.MenuURL = s.resolveURL(base, state.currentAnchorURL)
		}
	}
}

func (s *WebsiteScraper) checkHeatKeywords(text string, state *scrapeState) {
	for _, kw := range state.heatKeywords {
		state.heatCount += strings.Count(text, kw)
	}
}

func (s *WebsiteScraper) checkPaymentKeywords(text string, state *scrapeState) {
	for _, kw := range state.paymentKeywords {
		if strings.Contains(text, kw) {
			found := s.capitalizeKeyword(kw)
			if !s.contains(state.foundPayments, found) {
				state.foundPayments = append(state.foundPayments, found)
			}
		}
	}
}

func (s *WebsiteScraper) checkRegionalKeywords(text string, state *scrapeState) {
	for region, keywords := range state.regionalKeywords {
		for _, kw := range keywords {
			state.regionalCounts[region] += strings.Count(text, kw)
		}
	}
}

func (s *WebsiteScraper) mapHeatLevel(count int) int {
	if count > 5 {
		return 5
	}
	return count
}

func (s *WebsiteScraper) isLikelySignature(text string) bool {
	lower := strings.ToLower(text)
	indicators := []string{"signature", "special", "popular", "recommended", "famous", "dish"}
	for _, ind := range indicators {
		if strings.Contains(lower, ind) {
			return true
		}
	}
	return false
}

func (s *WebsiteScraper) isMenuLink(link string) bool {
	lower := strings.ToLower(link)
	indicators := []string{"menu", "order", "glassguide", "doordash", "ubereats", "grubhub", "chownow", "toasttab", "clover"}
	for _, ind := range indicators {
		if strings.Contains(lower, ind) {
			return true
		}
	}
	return strings.HasSuffix(lower, ".pdf") || strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg")
}

func (s *WebsiteScraper) isMenuText(text string) bool {
	lower := strings.ToLower(text)
	indicators := []string{"menu", "order", "order online"}
	for _, ind := range indicators {
		if strings.Contains(lower, ind) {
			return true
		}
	}
	return false
}

func (s *WebsiteScraper) resolveURL(base *url.URL, link string) string {
	if base == nil {
		return link
	}
	u, err := url.Parse(link)
	if err != nil {
		return link
	}
	return base.ResolveReference(u).String()
}

func (s *WebsiteScraper) capitalizeKeyword(kw string) string {
	switch kw {
	case "zelle":
		return "Zelle"
	case "venmo":
		return "Venmo"
	case "cashapp":
		return "CashApp"
	case "cash app":
		return "CashApp"
	default:
		return kw
	}
}

func (s *WebsiteScraper) contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
