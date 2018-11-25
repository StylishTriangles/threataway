package controller

import (
	"log"
	"net/http"

	"gowebapp/source/model/user"
	"gowebapp/source/shared/password"
	"gowebapp/source/view"
)

func loginGET(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "login", nil)
}

// TODO: cleanup
func loginPOST(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	passwd := r.FormValue("password")
	u, err := user.GetByHandle(login)
	if err != nil || !password.Compare(passwd, u.PasswordHash) {
		if err != nil && err != user.ErrNotFound { // database error
			log.Println(err)
			http.Error(w, "There was a problem processing your request, please try again later.", http.StatusInternalServerError)
			return
		}
		view.RenderTemplate(w, "login", struct {
			Login        string
			LoginFailure bool
		}{
			Login:        login,
			LoginFailure: true,
		})
		return
	}
	if u.Active == 0 { // Check if account is active
		w.Write([]byte("Please activate your account first"))
		return
	}
	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Values["user"] = u
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Session could not be saved", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/domains", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = false
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Session could not be saved", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/loggedout", http.StatusFound)
}

func loggedInHandler(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "loggedin", nil)
}

func loggedOutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}
