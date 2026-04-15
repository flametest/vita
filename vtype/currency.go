package vtype

import (
	"encoding/json"

	"github.com/flametest/vita/verrors"
)

type Currency string

const (
	// Major reserve currencies
	CurrencyUSD Currency = "USD" // US Dollar
	CurrencyEUR Currency = "EUR" // Euro
	CurrencyGBP Currency = "GBP" // British Pound Sterling
	CurrencyJPY Currency = "JPY" // Japanese Yen
	CurrencyCHF Currency = "CHF" // Swiss Franc
	CurrencyCAD Currency = "CAD" // Canadian Dollar
	CurrencyAUD Currency = "AUD" // Australian Dollar
	CurrencyNZD Currency = "NZD" // New Zealand Dollar

	// Asia-Pacific
	CurrencyCNY Currency = "CNY" // Chinese Yuan
	CurrencyHKD Currency = "HKD" // Hong Kong Dollar
	CurrencySGD Currency = "SGD" // Singapore Dollar
	CurrencyKRW Currency = "KRW" // South Korean Won
	CurrencyTWD Currency = "TWD" // New Taiwan Dollar
	CurrencyTHB Currency = "THB" // Thai Baht
	CurrencyMYR Currency = "MYR" // Malaysian Ringgit
	CurrencyPHP Currency = "PHP" // Philippine Peso
	CurrencyIDR Currency = "IDR" // Indonesian Rupiah
	CurrencyVND Currency = "VND" // Vietnamese Dong
	CurrencyINR Currency = "INR" // Indian Rupee

	// Middle East
	CurrencyAED Currency = "AED" // UAE Dirham
	CurrencySAR Currency = "SAR" // Saudi Riyal

	// Europe
	CurrencySEK Currency = "SEK" // Swedish Krona
	CurrencyNOK Currency = "NOK" // Norwegian Krone
	CurrencyDKK Currency = "DKK" // Danish Krone
	CurrencyPLN Currency = "PLN" // Polish Zloty
	CurrencyTRY Currency = "TRY" // Turkish Lira

	// Americas
	CurrencyMXN Currency = "MXN" // Mexican Peso
	CurrencyBRL Currency = "BRL" // Brazilian Real

	// Africa
	CurrencyZAR Currency = "ZAR" // South African Rand

	// Eurasia
	CurrencyRUB Currency = "RUB" // Russian Ruble
)

var supportedCurrencies = map[Currency]bool{
	// Major reserve currencies
	CurrencyUSD: true,
	CurrencyEUR: true,
	CurrencyGBP: true,
	CurrencyJPY: true,
	CurrencyCHF: true,
	CurrencyCAD: true,
	CurrencyAUD: true,
	CurrencyNZD: true,
	// Asia-Pacific
	CurrencyCNY: true,
	CurrencyHKD: true,
	CurrencySGD: true,
	CurrencyKRW: true,
	CurrencyTWD: true,
	CurrencyTHB: true,
	CurrencyMYR: true,
	CurrencyPHP: true,
	CurrencyIDR: true,
	CurrencyVND: true,
	CurrencyINR: true,
	// Middle East
	CurrencyAED: true,
	CurrencySAR: true,
	// Europe
	CurrencySEK: true,
	CurrencyNOK: true,
	CurrencyDKK: true,
	CurrencyPLN: true,
	CurrencyTRY: true,
	// Americas
	CurrencyMXN: true,
	CurrencyBRL: true,
	// Africa
	CurrencyZAR: true,
	// Eurasia
	CurrencyRUB: true,
}

func NewCurrency(s string) (Currency, error) {
	c := Currency(s)
	if !c.IsValid() {
		return "", verrors.BadRequestError("unsupported currency: " + s)
	}
	return c, nil
}

func (c Currency) IsValid() bool {
	return supportedCurrencies[c]
}

func (c Currency) String() string {
	return string(c)
}

// UnmarshalJSON is a hook for json.Unmarshal. When JSON is deserialized into a
// Currency field, Go calls this method instead of the default string assignment.
//
// Without this method, invalid currency values would pass through silently.
// With it, any JSON deserialization entry point (HTTP request body, config file, etc.)
// will automatically reject unsupported currencies.
//
// Example:
//
//	type Order struct {
//	    Currency vtype.Currency `json:"currency"`
//	}
//
//	var o Order
//	json.Unmarshal([]byte(`{"currency": "CNY"}`), &o)
//	// Without UnmarshalJSON: o.Currency = "CNY", no error
//	// With UnmarshalJSON:    returns BadRequestError
func (c *Currency) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	v, err := NewCurrency(s)
	if err != nil {
		return err
	}
	*c = v
	return nil
}
