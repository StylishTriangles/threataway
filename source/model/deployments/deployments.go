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
