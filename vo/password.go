package vo

import (
	"database/sql/driver"

	"github.com/flametest/vita/verrors"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hashedPwd string
}

func NewPassword(pwd string) (*Password, error) {
	//TODO: check the password strength
	hashedPwd, err := hashPwd(pwd)
	if err != nil {
		return nil, err
	}
	return &Password{hashedPwd: hashedPwd}, nil
}

func Generate(length int) string {
	// TODO: generate using different string set defined in vstring
	return ""
}

func hashPwd(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (p *Password) Validate(pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.hashedPwd), []byte(pwd))
}

// Scan implement Scanner
// https://gorm.io/docs/data_types.html
func (p *Password) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return verrors.ErrTypeAssertion
	}
	*p = Password{hashedPwd: string(bytes)}
	return nil
}

// Value implement Valuer
// https://gorm.io/docs/data_types.html
func (p *Password) Value() (driver.Value, error) {
	return p.hashedPwd, nil
}

// MarshalText implement TextMarshaler
func (p Password) MarshalText() ([]byte, error) {
	return []byte("******"), nil
}
