package controller

import (
	"net/http"
)

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	session, err := store.Get(r, "session")
	auth, ok := session.Values["authenticated"]
	if err != nil || !ok || !auth.(bool) {
		//http.Error(w, "Please log in to view this page", 403)
		//log.Println("Please log in to view this page")
		http.Redirect(w, r, "/login", 303)

		return false
	}
	return true
}
