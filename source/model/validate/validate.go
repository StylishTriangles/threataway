package validate

import (
	"bufio"
	"errors"
	"log"
	"os"
	"regexp"
)

// all lengths are in bytes
const (
	UsernameMaxLen = 32
	EmailMaxLen    = 255
	NameMaxLen     = 45
	PasswordMinLen = 8
	PasswordMaxLen = 56 // roughly equal to bcrypt's max key length
)

var (
	emailRegexp    = regexp.MustCompile("^[^\\s@]+@(?:[^\\s@.,]+\\.)+[^\\s@.,]{2,}$") // Minimalist, quick reqex (can match many incorrect emails though)
	usernameRegexp = regexp.MustCompile("^[A-Za-z0-9._-]+$")

	// a set of common passwords to speed up lookup
	commonPasswords map[string]struct{}

	// path to common passwords list
	commonPasswordsPath = "misc/passwords/common.txt"
)

var (
	// ErrCommonPassword is returned when password is in a list of common passwords
	ErrCommonPassword = errors.New("Input is a common password")
	// ErrEmpty is returned when input is an empty string
	ErrEmpty = errors.New("Input is empty")
	// ErrInvalidFormat is returned when input doesn't match a regex
	ErrInvalidFormat = errors.New("Input has invalid format")
	// ErrTooLong is returned when input has too many bytes
	ErrTooLong = errors.New("Input too long")
	// ErrTooShort is returned when input has too few bytes
	ErrTooShort = errors.New("Input too short")
)

// Username is used to validate usernames.
// Returns true if username meets hardcoded criteria.
// Otherwise an error is returned.
func Username(username string) (bool, error) {
	// check length
	l := len(username)
	if l == 0 {
		return false, ErrEmpty
	}
	if l > UsernameMaxLen {
		return false, ErrTooLong
	}
	ok := usernameRegexp.MatchString(username)
	if !ok {
		return false, ErrInvalidFormat
	}
	return true, nil
}

// Email is used to check for common typing errors, it does not however check if e-mail address exists.
// Returns true if input string is in the format _example@example.com_ and doesn't exceed maxlen.
// Otherwise an error is returned.
func Email(email string) (bool, error) {
	// check length
	l := len(email)
	if l == 0 {
		return false, ErrEmpty
	}
	if l > EmailMaxLen {
		return false, ErrTooLong
	}
	ok := emailRegexp.MatchString(email)
	if !ok {
		return false, ErrInvalidFormat
	}
	return true, nil
}

// Name is used to validate person's first/middle/last name.
// Returns true if name is non zero length and shorter than maxlen.
// Otherwise an error is returned.
func Name(name string) (bool, error) {
	l := len(name)
	if l == 0 {
		return false, ErrEmpty
	}
	if l > NameMaxLen {
		return false, ErrTooLong
	}
	return true, nil
}

// Password is used to validate user's password.
// Returns true if password meets specific criteria.
// Otherwise an error is returned.
func Password(password string) (bool, error) {
	l := len(password)
	if l == 0 {
		return false, ErrEmpty
	}
	if l < PasswordMinLen {
		return false, ErrTooShort
	}
	if l > PasswordMaxLen {
		return false, ErrTooLong
	}

	isCommon, err := isPasswordCommon(password)
	if err != nil {
		log.Fatal(err)
	}
	if isCommon {
		return false, ErrCommonPassword
	}
	return true, nil
}

// isPasswordCommon checks if specified password is in common passwords list.
// If the file containing those passwords cannot be opened an error is returned.
func isPasswordCommon(password string) (bool, error) {
	if commonPasswords == nil {
		commonPasswords = make(map[string]struct{})
		// load passwords from file
		file, err := os.Open(commonPasswordsPath)
		if err != nil {
			return false, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			commonPasswords[scanner.Text()] = struct{}{}
		}

		if err := scanner.Err(); err != nil {
			return false, err
		}
	}
	_, exists := commonPasswords[password]
	return exists, nil
}
