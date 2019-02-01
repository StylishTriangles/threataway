package controller

import (
	"threataway/source/view"
	"net/http"
)

func forbiddenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	view.New("forbidden").Render(w)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	view.RenderTemplate(w, "404", nil)
}

func errorHandler(w http.ResponseWriter, r *http.Request, code int) {
	if code == 404 {
		notFoundHandler(w, r)
	}
}
