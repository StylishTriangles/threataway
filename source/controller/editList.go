package controller

import (
	"gowebapp/source/model/domain"
	"gowebapp/source/model/lists"
	"gowebapp/source/model/user"
	"gowebapp/source/shared/database"
	"gowebapp/source/view"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func editListGET(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}

	mp := mux.Vars(r)
	v := view.New("editList")
	domains, err := domain.GetAll()
	domainsChecked, err := domain.GetFromList(mp["list_name"])

	if err != nil {
		http.Error(w, "There was an error processing your request"+err.Error(), 500)
		return
	}

	v.Vars["Domains"] = domains
	v.Vars["DomainsChecked"] = domainsChecked
	v.Vars["ListName"] = mp["list_name"]
	v.Render(w)
}

func editListPost(w http.ResponseWriter, r *http.Request) {
	// Check if user has write permissions
	session, err := store.Get(r, "session")
	auth, ok := session.Values["authenticated"]
	if err != nil || !ok || !auth.(bool) {
		http.Error(w, "Please log in to view this page", 403)
		log.Println("Please log in to view this page")
		return
	}
	user, ok := session.Values["user"].(*user.User)
	if !ok {
		log.Println("Corrupt cookie data")
		http.Error(w, "An error occoured: corrupt cookie data", 500)
		return
	}
	if user.Role < 1 {
		log.Println("Insufficient permisions (you need write access to perform this action)")
		http.Error(w, "Insufficient permisions (you need write access to perform this action)", 403)
		return
	}

	// Delete lists with specified IDs provided via json from frontend
	r.ParseForm()
	ids, ok := r.Form["ids[]"]
	if !ok {
		http.Error(w, "Invalid request", 500)
		log.Println("Invalid request")
		log.Println(r.Form)
		return
	}
	var parsedIDs []uint32
	for _, sid := range ids {
		val, _ := strconv.Atoi(sid)
		parsedIDs = append(parsedIDs, uint32(val))
	}

	var id uint32
	name := r.FormValue("listName")

	stmt, err := database.DB.Prepare(`SELECT idList FROM lists WHERE name = ?`)

	if err != nil {
		log.Println("Error occured while connecting to database")
		http.Error(w, "Cannot execute name - id relation", 500)

		return
	}
	defer stmt.Close()
	log.Println(name)
	stmt.QueryRow(name).Scan(&id)

	lists.DeleteAllDomainsFromList(id)
	lists.AddDomainsToList(parsedIDs, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Println(err.Error())
		return
	}
	w.Write([]byte("Successfully edited list!"))
}
