package ogcard

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"strings"
	"sync"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed assets/logo.png
var logoPNG []byte

const (
	CardWidth  = 1200
	CardHeight = 630
	margin     = 40
)

var (
	colorBg     = color.RGBA{R: 0x1A, G: 0x12, B: 0x0E, A: 0xFF} // #1A120E (earth-dark)
	colorCardBg = color.RGBA{R: 0x24, G: 0x1B, B: 0x16, A: 0xFF} // subtle inner card
	colorBorder = color.RGBA{R: 0x3D, G: 0x2E, B: 0x24, A: 0xFF}
	colorOrange = color.RGBA{R: 0xFF, G: 0x8A, B: 0x00, A: 0xFF} // #FF8A00 (primary)
	colorCream  = color.RGBA{R: 0xFA, G: 0xF8, B: 0xF1, A: 0xFF} // #FAF8F1 (text)
	colorSand   = color.RGBA{R: 0xF4, G: 0xEB, B: 0xD0, A: 0xFF} // #F4EBD0 (subtext)
	colorOchre  = color.RGBA{R: 0xCC, G: 0x77, B: 0x22, A: 0xFF} // #CC7722
	colorMuted  = color.RGBA{R: 0x9E, G: 0x8B, B: 0x7E, A: 0xFF}
)

// CardData encapsulates listing attributes for rendering.
// Struct fields are ordered for 64-bit alignment (fieldalignment).
type CardData struct {
	UpdatedAt   time.Time
	Title       string
	City        string
	Category    string
	Tagline     string
	Rating      float64
	ReviewCount int
}

// Generator defines the contract for rendering OG cards.
type Generator interface {
	Generate(data CardData) ([]byte, error)
	GenerateFallback() []byte
	GenerateBrandCard() []byte
}

// Service renders branded 1200x630 OG share cards with in-process caching.
// Struct fields are ordered for 64-bit alignment (fieldalignment).
type Service struct {
	logoImg      image.Image
	titleFace    font.Face
	subFace      font.Face
	bodyFace     font.Face
	smallFace    font.Face
	cache        map[string][]byte
	fallbackPNG  []byte
	mu           sync.RWMutex
	fallbackOnce sync.Once
}

// NewService initializes fonts, logo asset, and fallback card.
func NewService() *Service {
	s := &Service{
		cache: make(map[string][]byte),
	}

	if len(logoPNG) > 0 {
		if img, err := png.Decode(bytes.NewReader(logoPNG)); err == nil {
			s.logoImg = img
		}
	}

	loadServiceFonts(s)
	s.ensureFallback()
	return s
}

func loadServiceFonts(s *Service) {
	loadBoldFaces(s)
	loadRegularFaces(s)
}

func loadBoldFaces(s *Service) {
	fBold, err := opentype.Parse(gobold.TTF)
	if err != nil {
		return
	}
	if face, faceErr := opentype.NewFace(fBold, &opentype.FaceOptions{
		Size:    76,
		DPI:     72,
		Hinting: font.HintingFull,
	}); faceErr == nil {
		s.titleFace = face
	}
	if face, faceErr := opentype.NewFace(fBold, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingFull,
	}); faceErr == nil {
		s.smallFace = face
	}
}

func loadRegularFaces(s *Service) {
	fReg, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return
	}
	if face, faceErr := opentype.NewFace(fReg, &opentype.FaceOptions{
		Size:    40,
		DPI:     72,
		Hinting: font.HintingFull,
	}); faceErr == nil {
		s.subFace = face
	}
	if face, faceErr := opentype.NewFace(fReg, &opentype.FaceOptions{
		Size:    22,
		DPI:     72,
		Hinting: font.HintingFull,
	}); faceErr == nil {
		s.bodyFace = face
	}
}

