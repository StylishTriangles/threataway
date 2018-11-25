package controller

import (
	"encoding/json"
	"fmt"
	"gowebapp/source/model/templates"
	"gowebapp/source/model/user"
	"gowebapp/source/view"
	"log"
	"net/http"
	"strconv"
)

func templatesHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}

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

	r.ParseForm()
	name := r.FormValue("name")
	header := r.FormValue("header")
	footer := r.FormValue("footer")
	urlTemplate := r.FormValue("urlTemplate")
	if name == "" {
		http.Error(w, "Name cannot be empty", 500)
		return
	}

	err = templates.CreateNewTemplate(name, header, footer, urlTemplate)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Println(err.Error())
		return
	}
	w.Write([]byte("Successfully created new template!"))
}

func templateLoad(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}

	tID := r.FormValue("template_id")
	tIDi, err := strconv.Atoi(tID)
	if err != nil {
		http.Error(w, "Could not resolve "+tID+" to integer value", http.StatusBadRequest)
		return
	}
	t, err := templates.GetTemplateByID(uint32(tIDi))
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Printf("templates.go@templateLoad: " + err.Error())
		return
	}

	vals := make(map[string]string)
	vals["name"] = t.Name
	vals["header"] = t.Header
	vals["footer"] = t.Footer
	vals["urlTemplate"] = t.UrlTemplate
	bytes, _ := json.Marshal(vals)
	w.Write(bytes)
}
