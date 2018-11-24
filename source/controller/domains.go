package controller

import (
	"gowebapp/source/model/domain"
	"gowebapp/source/view"
	"net/http"
)

func domainsGET(w http.ResponseWriter, r *http.Request) {
	v := view.New("domains")
	domains, err := domain.GetAll()
	if err != nil {
		http.Error(w, "There was an error processing your request"+err.Error(), 500)
	}
	v.Vars["Domains"] = domains
	v.Render()
}
