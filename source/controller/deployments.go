package controller

import (
	"fmt"
	"gowebapp/source/model/deployments"
	"gowebapp/source/model/lists"
	"gowebapp/source/model/templates"
	"gowebapp/source/model/user"
	"gowebapp/source/view"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var urlRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-+]+$`)

func deploymentsHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}

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

func deploymentsAdd(w http.ResponseWriter, r *http.Request) {
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

	lID := r.FormValue("list_id")
	tID := r.FormValue("template_id")
	url := r.FormValue("url")
	url = strings.TrimSpace(url)

	lIDi, err := strconv.Atoi(lID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	tIDi, err := strconv.Atoi(tID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	if urlRegex.FindString(url) == "" {
		http.Error(w, "URL must be in this format: [a-zA-Z0-9_-+]", 400)
		log.Println(url)
		return
	}
	err = deployments.CreateNewDeployment(url, uint32(lIDi), uint32(tIDi))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte(fmt.Sprintf("Success publishing new list, it will be available at /d/%s", url)))
}

func deploymentsDeleteHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid request", 400)
		log.Println("Invalid request")
		log.Println(r.Form)
		return
	}
	var parsedIDs []uint32
	for _, sid := range ids {
		val, _ := strconv.Atoi(sid)
		parsedIDs = append(parsedIDs, uint32(val))
	}

	err = deployments.DeleteDeployments(parsedIDs)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte("Delete success"))
}
