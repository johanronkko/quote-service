package quote

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")
	// ErrNotFound is used when a specific Quote is requested but does not exist.
	ErrNotFound = errors.New("not found")
)

// Quote manages the set of API's for quote access.
type Quote struct {
	db   *sqlx.DB
	calc ShipmentCostCalculator
}

// New constructs a Quote for api access.
func New(db *sqlx.DB, calc ShipmentCostCalculator) Quote {
	return Quote{
		db,
		calc,
	}
}

// Create adds a quote to the database.
func (q Quote) Create(ctx context.Context, nq NewQuote) (Info, error) {

	// TODO: validate new quote

	cost, err := q.calc.ShipmentCost(nq.Weight, nq.From.CountryCode)
	if err != nil {
		return Info{}, fmt.Errorf("calculating shipment cost: %w", err)
	}

	info := Info{
		ID:           generateID(),
		To:           nq.To,
		From:         nq.From,
		Weight:       nq.Weight,
		ShipmentCost: cost,
	}

	const query = `
	INSERT INTO quotes
		(quote_id, package_weight, shipment_cost, to_name, to_email, to_address, to_country_code, from_name, from_email, from_address, from_country_code)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	if _, err := q.db.ExecContext(ctx, query, info.ID, info.Weight, info.ShipmentCost, info.To.Name, info.To.Email, info.To.Address, info.To.CountryCode, info.From.Name, info.From.Email, info.From.Address, info.From.CountryCode); err != nil {
		return Info{}, fmt.Errorf("inserting quote: %w", err)
	}

	return info, nil
}

// Query retrieves a list of existing quotes from the database.
func (q Quote) Query(ctx context.Context) ([]Info, error) {

	const query = `
	SELECT
		*
	FROM
		quotes`

	queryQuotes := []queryQuote{}
	if err := q.db.SelectContext(ctx, &queryQuotes, query); err != nil {
		return nil, fmt.Errorf("selecting quotes: %w", err)
	}

	quotes := []Info{}
	for _, qq := range queryQuotes {
		quotes = append(quotes, qq.toInfo())
	}

	return quotes, nil
}

// QueryByID gets the specified quote from the database.
func (q Quote) QueryByID(ctx context.Context, quoteID ID) (Info, error) {
	if err := quoteID.validate(); err != nil {
		return Info{}, ErrInvalidID
	}

	const query = `
	SELECT
		*
	FROM
		quotes
	WHERE 
		quote_id = $1`

	var queryQuote queryQuote
	if err := q.db.GetContext(ctx, &queryQuote, query, quoteID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, fmt.Errorf("selecting quote %q: %w", quoteID, err)
	}

	return queryQuote.toInfo(), nil
}

type queryQuote struct {
	ID              ID      `db:"quote_id"`
	Weight          int     `db:"package_weight"`
	ShipmentCost    float64 `db:"shipment_cost"`
	ToName          string  `db:"to_name"`
	ToEmail         string  `db:"to_email"`
	ToAddress       string  `db:"to_address"`
	ToCountryCode   string  `db:"to_country_code"`
	FromName        string  `db:"from_name"`
	FromEmail       string  `db:"from_email"`
	FromAddress     string  `db:"from_address"`
	FromCountryCode string  `db:"from_country_code"`
}

func (qq queryQuote) toInfo() Info {
	return Info{
		ID:           qq.ID,
		Weight:       qq.Weight,
		ShipmentCost: qq.ShipmentCost,
		To: Customer{
			Name:        qq.ToName,
			Email:       qq.ToEmail,
			Address:     qq.ToAddress,
			CountryCode: qq.ToCountryCode,
		},
		From: Customer{
			Name:        qq.FromName,
			Email:       qq.FromEmail,
			Address:     qq.FromAddress,
			CountryCode: qq.FromCountryCode,
		},
	}
}

func (id ID) validate() error {
	if _, err := uuid.Parse(string(id)); err != nil {
		return ErrInvalidID
	}
	return nil
}

func generateID() ID {
	return ID(uuid.New().String())
}
