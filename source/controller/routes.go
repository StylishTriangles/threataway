// Package controller provides handlers and a router for path matching
package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	store *sessions.CookieStore
)

func init() {
	// random keys are generated on startup, no session persistence
	authKey := securecookie.GenerateRandomKey(32)
	encKey := securecookie.GenerateRandomKey(24)
	if authKey == nil || encKey == nil {
		log.Fatal("Could not create session authentication or encryption key")
	}
	store = sessions.NewCookieStore(authKey, encKey)
}

// GetRouter creates a Gorilla Toolkit router with all routes from this controller package
func GetRouter() *mux.Router {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	r.HandleFunc("/", domainsGET)
	r.HandleFunc("/forbidden", forbiddenHandler)
	r.HandleFunc("/lists", listsHandler).Methods("GET")
	r.HandleFunc("/lists", listsDeleteHandler).Methods("POST")
	r.HandleFunc("/deployments", deploymentsHandler).Methods("GET")
	r.HandleFunc("/deployments/add", deploymentsAdd).Methods("POST")
	r.HandleFunc("/deployments/delete", deploymentsDeleteHandler).Methods("POST")
	r.HandleFunc("/loggedin", loggedInHandler)
	r.HandleFunc("/loggedout", loggedOutHandler)
	r.HandleFunc("/login", loginGET).Methods("GET")
	r.HandleFunc("/login", loginPOST).Methods("POST")
	r.HandleFunc("/logout", logoutHandler)
	r.HandleFunc("/templates", templatesHandler).Methods("GET")
	r.HandleFunc("/templates", templatesDeleteHandler).Methods("POST")
	r.HandleFunc("/templates/add", templateCreate).Methods("POST")
	r.HandleFunc("/templates/edit_load", templateLoad)
	r.HandleFunc("/register", registerGET).Methods("GET")
	r.HandleFunc("/register", registerPOST).Methods("POST")
	r.HandleFunc("/registered", registeredHandler)
	r.HandleFunc("/secret", secretGET).Methods("GET")
	r.HandleFunc("/account", accountView) //added for convenience
	r.HandleFunc("/account/view", accountView)
	r.HandleFunc("/account/password", accountPasswordGET).Methods("GET")
	r.HandleFunc("/account/password", accountPasswordPOST).Methods("POST")
	r.HandleFunc("/account/activate", accountActivate).Queries("key", "{key}")
	r.HandleFunc("/domains", domainsGET).Methods("GET")
	r.HandleFunc("/domains/add", domainsAdd).Methods("POST")
	r.HandleFunc("/domains/delete", domainsDelete).Methods("POST")
	r.HandleFunc("/domains/list", domainsCreateList).Methods("POST")

	// Deployment
	r.HandleFunc("/d/{deploy_name}", deployedListHandler)
	r.HandleFunc("/lists/{list_name}", displayListGET).Methods("GET")
	r.HandleFunc("/lists/edit/{list_name}", editListGET).Methods("GET")
	r.HandleFunc("/lists/edit", editListPost).Methods("POST")

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	return r
}
