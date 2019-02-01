package controller

import (
	"threataway/source/model/user"
	"threataway/source/shared/password"
	"threataway/source/view"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func accountView(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	details, ok := session.Values["user"].(*user.User)
	if !ok {
		http.Redirect(w, r, "/logout", http.StatusFound)
		return
	}

	view.RenderTemplate(w, "account/view", details)
}

func accountPasswordGET(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "account/change_password", nil)
}

// TODO: make this work
func accountPasswordPOST(w http.ResponseWriter, r *http.Request) {
	form := struct{ oldPassword, newPassword1, newPassword2 string }{
		oldPassword:  r.FormValue("oldpassword"),
		newPassword1: r.FormValue("newpassword1"),
		newPassword2: r.FormValue("newpassword2"),
	}
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	details, ok := session.Values["user"].(*user.User)
	if !ok {
		http.Redirect(w, r, "/logout", http.StatusFound)
		return
	}
	if !ok || !password.Compare(form.oldPassword, details.PasswordHash) {
		// Invalid password
	}
	if form.newPassword1 != form.newPassword2 {
		// Password mismatch
	}
	view.RenderTemplate(w, "password_changed", nil)
}

func accountActivate(w http.ResponseWriter, r *http.Request) {
	v := view.New("account/activate")
	vars := mux.Vars(r)
	k, ok := vars["key"]
	if ok {
		err := user.Confirm(k)
		switch err {
		case nil:
			break
		case user.ErrInvalidActivator:
			v.Vars["InvalidKey"] = true
		case user.ErrActivatorExpired:
			v.Vars["ExpiredKey"] = true
		default: // Database error
			http.Error(w, "Could not process your request at this time, please try again later", http.StatusInternalServerError)
			log.Println("accountActivate error:", err)
			return
		}
	} else {
		v.Vars["InvalidKey"] = true
	}
	v.Render(w)
}
