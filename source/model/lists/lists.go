package lists

import (
	"gowebapp/source/shared/database"
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
