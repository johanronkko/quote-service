package region

import (
	"errors"
	"strings"
)

var (
	// ErrUnsupportedCountryCode occurs when provided country code not supported.
	// Country codes are expected to follow the ISO-3166-1 alpha-2 standard.
	ErrUnsupportedCountryCode = errors.New("country code not supported")
)

// Region represents a region in the world. Value of region corresponds to its
// factor value when calculating shipment cost.
type Region float64

// Regions.
const (
	Nordic    Region = 1
	WithinEU         = 1.5
	OutsideEU        = 2.5
)

// From returns region given country code. Errors if invalid country code.
func From(ccode string) (Region, error) {
	r, ok := regions[strings.ToLower(ccode)]
	if !ok {
		return 0, ErrUnsupportedCountryCode
	}
	return r, nil
}

// Regions are hard coded. In the future we could instead do the country code to
// region mapping via an API call to a centralized database, externally or
// internally.
var regions map[string]Region

func init() {
	regions = map[string]Region{
		"sv": Nordic, // Sweden
		"no": Nordic, // Norway
		"dk": Nordic, // Denmark
		"fi": Nordic, // Finland

		"fr": WithinEU, // France
		"de": WithinEU, // Germany
		"nl": WithinEU, // Netherlands
		"it": WithinEU, // Italy
		"pt": WithinEU, // Portugal
		"at": WithinEU, // Austria
		"be": WithinEU, // Belgium
		"lv": WithinEU, // Latvia
		"bg": WithinEU, // Bulgaria
		"lt": WithinEU, // Lithuania
		"hr": WithinEU, // Croatia
		"lu": WithinEU, // Luxembourg
		"cy": WithinEU, // Cyprus
		"mt": WithinEU, // Malta
		"cz": WithinEU, // Czechia
		"pl": WithinEU, // Poland
		"ee": WithinEU, // Estonia
		"ro": WithinEU, // Romania
		"sk": WithinEU, // Slovakia
		"si": WithinEU, // Slovenia
		"gr": WithinEU, // Greece
		"es": WithinEU, // Spain
		"hu": WithinEU, // Hungary
		"ie": WithinEU, // Ireland

		"us": OutsideEU, // United States of America
		"ca": OutsideEU, // Canada
		"cn": OutsideEU, // China
		"jp": OutsideEU, // Japan
		"th": OutsideEU, // Thailand
		"br": OutsideEU, // Brazil
		"ar": OutsideEU, // Argentina
	}
}
