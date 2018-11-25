package controller

import (
	"fmt"
	"gowebapp/source/model/templates"
	"gowebapp/source/model/user"
	"gowebapp/source/view"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func templatesHandler(w http.ResponseWriter, r *http.Request) {
	v := view.New("templates")
	templates, err := templates.GetAllTemplates()

	if err != nil {
		http.Error(w, "There was an error processing your request "+err.Error(), 500)
		return
	}

	v.Vars["Templates"] = templates
	v.Render(w)
}

func templatesDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	// Delete templates with specified IDs provided via json from frontend
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

	deleted, err := templates.DeleteTemplates(parsedIDs)
	if err != nil {
		log.Println(err.Error())
	}
	w.Write([]byte(fmt.Sprintf("Deleted %d records", deleted)))
}

func templateCreate(w http.ResponseWriter, r *http.Request) {
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
	name, ok_n := r.Form["name"]
	header, ok_h := r.Form["header"]
	footer, ok_f := r.Form["footer"]
	urlTemplate, ok_u := r.Form["urlTemplate"]
	if !(ok_n && ok_h && ok_f && ok_u) || name[0] == "" || !strings.Contains(urlTemplate[0], "{URL}") {
		http.Error(w, "Invalid request", 500)
		log.Println("Invalid request")
		log.Println(r.Form)
		return
	}
	err = templates.CreateNewTemplate(name[0], header[0], footer[0], urlTemplate[0])
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Println(err.Error())
		return
	}
	w.Write([]byte("Successfully created new template!"))

}
