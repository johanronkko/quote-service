package quote

// ID represents a quote ID.
type ID string

// Info represents an individual quote.
type Info struct {
	ID     `json:"id"`
	To     Customer `json:"to"`
	From   Customer `json:"from"`
	Weight int      `json:"weight"`
	Price  float64  `json:"price"`
}

// NewQuote contains information needed to create a new Quote.
type NewQuote struct {
	To     Customer `json:"to"`
	From   Customer `json:"from"`
	Weight int      `json:"weight"`
}

// Customer represents a customer associated with a quote.
type Customer struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	CountryCode string `json:"country_code"`
}
