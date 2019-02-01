package lists

import (
	"threataway/source/shared/database"
	"log"
	"sort"
)

// List may contain one row from lists table
type List struct {
	ID    uint32 `db:"idList"`
	Name  string `db:"name"`
	Owner uint32 `db:"ownerID"`
}

// New creates new user
func New() *List {
	return &List{}
}

// GetAllLists does smth
func GetAllLists() ([]List, error) {

	stmt, err := database.DB.Prepare(`SELECT idList, name, ownerID FROM lists`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var returnList []List

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		list := New()
		rows.Scan(&list.ID, &list.Name, &list.Owner)

		returnList = append(returnList, *list)
	}

	return returnList, nil

}

// DeleteLists deletes provided list IDs from DB
func DeleteLists(listIDs []uint32, userID uint32) (int, error) {
	sort.Slice(listIDs, func(i, j int) bool { return listIDs[i] < listIDs[j] })
	deleted := 0
	tx, err := database.DB.Begin()
	if err != nil {
		return deleted, err
	}
	defer tx.Rollback()
	for _, v := range listIDs {
		stmt, err := tx.Prepare(`DELETE FROM lists WHERE idList = ? AND ownerID = ?`)
		if err != nil {
			return deleted, err
		}
		res, err := stmt.Exec(v, userID)
		if err != nil {
			stmt.Close()
			return deleted, err
		}
		v, _ := res.RowsAffected()
		deleted += int(v)
		stmt.Close()
	}
	return deleted, tx.Commit()
}

// DeleteAllDomainsFromList clears domains in list
func DeleteAllDomainsFromList(listID uint32) {
	stmt, err := database.DB.Prepare(`DELETE from listlists WHERE idList = ?`)
	if err != nil {
		log.Println("Cannot delete urls from list")

		return
	}
	defer stmt.Close()

	stmt.Query(listID)
}

// AddDomainsToList adds domains to existing list
func AddDomainsToList(domainIDs []uint32, listID uint32) error {

	for _, v := range domainIDs {
		stmt, err := database.DB.Prepare("INSERT INTO listlists(idUrl, idList) VALUES(?, ?)")
		if err != nil {
			return err
		}

		_, err = stmt.Query(v, listID)
		if err != nil {
			log.Println("Cannot insert id to database")
			stmt.Close()
			return err
		}
		stmt.Close()
	}

	return nil
}

// CreateNewList creates new list in db with ids of given domains, specific name and ownerID also must be provided
func CreateNewList(domainIDs []uint32, name string, userID uint32) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO lists(name, ownerID) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, userID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	idu := uint32(id)
	for _, v := range domainIDs {
		stmt2, err := tx.Prepare("INSERT INTO listlists(idUrl, idList) VALUES(?, ?)")
		if err != nil {
			return err
		}

		_, err = stmt2.Exec(v, idu)
		if err != nil {
			stmt2.Close()
			return err
		}
		stmt2.Close()
	}
	return tx.Commit()
}
