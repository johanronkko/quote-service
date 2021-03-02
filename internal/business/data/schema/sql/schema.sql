-- Version: 1.1
-- Description: Create table quotes
CREATE TABLE quotes (
	quote_id            TEXT,
	package_weight      INT NOT NULL,
    shipment_cost       REAL NOT NULL,
    to_name             TEXT NOT NULL,
    to_email            TEXT NOT NULL,
    to_address          TEXT NOT NULL,
    to_country_code     TEXT NOT NULL,
    from_name           TEXT NOT NULL,
    from_email          TEXT NOT NULL,
    from_address        TEXT NOT NULL,
    from_country_code   TEXT NOT NULL,
	PRIMARY KEY (quote_id)
);
