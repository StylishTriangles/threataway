package controller

import (
	"fmt"
	"threataway/source/model/domain"
	"threataway/source/model/lists"
	"threataway/source/model/user"
	"threataway/source/shared/database"
	"threataway/source/view"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

var ipRegex = regexp.MustCompile(`^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4[0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
var ipShort = regexp.MustCompile(`^(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])$`)
var illegalDomainChars = regexp.MustCompile("[;/?:@=&]")

func domainsGET(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}

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
	// Check if user has write permissions
	session, err := store.Get(r, "session")
	auth, ok := session.Values["authenticated"]
	if err != nil || !ok || !auth.(bool) {
		http.Redirect(w, r, "/login", 303)
		//http.Error(w, "Please log in to view this page", 403)
		//log.Println("Please log in to view this page")
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

	domainName := r.FormValue("name")
	domainName = strings.TrimSpace(domainName)
	if illegalDomainChars.FindString(domainName) != "" {
		http.Error(w, "Please provide domain name in the simplest form (ex. \"onet.pl\" instead of \"http://poczta.onet.pl/\"", http.StatusBadRequest)
		return
	}
	if ipRegex.FindString(domainName) != "" {
		if ipShort.FindString(domainName) == "" {
			http.Error(w, "Please remove all leading zeros from IP address (ex. Change 10.01.01.0 to 10.1.1.0)", http.StatusBadRequest)
			return
		}
	} else {
		domainClean, err := publicsuffix.Domain(domainName)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if domainClean != domainName {
			http.Error(w, fmt.Sprintf("Wrong domain name, did you mean \"%s\"?", domainClean), http.StatusBadRequest)
			return
		}
	}

	err = domain.RegisterNew(domainName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("Successfully added %s, assessment will be available shortly.", domainName)))

	stmt, err := database.DB.Prepare(`SELECT idUrl FROM urls WHERE domain = ?`)
	if err != nil {
		return
	}
	defer stmt.Close()

	var id uint32

	stmt.QueryRow(domainName).Scan(&id)

	path, _ := filepath.Abs("./scrappers")

	cmd := exec.Command("python", "main.py", fmt.Sprintf("%v", id), domainName)
	cmd.Dir = path
	log.Println("Updating " + domainName)
	go cmd.Run()
}

func domainsCreateList(w http.ResponseWriter, r *http.Request) {
	// Check if user has write permissions
	session, err := store.Get(r, "session")
	auth, ok := session.Values["authenticated"]
	if err != nil || !ok || !auth.(bool) {
		http.Error(w, "ePlease log in to view this page", 403)
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

	err = lists.CreateNewList(parsedIDs, r.FormValue("name"), user.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Println(err.Error())
		return
	}
	w.Write([]byte("Successfully created new list!"))
}

func domainsDelete(w http.ResponseWriter, r *http.Request) {
	// Check if user has write permissions
	session, err := store.Get(r, "session")
	auth, ok := session.Values["authenticated"]
	if err != nil || !ok || !auth.(bool) {
		http.Error(w, "aPlease log in to view this page", 403)
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

	err = domain.DeleteDomains(parsedIDs)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte("Domains successfully deleted"))
}
