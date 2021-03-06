package controller

import (
	"threataway/source/model/env"
	"threataway/source/model/user"
	"threataway/source/view"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
)

func timerGET(w http.ResponseWriter, r *http.Request) {
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

	v := view.New("timer")
	tmr := env.GetEnv("timer")
	if tmr != "" {
		v.Vars["timer"] = tmr
	}
	v.Render(w)
}

func timerPOST(w http.ResponseWriter, r *http.Request) {
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

	tmr := r.FormValue("timer")
	env.SetEnv("timer", tmr)
	v := view.New("timer")
	envtmr := env.GetEnv("timer")
	v.Vars["timer"] = envtmr
	v.Vars["ok"] = (tmr == envtmr)
	v.Render(w)
	t := env.GetEnv("timer")
	if t == "" {
		t = "10"
	}
	env.ChangeTimeout(t, func() {
		path, _ := filepath.Abs("./scrappers")

		cmd := exec.Command("python", "main.py", "all")
		cmd.Dir = path
		go cmd.Run()
	})
}
