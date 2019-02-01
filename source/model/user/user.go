package user

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"threataway/source/shared/database"
	"threataway/source/shared/email"
	"threataway/source/shared/password"
	"time"
)

var (
	// ErrAlreadyRegistered is user already registered error
	ErrAlreadyRegistered = errors.New("Username or email already in use")
	// ErrNotFound is not found in database error
	ErrNotFound = errors.New("User not found")
	// ErrInvalidActivator is invalid activation key error
	ErrInvalidActivator = errors.New("Invalid activation key")
	// ErrActivatorExpired is an expired activation key error
	ErrActivatorExpired = errors.New("Activation key expired")
	// ErrInvalidEmail is an invalid email address error
	ErrInvalidEmail = errors.New("Invalid email address")
)

// User may contain one row from users table
type User struct {
	ID           uint32    `db:"idUser"`
	Login        string    `db:"login"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password"`
	FirstName    string    `db:"firstName"`
	LastName     string    `db:"lastName"`
	LastLogin    time.Time `db:"lastLogin"`
	Active       uint8     `db:""`
	Role         uint8     `db:""`
}

// New creates new user
func New() *User {
	return &User{}
}

type confirm struct {
	IDUser  uint32
	Key     string
	Expires time.Time
	Used    uint8
}

func init() {
	gob.Register(&User{})
}

// Confirm confirms certain action bound by an activation key
func Confirm(key string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Search for activation key in db
	stmt, err := tx.Prepare("SELECT `idUser`, `expires`, `used` FROM `confirm` WHERE `key` = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	c := confirm{}
	var exp string
	err = stmt.QueryRow(key).Scan(&c.IDUser, &exp, &c.Used)
	c.Expires, err = time.Parse("2006-01-02 15:04:05", exp)
	if err != nil {
		return err
	}
	// Check if provided key is valid
	if err != nil || c.Used != 0 {
		return ErrInvalidActivator
	}
	// Check if key hasn't expired yet
	if t := time.Now(); t.After(c.Expires) {
		return ErrActivatorExpired
	}

	// Activate user account
	stmt2, err := tx.Prepare("UPDATE `users` SET `active` = '1' WHERE `idUser` = ?;")
	if err != nil {
		return err
	}
	defer stmt2.Close()
	_, err = stmt2.Exec(c.IDUser)
	if err != nil {
		return err
	}
	// Mark activation key as used
	stmt3, err := tx.Prepare("UPDATE `confirm` SET `used` = '1' WHERE `key` = ?;")
	if err != nil {
		return err
	}
	defer stmt3.Close()
	_, err = stmt3.Exec(key)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func createActivationKey() (string, error) {
	bin := make([]byte, 32)
	_, err := rand.Read(bin)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bin), nil
}

func sendActivationEmail(address, key string) error {
	e := email.New()
	// TODO: make configurable
	al := "https://legoracer.ddns.net/account/activate?key=" + key
	e.SetRecipient(address)
	e.SetSubject("Action required: activate your account")
	e.SetTemplate("verify_email.html")
	e.SetVar("ActivationLink", al)
	return e.Send()
}

// Register creates new database record for a user, including valid records in referenced tables.
// It also assumes entered data is formatted correctly.
// Returns nil when function succeeds
func Register(login, email, passwd string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if username already exists in db
	stmt, err := tx.Prepare("SELECT 1 FROM users WHERE login = ? OR email = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var exists bool
	err = stmt.QueryRow(login, email).Scan(&exists)
	if err == nil {
		return ErrAlreadyRegistered
	}

	// Generate password hash
	hash, err := password.Hash(passwd)
	if err != nil {
		return err
	}

	// get registration time
	tm := time.Now()
	// Insert new user into DB
	// IMPORTANT! Accounts will always start activated
	stmt2, err := tx.Prepare("INSERT INTO users(login, email, password, createDate, active, role) VALUES(?, ?, ?, ?, 1, 0)")
	if err != nil {
		return err
	}
	defer stmt2.Close()
	res, err := stmt2.Exec(login, email, hash, tm)
	if err != nil {
		return err
	}
	uid, err := res.LastInsertId() // get user id
	if err != nil {
		return err
	}
	// Create activation key
	key, err := createActivationKey()
	if err != nil {
		return err
	}
	// Put activation key in database
	stmtConfirm, err := tx.Prepare("INSERT INTO `confirm` (`idUser`, `key`) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmtConfirm.Close()
	_, err = stmtConfirm.Exec(uid, key)
	if err != nil {
		return err
	}
	// Send activation email
	err = sendActivationEmail(email, key)
	// IMPORTANT! Doesn't matter if email is sent because accounts are always activated
	// if err != nil {
	// 	return ErrInvalidEmail
	// }

	return tx.Commit()
}

// GetByHandle returns a *User struct of a user with specified email or username.
// An error (ErrNotFound) is returned, when user was not found in DB
func GetByHandle(emailOrLogin string) (*User, error) {
	stmt, err := database.DB.Prepare(`
SELECT idUser, login, email, password, firstName, lastName, active, role
FROM users
WHERE email = ? or login = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	u := New()
	err = stmt.QueryRow(emailOrLogin, emailOrLogin).
		Scan(&u.ID, &u.Login, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Active, &u.Role)
	if err != nil {
		return nil, ErrNotFound
	}
	return u, nil
}
