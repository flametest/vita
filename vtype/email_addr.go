package vtype

import (
	"encoding/json"
	"regexp"

	"github.com/flametest/vita/verrors"
	"github.com/flametest/vita/vstring"
)

type EmailAddr struct {
	address string
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(vstring.EmailRegexString)
	if emailRegex.MatchString(email) {
		return true
	}
	return false
}

func NewEmailAddr(address string) (*EmailAddr, error) {
	if !IsValidEmail(address) {
		return nil, verrors.BadRequestError("Invalid email address")
	}
	return &EmailAddr{address: address}, nil
}

func (e *EmailAddr) Address() string {
	return e.address
}

func (e *EmailAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.address)
}

func (e *EmailAddr) UnmarshalJSON(bytes []byte) error {
	var address string
	err := json.Unmarshal(bytes, &address)
	if err != nil {
		return err
	}
	if !IsValidEmail(address) {
		return verrors.BadRequestError("Invalid email address")
	}
	e.address = address
	return nil
}

type EmailAddrList []*EmailAddr

func NewEmailAddrList(addresses []string) (EmailAddrList, error) {
	emailAddrList := make(EmailAddrList, len(addresses))
	for i, address := range addresses {
		emailAddr, err := NewEmailAddr(address)
		if err != nil {
			return nil, err
		}
		emailAddrList[i] = emailAddr
	}
	return emailAddrList, nil
}

func (el *EmailAddrList) Add(emailAddr *EmailAddr) *EmailAddrList {
	*el = append(*el, emailAddr)
	return el
}

func (el *EmailAddrList) Remove(emailAddr *EmailAddr) {
	for i, addr := range *el {
		if addr == emailAddr {
			*el = append((*el)[:i], (*el)[i+1:]...)
			break
		}
	}
}

func (el *EmailAddrList) Addresses() []string {
	addresses := make([]string, len(*el))
	for i, emailAddr := range *el {
		addresses[i] = emailAddr.address
	}
	return addresses
}

func (el *EmailAddrList) Empty() bool {
	return len(*el) == 0
}
