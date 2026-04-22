package domain

var ValidOrigins = map[string]bool{
	// West Africa
	"Benin":         true,
	"Burkina Faso":  true,
	"Cabo Verde":    true,
	"Cote d'Ivoire": true,
	"Gambia":        true,
	"Ghana":         true,
	"Guinea":        true,
	"Guinea-Bissau": true,
	"Liberia":       true,
	"Mali":          true,
	"Mauritania":    true,
	"Niger":         true,
	"Nigeria":       true,
	"Senegal":       true,
	"Sierra Leone":  true,
	"Togo":          true,
	// North Africa
	"Algeria":        true,
	"Egypt":          true,
	"Libya":          true,
	"Morocco":        true,
	"Sudan":          true,
	"Tunisia":        true,
	"Western Sahara": true,
	// East Africa
	"Burundi":     true,
	"Comoros":     true,
	"Djibouti":    true,
	"Eritrea":     true,
	"Ethiopia":    true,
	"Kenya":       true,
	"Madagascar":  true,
	"Malawi":      true,
	"Mauritius":   true,
	"Mozambique":  true,
	"Rwanda":      true,
	"Seychelles":  true,
	"Somalia":     true,
	"South Sudan": true,
	"Tanzania":    true,
	"Uganda":      true,
	"Zambia":      true,
	"Zimbabwe":    true,
	// Central Africa
	"Angola":                           true,
	"Cameroon":                         true,
	"Central African Republic":         true,
	"Chad":                             true,
	"Congo":                            true,
	"Democratic Republic of the Congo": true,
	"Equatorial Guinea":                true,
	"Gabon":                            true,
	"Sao Tome and Principe":            true,
	// Southern Africa
	"Botswana":     true,
	"Eswatini":     true,
	"Lesotho":      true,
	"Namibia":      true,
	"South Africa": true,
	// Other
	"Other": true,
}

func (l *Listing) validateOrigin() error {
	if l.OwnerOrigin == "" {
		return ErrMissingOrigin
	}
	if !ValidOrigins[l.OwnerOrigin] {
		return ErrInvalidOrigin
	}
	return nil
}
