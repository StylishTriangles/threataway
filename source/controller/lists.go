package controller

import (
	"encoding/json"
	"gowebapp/source/model/lists"
	"gowebapp/source/model/user"
	"gowebapp/source/view"
	"log"
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

func listsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user has write permissions
	session, err := store.Get(r, "session")
	if err != nil || !session.Values["authenticated"].(bool) {
		http.Error(w, "Please log in to view this page", 403)
		return
	}
	user, ok := session.Values["user"].(*user.User)
	if !ok {
		log.Println("Corrupt cookie data")
		http.Error(w, "An error occoured: corrupt cookie data", 500)
		return
	}
	if user.Role < 1 {
		http.Error(w, "Insufficient permisions (you need write access to perform this action", 403)
	}

	// Delete lists with specified IDs provided via json from frontend
	jsonData := r.FormValue("ids")
	dummy := struct {
		ids []uint32
	}{}
	log.Println(jsonData) // Debug
	json.Unmarshal([]byte(jsonData), &dummy)

}
