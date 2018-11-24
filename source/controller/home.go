package controller

import (
	"gowebapp/source/view"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "home", nil)
}
