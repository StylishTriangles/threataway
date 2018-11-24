package controller

import (
	"gowebapp/source/model/lists"
	"gowebapp/source/model/user"
	"gowebapp/source/view"
	"log"
	"net/http"
	"strconv"
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
	auth, ok := session.Values["authenticated"]
	if err != nil || !ok || !auth.(bool) {
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
	r.ParseForm()
	ids, ok := r.Form["ids"]
	if !ok {
		http.Error(w, "Invalid request", 500)
	}
	var parsedIDs []uint32
	for _, sid := range ids {
		val, _ := strconv.Atoi(sid)
		parsedIDs = append(parsedIDs, uint32(val))
	}
}