// Generate renders or fetches from cache a 1200x630 PNG card.
func (s *Service) Generate(data CardData) ([]byte, error) {
	key := s.cacheKey(data)

	s.mu.RLock()
	if cached, ok := s.cache[key]; ok {
		s.mu.RUnlock()
		return cached, nil
	}
	s.mu.RUnlock()

	pngBytes, err := s.render(data)
	if err != nil {
		return s.GenerateFallback(), nil
	}

	s.mu.Lock()
	if len(s.cache) > 1000 {
		s.cache = make(map[string][]byte)
	}
	s.cache[key] = pngBytes
	s.mu.Unlock()

	return pngBytes, nil
}

// GenerateFallback returns the pre-rendered generic branded card.
func (s *Service) GenerateFallback() []byte {
	s.ensureFallback()
	return s.fallbackPNG
}

// GenerateBrandCard returns the pre-rendered site-wide brand card.
func (s *Service) GenerateBrandCard() []byte {
	return s.GenerateFallback()
}

func (s *Service) ensureFallback() {
	s.fallbackOnce.Do(func() {
		fallbackData := CardData{
			Title:   "agbalumo",
			City:    "Dallas & Beyond",
			Tagline: "Find African food in under 60 seconds.",
		}
		res, err := s.render(fallbackData)
		if err == nil {
			s.fallbackPNG = res
		}
	})
}

func (s *Service) cacheKey(d CardData) string {
	return fmt.Sprintf("%s|%s|%s|%.1f|%d|%d", d.Title, d.City, d.Category, d.Rating, d.ReviewCount, d.UpdatedAt.Unix())
}

func (s *Service) render(data CardData) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, CardWidth, CardHeight))

	draw.Draw(img, img.Bounds(), &image.Uniform{C: colorBg}, image.Point{}, draw.Src)

	cardRect := image.Rect(margin, margin, CardWidth-margin, CardHeight-margin)
	drawRect(img, cardRect, colorCardBg)
	strokeRect(img, cardRect, colorBorder, 2)

	topAccent := image.Rect(margin, margin, CardWidth-margin, margin+6)
	draw.Draw(img, topAccent, &image.Uniform{C: colorOrange}, image.Point{}, draw.Src)

	s.drawHeader(img, data)
	s.drawBody(img, data)
	s.drawFooter(img, data)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Service) drawHeader(img *image.RGBA, data CardData) {
	logoX := margin + 48
	logoY := margin + 36

	if s.logoImg != nil {
		targetH := 48
		bounds := s.logoImg.Bounds()
		targetW := (bounds.Dx() * targetH) / bounds.Dy()
		scaleDraw(img, image.Rect(logoX, logoY, logoX+targetW, logoY+targetH), s.logoImg)
	} else {
		s.drawString(img, s.titleFace, "agbalumo", logoX, logoY+40, colorOrange)
	}

	badgeText := "VERIFIED AFRICAN SPOT"
	if data.Category != "" {
		badgeText = strings.ToUpper(data.Category)
	}
	badgeW := font.MeasureString(s.smallFace, badgeText).Ceil()
	badgeX := CardWidth - 75 - badgeW
	s.drawString(img, s.smallFace, badgeText, badgeX, logoY+32, colorOchre)

	divY := margin + 104
	divRect := image.Rect(logoX, divY, CardWidth-margin-48, divY+1)
	draw.Draw(img, divRect, &image.Uniform{C: colorBorder}, image.Point{}, draw.Src)
}

