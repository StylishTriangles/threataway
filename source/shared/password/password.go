package password

import (
	"golang.org/x/crypto/bcrypt"
)

// HashCost - as it may seem
const HashCost = 11

// Hash returns hashed password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), HashCost)
}

// Compare password and hash
// returns true if password matches hash
func Compare(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
