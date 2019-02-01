package templates

import "threataway/source/shared/database"

// Template may contain one row from templates table
type Template struct {
	ID          uint32 `db:"templateID"`
	Name        string `db:"name"`
	Header      string `db:"header"`
	Footer      string `db:"footer"`
	UrlTemplate string `db:"urlTemplate"`
}

// New returns newly created template pointer
func New() *Template {
	return &Template{}
}

// GetAllTemplates fetches all templates from DB
func GetAllTemplates() ([]Template, error) {
	stmt, err := database.DB.Prepare(`SELECT templateID, name, header, footer, urlTemplate FROM templates`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var returnList []Template

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		template := New()
		rows.Scan(&template.ID, &template.Name, &template.Header, &template.Footer, &template.UrlTemplate)

		returnList = append(returnList, *template)
	}

	return returnList, nil

}

// GetTemplateByID returns template which has given ID
func GetTemplateByID(ID uint32) (Template, error) {
	ret := Template{}
	stmt, err := database.DB.Prepare("SELECT templateID, name, header, footer, urlTemplate FROM templates WHERE templateID = ?")
	if err != nil {
		return ret, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(ID)
	err = row.Scan(&ret.ID, &ret.Name, &ret.Header, &ret.Footer, &ret.UrlTemplate)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

// DeleteTemplates deletes provided templateIDs from DB
func DeleteTemplates(templateIDs []uint32) (int, error) {
	deleted := 0
	for _, v := range templateIDs {
		stmt, err := database.DB.Prepare(`DELETE FROM templates WHERE templateID = ?`)
		if err != nil {
			return deleted, err
		}
		res, err := stmt.Exec(v)
		if err != nil {
			stmt.Close()
			return deleted, err
		}
		v, _ := res.RowsAffected()
		deleted += int(v)
		stmt.Close()
	}
	return deleted, nil
}

// CreateNewTemplate creates new template in db with ids of given domains, specific name and ownerID also must be provided
func CreateNewTemplate(name, header, footer, urlTemplate string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO templates(name, header, footer, urlTemplate) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(name, header, footer, urlTemplate); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateTemplate overrides template in database
func UpdateTemplate(updatedTemplate Template) error {
	ut := updatedTemplate
	stmt, err := database.DB.Prepare("UPDATE templates SET name = ?, header = ?, footer = ?, urlTemplate = ? WHERE templateID = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ut.Name, ut.Header, ut.Footer, ut.UrlTemplate, ut.ID)
	return err
}
