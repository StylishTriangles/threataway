package controller

import (
	"gowebapp/source/model/domain"
	"gowebapp/source/model/user"
	"gowebapp/source/view"
	"log"
	"net/http"
	"strconv"
)

func domainsGET(w http.ResponseWriter, r *http.Request) {
	// TODO: SECURITY!!!!
	v := view.New("domains")
	domains, err := domain.GetAll()
	if err != nil {
		http.Error(w, "There was an error processing your request"+err.Error(), 500)
		return
	}
	v.Vars["Domains"] = domains
	v.Render(w)
}

func domainsAdd(w http.ResponseWriter, r *http.Request) {
	// TODO: SECURITY!!!!

}

func domainsCreateList(w http.ResponseWriter, r *http.Request) {
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
		log.Println("Insufficient permisions (you need write access to perform this action")
		http.Error(w, "Insufficient permisions (you need write access to perform this action", 403)
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

	// deleted, err := lists.DeleteLists(parsedIDs, user.ID)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// w.Write([]byte(fmt.Sprintf("Deleted %d records", deleted)))
}
