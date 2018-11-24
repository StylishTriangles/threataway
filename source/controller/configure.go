package controller

import (
	"gowebapp/source/view"
	"net/http"
)

func configureGET(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "configure", nil)
}

func configurePOST(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "configure", nil)
}

func configuredHandler(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "configure", nil)
}
