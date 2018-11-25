package deployments

import (
	"gowebapp/source/shared/database"
)

// Deployment may contain one row from lists table
type Deployment struct {
	ID           uint32 `db:"deploymentID"`
	URL          string `db:"url"`
	ListID       uint32 `db:"listID"`
	ListName     string `db:"listName"`
	TemplateID   uint32 `db:"templateID"`
	TemplateName string `db:"templateName"`
}

// New creates new Deployments
func New() *Deployment {
	return &Deployment{}
}

// GetAllDeployments does smth
func GetAllDeployments() ([]Deployment, error) {

	stmt, err := database.DB.Prepare(`SELECT templates.name, lists.name, deployments.deploymentID, deployments.url, deployments.listID, deployments.templateID FROM deployments 
	INNER JOIN lists ON lists.idList = deployments.listID INNER JOIN templates ON templates.templateID = deployments.templateID`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var returnList []Deployment

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		deployment := New()
		rows.Scan(&deployment.TemplateName, &deployment.ListName, &deployment.ID, &deployment.URL, &deployment.ListID, &deployment.TemplateID)

		returnList = append(returnList, *deployment)
	}

	return returnList, nil

}

// CreateNewDeployment creates new deployment in the database
func CreateNewDeployment(deployURL string, listID, tmplID uint32) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmti, err := tx.Prepare("INSERT INTO deployments(url, listID, templateID) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmti.Close()
	_, err = stmti.Exec(deployURL, listID, tmplID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// DeleteDeployments does exactly what it says
func DeleteDeployments(deploymentIDs []uint32) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, v := range deploymentIDs {
		stmt, err := tx.Prepare(`DELETE FROM deployments WHERE deploymentID = ?`)
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
