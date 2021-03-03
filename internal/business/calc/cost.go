package calc

import (
	"errors"
	"strings"
)

// ShipmentCost calculates shipment costs.
type ShipmentCost struct {
	regions map[string]region
}

// NewShipmentCost constructs a ShipmentCost for shipment cost calculation.
func NewShipmentCost() ShipmentCost {
	// ShipmentCost is currently done from a hardcode map with country code as key
	// and region as value. In the future this could instead be done via an API
	// call to a centralized database, externally or internally.
	regions := map[string]region{
		"sv": nordic, // Sweden
		"no": nordic, // Norway
		"dk": nordic, // Denmark
		"fi": nordic, // Finland

		"fr": withinEU, // France
		"de": withinEU, // Germany
		"nl": withinEU, // Netherlands
		"it": withinEU, // Italy
		"pt": withinEU, // Portugal
		"at": withinEU, // Austria
		"be": withinEU, // Belgium
		"lv": withinEU, // Latvia
		"bg": withinEU, // Bulgaria
		"lt": withinEU, // Lithuania
		"hr": withinEU, // Croatia
		"lu": withinEU, // Luxembourg
		"cy": withinEU, // Cyprus
		"mt": withinEU, // Malta
		"cz": withinEU, // Czechia
		"pl": withinEU, // Poland
		"ee": withinEU, // Estonia
		"ro": withinEU, // Romania
		"sk": withinEU, // Slovakia
		"si": withinEU, // Slovenia
		"gr": withinEU, // Greece
		"es": withinEU, // Spain
		"hu": withinEU, // Hungary
		"ie": withinEU, // Ireland

		"us": outsideEU, // United States of America
		"ca": outsideEU, // Canada
		"cn": outsideEU, // China
		"jp": outsideEU, // Japan
		"th": outsideEU, // Thailand
		"br": outsideEU, // Brazil
		"ar": outsideEU, // Argentina
	}
	return ShipmentCost{regions}
}

// ShipmentCost calculates shipment cost as the multiplication of a package's
// weight class factor and a region factor determined by the country code.
// Errors if country code not supported or if package not within a valid weight
// class.
//
// We have four classes of weight: small (0 - 10kg), 100sek; medium (10 - 25kg),
// 300sek; large (25 - 50kg), 500sek; huge (50 - 1000kg), 2000sek. If country
// code is Nordic, weight class is multiplied with 1, within EU with 1.5 and
// outside EU with 2.5.
func (c *ShipmentCost) ShipmentCost(weight int, ccode string) (float64, error) {
	region, ok := c.regions[strings.ToLower(ccode)]
	if !ok {
		return 0, errors.New("unsupported country code")
	}
	if weight < 0 || weight > 1000 {
		return 0, errors.New("invalid weight")
	}
	regionFactor := float64(region)
	if weight <= 10 {
		return 100 * regionFactor, nil
	}
	if weight <= 25 {
		return 300 * regionFactor, nil
	}
	if weight <= 50 {
		return 500 * regionFactor, nil
	}
	return 2000 * regionFactor, nil
}

type region float64

const (
	nordic    region = 1
	withinEU         = 1.5
	outsideEU        = 2.5
)
