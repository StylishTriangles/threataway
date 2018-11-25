package controller

import (
	"gowebapp/source/model/domain"
	"gowebapp/source/view"
	"net/http"

	"github.com/gorilla/mux"
)

func displayListGET(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}

	mp := mux.Vars(r)

	v := view.New("displayList")
	domains, err := domain.GetFromList(mp["list_name"])
	if err != nil {
		http.Error(w, "There was an error processing your request"+err.Error(), 500)
		return
	}
	v.Vars["Domains"] = domains
	v.Render(w)
}
