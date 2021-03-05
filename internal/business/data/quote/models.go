package quote

// Info represents an individual quote.
type Info struct {
	ID           string   `json:"id"`
	To           Customer `json:"to"`
	From         Customer `json:"from"`
	Weight       int      `json:"weight"`
	ShipmentCost float64  `json:"shipment_cost"`
}

// NewQuote contains information needed to create a new Quote.
type NewQuote struct {
	To     Customer `json:"to" validate:"required,dive"`
	From   Customer `json:"from" validate:"required,dive"`
	Weight int      `json:"weight" validate:"required,gte=0,lte=1000"`
}

// Customer contains information about a customer associated with a quote.
type Customer struct {
	Name        string `json:"name" validate:"required,personname"`
	Email       string `json:"email" validate:"required,email"`
	Address     string `json:"address" validate:"required,max=100"`
	CountryCode string `json:"country_code" validate:"required,iso3166_1_alpha2"`
}