func (s *Service) drawBody(img *image.RGBA, data CardData) {
	logoX := margin + 48
	title := data.Title
	if title == "" {
		title = "agbalumo"
	}

	maxTitleWidth := CardWidth - 75 - logoX
	lines := s.wrapTitle(title, maxTitleWidth)

	var titleY int
	lineSpacing := 84
	subGap := 68

	if len(lines) == 1 {
		titleY = 320
	} else {
		titleY = 276
	}

	for i, line := range lines {
		s.drawString(img, s.titleFace, line, logoX, titleY+(i*lineSpacing), colorCream)
	}

	metaY := titleY + ((len(lines) - 1) * lineSpacing) + subGap

	locText := data.City
	if locText == "" {
		locText = "Verified Location"
	}

	pinX := logoX + 12
	drawPin(img, pinX, metaY-14, colorOrange)
	cityX := pinX + 22
	s.drawString(img, s.subFace, locText, cityX, metaY, colorSand)

	if data.Rating > 0 {
		cityW := font.MeasureString(s.subFace, locText).Ceil()
		starX := cityX + cityW + 36 + 14
		drawStar(img, starX, metaY-14, 15, 6.5, colorOrange)
		ratingStr := fmt.Sprintf("%.1f", data.Rating)
		if data.ReviewCount > 0 {
			ratingStr += fmt.Sprintf(" (%d reviews)", data.ReviewCount)
		}
		s.drawString(img, s.subFace, ratingStr, starX+24, metaY, colorCream)
	}
}

func (s *Service) drawFooter(img *image.RGBA, data CardData) {
	logoX := margin + 48
	footerDivY := CardHeight - margin - 72
	draw.Draw(img, image.Rect(logoX, footerDivY, CardWidth-margin-48, footerDivY+1), &image.Uniform{C: colorBorder}, image.Point{}, draw.Src)

	tagline := data.Tagline
	if tagline == "" {
		tagline = "Find African food in under 60 seconds."
	}
	footerY := CardHeight - margin - 28
	s.drawString(img, s.bodyFace, tagline, logoX, footerY, colorMuted)
	siteW := font.MeasureString(s.smallFace, "agbalumo.com").Ceil()
	s.drawString(img, s.smallFace, "agbalumo.com", CardWidth-75-siteW, footerY, colorOrange)
}

func (s *Service) drawString(dst *image.RGBA, face font.Face, text string, x, y int, col color.Color) {
	if face == nil {
		return
	}
	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}
	d.DrawString(text)
}

func drawRect(dst *image.RGBA, r image.Rectangle, col color.Color) {
	draw.Draw(dst, r, &image.Uniform{C: col}, image.Point{}, draw.Src)
}

func strokeRect(dst *image.RGBA, r image.Rectangle, col color.Color, strokeWidth int) {
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+strokeWidth), &image.Uniform{C: col}, image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Max.Y-strokeWidth, r.Max.X, r.Max.Y), &image.Uniform{C: col}, image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Min.X+strokeWidth, r.Max.Y), &image.Uniform{C: col}, image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(r.Max.X-strokeWidth, r.Min.Y, r.Max.X, r.Max.Y), &image.Uniform{C: col}, image.Point{}, draw.Src)
}

func scaleDraw(dst *image.RGBA, targetRect image.Rectangle, src image.Image) {
	srcBounds := src.Bounds()
	targetW := targetRect.Dx()
	targetH := targetRect.Dy()
	if targetW <= 0 || targetH <= 0 {
		return
	}

	for y := 0; y < targetH; y++ {
		srcY := srcBounds.Min.Y + (y * srcBounds.Dy() / targetH)
		for x := 0; x < targetW; x++ {
			srcX := srcBounds.Min.X + (x * srcBounds.Dx() / targetW)
			r, g, b, a := src.At(srcX, srcY).RGBA()
			if a > 0 {
				dst.Set(targetRect.Min.X+x, targetRect.Min.Y+y, color.RGBA64{
					R: uint16(r & 0xffff), // #nosec G115 - masked to 16-bit
					G: uint16(g & 0xffff), // #nosec G115 - masked to 16-bit
					B: uint16(b & 0xffff), // #nosec G115 - masked to 16-bit
					A: uint16(a & 0xffff), // #nosec G115 - masked to 16-bit
				})
			}
		}
	}
}

