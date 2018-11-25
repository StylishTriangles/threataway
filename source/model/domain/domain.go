package domain

import (
	"gowebapp/source/shared/database"
)

// Domain represents a single domain in database
type Domain struct {
	ID     uint32
	URL    string
	Rating float32
}

// GetAll returns list of all tracked domains
func GetAll() ([]Domain, error) {
	var ret []Domain
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check if username already exists in db
	stmt, err := tx.Prepare("SELECT idUrl, domain, rating FROM urls")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		d := Domain{}
		err = rows.Scan(&d.ID, &d.URL, &d.Rating)
		if err != nil {
			return nil, err
		}
		ret = append(ret, d)
	}

	tx.Commit()
	return ret, nil
}

// RegisterNew creates a new database record for given domain
func RegisterNew(domain string) error {
	stmt, err := database.DB.Prepare("INSERT INTO urls(domain) VALUES(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(domain)

	return err
}
