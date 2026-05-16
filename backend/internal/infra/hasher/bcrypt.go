// Package hasher provides the production PasswordHasher implementation.
package hasher

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Bcrypt implements usecase.PasswordHasher using a fixed bcrypt cost.
type Bcrypt struct{ Cost int }

// New returns a Bcrypt hasher. Pass 0 for default cost (bcrypt.DefaultCost).
func New(cost int) Bcrypt {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return Bcrypt{Cost: cost}
}

// Hash returns the bcrypt hash of plain.
func (b Bcrypt) Hash(plain string) (string, error) {
	out, err := bcrypt.GenerateFromPassword([]byte(plain), b.Cost)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Verify reports nil when plain matches the stored hash.
func (b Bcrypt) Verify(hash, plain string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)); err != nil {
		return errors.New("password mismatch")
	}
	return nil
}
