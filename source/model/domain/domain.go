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

func GetFromList(listName string) ([]Domain, error) {
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

// DeleteDomains deletes provided domainIDs from database
func DeleteDomains(domainIDs []uint32) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, v := range domainIDs {
		stmt, err := tx.Prepare(`DELETE FROM urls WHERE idUrl = ?`)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(v)
		if err != nil {
			stmt.Close()
			return err
		}
		stmt.Close()
	}
	return tx.Commit()
}
