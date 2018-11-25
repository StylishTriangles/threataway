package controller

import (
	"gowebapp/source/model/deployments"
	"gowebapp/source/model/lists"
	"gowebapp/source/model/templates"
	"gowebapp/source/view"
	"net/http"
)

func deploymentsHandler(w http.ResponseWriter, r *http.Request) {
	v := view.New("deployments")
	deployments, err := deployments.GetAllDeployments()
	lists, err := lists.GetAllLists()
	templates, err := templates.GetAllTemplates()

	if err != nil {
		http.Error(w, "There was an error processing your request "+err.Error(), 500)
		return
	}

	v.Vars["Deployments"] = deployments
	v.Vars["Lists"] = lists
	v.Vars["Templates"] = templates
	v.Render(w)
}
