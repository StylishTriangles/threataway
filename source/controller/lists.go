package controller

import (
	"gowebapp/source/model/lists"
	"gowebapp/source/view"
	"net/http"
)

func listsHandler(w http.ResponseWriter, r *http.Request) {
	v := view.New("lists")
	lists, err := lists.GetAllLists()

	if err != nil {
		http.Error(w, "There was an error processing your request"+err.Error(), 500)
	}

	v.Vars["Lists"] = lists
	v.Render(w)
}
