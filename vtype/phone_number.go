package vtype

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/flametest/vita/verrors"
	"github.com/flametest/vita/vstring"
	"github.com/nyaruka/phonenumbers"
)

type PhoneNumber struct {
	countryCode      string
	subscriberNumber string
}

// IsValidPhoneNumber E.164 format: +[country code][subscriber number]
func IsValidPhoneNumber(number string) bool {
	// remove all the spaces first
	number = strings.ReplaceAll(number, " ", "")
	e164Regex := regexp.MustCompile(vstring.E164Regex)
	return e164Regex.MatchString(number)

}

func NewPhoneNumber(number string) (*PhoneNumber, error) {
	pn, err := phonenumbers.Parse(number, "")
	if err != nil {
		return nil, err
	}

	if !phonenumbers.IsValidNumber(pn) {
		return nil, verrors.BadRequestError("Invalid phone number")
	}

	return &PhoneNumber{
		countryCode:      strconv.FormatUint(uint64(*pn.CountryCode), 10),
		subscriberNumber: strconv.FormatUint(*pn.NationalNumber, 10),
	}, nil
}

// NewPhoneNumberWithCountryCode creates a PhoneNumber from separate country code and subscriber number.
// countryCode: e.g., "+1", "+86"
// subscriberNumber: digits only, e.g., "4155552671", "13812345678"
func NewPhoneNumberWithCountryCode(countryCode, subscriberNumber string) (*PhoneNumber, error) {
	return &PhoneNumber{countryCode: countryCode, subscriberNumber: subscriberNumber}, nil
}

func (p *PhoneNumber) SubscriberNumber() string {
	return p.subscriberNumber
}

// E164Format returns the number in E.164 format
func (p *PhoneNumber) E164Format() string {
	return "+" + p.countryCode + p.subscriberNumber
}

func (p *PhoneNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.E164Format())
}

func (p *PhoneNumber) UnmarshalJSON(bytes []byte) error {
	var number string
	err := json.Unmarshal(bytes, &number)
	if err != nil {
		return err
	}
	if !IsValidPhoneNumber(number) {
		return verrors.BadRequestError("Invalid phone number, must be E.164 format (e.g., +14155552671)")
	}
	pn, err := phonenumbers.Parse(number, "")
	if err != nil {
		return err
	}
	p.countryCode = strconv.FormatUint(uint64(*pn.CountryCode), 10)
	p.subscriberNumber = strconv.FormatUint(*pn.NationalNumber, 10)
	return nil
}