func drawPin(dst *image.RGBA, cx, cy int, col color.Color) {
	for y := -10; y <= 4; y++ {
		for x := -9; x <= 9; x++ {
			distSq := x*x + y*y
			if distSq <= 81 && distSq >= 8 {
				dst.Set(cx+x, cy+y, col)
			}
		}
	}
	for y := 3; y <= 12; y++ {
		halfW := (12 - y) / 2
		for x := -halfW; x <= halfW; x++ {
			dst.Set(cx+x, cy+y, col)
		}
	}
}

func drawStar(dst *image.RGBA, cx, cy int, rOuter, rInner float64, col color.Color) {
	var poly [10][2]float64
	for i := 0; i < 10; i++ {
		angle := -math.Pi/2 + float64(i)*math.Pi/5
		r := rInner
		if i%2 == 0 {
			r = rOuter
		}
		poly[i][0] = float64(cx) + r*math.Cos(angle)
		poly[i][1] = float64(cy) + r*math.Sin(angle)
	}

	minX := int(math.Floor(float64(cx) - rOuter))
	maxX := int(math.Ceil(float64(cx) + rOuter))
	minY := int(math.Floor(float64(cy) - rOuter))
	maxY := int(math.Ceil(float64(cy) + rOuter))

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if pointInPoly(float64(x)+0.5, float64(y)+0.5, poly[:]) {
				dst.Set(x, y, col)
			}
		}
	}
}

func pointInPoly(x, y float64, poly [][2]float64) bool {
	inside := false
	n := len(poly)
	j := n - 1
	for i := 0; i < n; i++ {
		xi, yi := poly[i][0], poly[i][1]
		xj, yj := poly[j][0], poly[j][1]
		if ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}
	return inside
}

func (s *Service) wrapTitle(title string, maxWidth int) []string {
	if s.titleFace == nil || font.MeasureString(s.titleFace, title).Ceil() <= maxWidth {
		return []string{title}
	}

	words := strings.Fields(title)
	if len(words) == 0 {
		return []string{""}
	}

	line1, nextIdx := s.fillFirstLine(words, maxWidth)
	if nextIdx >= len(words) {
		return []string{line1}
	}

	line2 := s.fillSecondLine(words[nextIdx:], maxWidth)
	return []string{line1, line2}
}

func (s *Service) fillFirstLine(words []string, maxWidth int) (string, int) {
	var line string
	for i, w := range words {
		candidate := joinWords(line, w)
		if font.MeasureString(s.titleFace, candidate).Ceil() <= maxWidth {
			line = candidate
			continue
		}
		if line == "" {
			split1, rem := splitLongWord(s.titleFace, w, maxWidth)
			if rem != "" {
				words[i] = rem
				return split1, i
			}
			return split1, i + 1
		}
		return line, i
	}
	return line, len(words)
}

func (s *Service) fillSecondLine(words []string, maxWidth int) string {
	var line string
	for _, w := range words {
		candidate := joinWords(line, w)
		if font.MeasureString(s.titleFace, candidate).Ceil() <= maxWidth {
			line = candidate
			continue
		}
		return truncateToFit(s.titleFace, candidate, maxWidth, "…")
	}
	return line
}

func joinWords(a, b string) string {
	if a == "" {
		return b
	}
	return a + " " + b
}

func splitLongWord(face font.Face, word string, maxWidth int) (string, string) {
	runes := []rune(word)
	for i := len(runes) - 1; i >= 1; i-- {
		part1 := string(runes[:i])
		if font.MeasureString(face, part1).Ceil() <= maxWidth {
			return part1, string(runes[i:])
		}
	}
	return word, ""
}

func truncateToFit(face font.Face, text string, maxWidth int, suffix string) string {
	text = strings.TrimSpace(text)
	if font.MeasureString(face, text+suffix).Ceil() <= maxWidth {
		return text + suffix
	}
	runes := []rune(text)
	for len(runes) > 0 {
		runes = runes[:len(runes)-1]
		candidate := strings.TrimSpace(string(runes)) + suffix
		if font.MeasureString(face, candidate).Ceil() <= maxWidth {
			return candidate
		}
	}
	return suffix
}
