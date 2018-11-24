package controller

import (
	"fmt"
	"gowebapp/source/model/user"
	"gowebapp/source/model/validate"
	"gowebapp/source/view"
	"log"
	"net/http"
	"strings"
)

func registerGET(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "register", nil)
}

func registerPOST(w http.ResponseWriter, r *http.Request) {
	v := view.New("register")
	v.Vars["Username"] = strings.TrimSpace(r.FormValue("username"))
	v.Vars["Password"] = r.FormValue("password")
	v.Vars["Password2"] = r.FormValue("password2")
	v.Vars["Email"] = strings.TrimSpace(r.FormValue("email"))

	// Validate email, username, password
	errCount := 0
	if r.FormValue("password") != r.FormValue("password2") {
		errCount++
		v.Vars["Password2Modal"] = "Passwords don't match"
	}
	if ok, err := validate.Email(v.Vars["Email"].(string)); !ok {
		errCount++
		switch err {
		case validate.ErrEmpty:
			v.Vars["EmailModal"] = "This field is required"
		case validate.ErrInvalidFormat:
			v.Vars["EmailModal"] = "Please enter a valid email address"
		case validate.ErrTooLong:
			v.Vars["EmailModal"] = fmt.Sprintf("E-mail cannot be longer than %d characters", validate.EmailMaxLen)
		}
	}
	if ok, err := validate.Username(v.Vars["Username"].(string)); !ok {
		errCount++
		switch err {
		case validate.ErrEmpty:
			v.Vars["UsernameModal"] = "This field is required"
		case validate.ErrTooLong:
			v.Vars["UsernameModal"] = fmt.Sprintf("Username cannot be longer than %d characters", validate.UsernameMaxLen)
		case validate.ErrInvalidFormat:
			v.Vars["UsernameModal"] = "Username must only contain alphanumeric characters and ._-"
		}
	}
	if ok, err := validate.Password(v.Vars["Password"].(string)); !ok {
		errCount++
		switch err {
		case validate.ErrEmpty:
			v.Vars["PasswordModal"] = "This field is required"
		case validate.ErrTooShort:
			v.Vars["PasswordModal"] = fmt.Sprintf("Password must be at least %d characters long", validate.PasswordMinLen)
		case validate.ErrTooLong:
			v.Vars["PasswordModal"] = fmt.Sprintf("Password cannot be longer than %d characters", validate.PasswordMaxLen)
		case validate.ErrCommonPassword:
			v.Vars["PasswordModal"] = "Password is too trivial, please avoid using common passwords"
		}
	}

	if errCount != 0 {
		v.AddFlash(fmt.Sprintf("Error: %d inputs are invalid", errCount), view.FlashError)
		v.Render(w)
	} else { // Attemt to register user when all inputs are correct
		err := user.Register(v.Vars["Username"].(string), v.Vars["Email"].(string), v.Vars["Password"].(string))
		switch err {
		case nil: // no problem, registration success
			http.Redirect(w, r, "/registered", http.StatusFound)
		case user.ErrAlreadyRegistered: // User already exists in database
			v.Vars["UsernameModal"] = "Username or email already in use"
			v.Render(w)
		case user.ErrInvalidEmail:
			v.Vars["EmailModal"] = "Please enter a valid email address"
			v.Render(w)
		default: // sql error/unknown error
			log.Println(r.RemoteAddr, err)
			http.Error(w, "Your request could not be processed at this time", http.StatusInternalServerError)
		}
	}
}

func registeredHandler(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "registered", nil)
}
